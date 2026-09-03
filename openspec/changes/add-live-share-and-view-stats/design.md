# 设计：直播分享与在线观看统计

## 模块归属

- 全部变更落在业务插件`linapro-live-manage`内部：后端新服务组件、插件 H5 静态页、插件管理端页面。不触碰`apps/lina-core`核心领域契约，不新增插件间契约。
- 观看统计是独立职责（统计域），新建`view`服务组件承载；不并入`play`服务（`play`只负责播放信息查询），避免单一组件职责膨胀。跨组件复用播放选择逻辑（按房间解析当前可播直播）通过同包内部函数共享，不新造窄接口。

## 数据库设计

新增表`plugin_linapro_live_manage_view_session`（观看会话）：

| 列 | 类型 | 说明 |
| --- | --- | --- |
| `id` | `BIGSERIAL` | 自增主键 |
| `tenant_id` | `INT NOT NULL DEFAULT 0` | 租户隔离，语义与既有房间/直播表一致 |
| `live_id` | `BIGINT NOT NULL` | 关联`plugin_linapro_live_manage_live.id` |
| `session_key` | `VARCHAR(64) NOT NULL` | 匿名观众会话键，浏览器侧生成并持久于`sessionStorage` |
| `created_at` / `updated_at` | `TIMESTAMPTZ` | GoFrame 自动时间维护 |

- **软删除决策**：不设`deleted_at`。会话是统计流水，无"恢复误删"诉求；过期数据按保留策略直接硬删。规则要求的软删除确认点在此明确记录。
- **唯一约束**：`uk_plugin_linapro_live_manage_view_session_session_key (session_key)`，是心跳`upsert`的真实业务键（一个浏览器会话一行）。
- **索引**：`idx_plugin_linapro_live_manage_view_session_live_updated (live_id, updated_at)`，同时服务在线计数（`live_id`等值 + `updated_at`范围）与累计计数（`live_id`前缀），避免统计聚合全表扫描。
- **时间窗口跨库通用**：在线判定阈值由 Go 侧计算`time.Now().Add(-60*time.Second)`后作为参数传入`GT(updatedAt, threshold)`，不使用数据库方言函数。
- **SQL 文件**：插件`manifest/sql/`新增`{序号}-live-share-and-view-stats.sql`（`CREATE TABLE IF NOT EXISTS` + `CREATE INDEX IF NOT EXISTS`），`manifest/sql/uninstall/`同步`DROP TABLE IF EXISTS`；随后`make db.init`、`make dao`重新生成。无 Seed DML 与 Mock 数据（统计表由运行时写入）。
- **保留策略**：不引入定时清理任务（教会场景会话规模有界且单行极小）；如未来需要，可独立迭代增加清理 Cron，本期不做。

## 后端设计

### 公开心跳接口

`POST /x/linapro-live-manage/api/v1/view/heartbeat`，挂既有公开路由分组（`NeverDoneCtx`、`CORS`、`Ctx`、`HandlerResponse`、`RequestBodyLimit`，无`Auth`/`Tenancy`/`Permission`）。

- 请求：`roomCode`（必填）、`tenantId`（可选，租户插件启用时必填并校验存在且未停用，未启用时忽略并固定`0`，与`play`接口语义一致）、`sessionKey`（必填，长度 1~64）。
- 处理：复用播放选择逻辑解析当前可播直播（进行中优先，回退最近已结束）；无可播直播或房间不存在时返回零值统计（不报错、不落库）——对无效`roomCode`与"暂无可播直播"返回相同形态，不泄露房间存在性；可播则按`session_key`幂等 upsert 会话（已存在则必要时更新`live_id`，不存在则插入，唯一索引兜底并发）。
- 响应：`onlineCount`（该直播最近 60 秒有心跳的会话数）、`totalViews`（该直播累计去重会话数）。均为`int64`，单条聚合 SQL 内完成（`COUNT`+`SUM CASE`，不用 PG 方言`FILTER`）。
- 时间字段契约：响应无时间点字段；无新增业务错误码（失败形态只有参数校验与静默零值），`backend-go`的`bizerr`要求不触发新码定义。
- 数据库访问次数：每次心跳固定 3 次以内（解析可播直播 1 次 + upsert 1 次 + 聚合 1 次），与观众规模无关。

### 管理端列表统计装配

- 既有直播内容列表接口响应新增`onlineCount`、`totalViews`两个字段；在列表 service 装配阶段收集当前页`liveId`集合，执行单条聚合 SQL（`WHERE live_id IN (...) GROUP BY live_id`）后内存映射合并，空列表跳过查询。
- 明确禁止逐行查询；数据库访问次数恒为 1 次，与页大小无关，满足接口性能与`N+1`治理要求。

### 数据权限例外说明

- 公开心跳接口为匿名写入的合理例外：权威边界是"`roomCode`对应房间存在且有公开可播直播"（复用`play`的公开可见性判定），响应仅含聚合计数，不返回任何记录明细；无效输入与无直播同形态返回，不泄露范围外数据存在性。对应验证由契约测试覆盖。
- 管理端统计列随既有直播内容列表接口一并返回，权限边界与该列表完全一致（同一响应装配，不新增独立数据面），不构成新的数据权限面。

## 前端设计

### H5 观播页（`frontend/h5/`）

- **分享操作条**：信息面板底部新增`.share-bar`。`navigator.share`可用时展示"分享直播"按钮（调用`navigator.share({title, url})`）；不可用时展示"复制链接"按钮（`navigator.clipboard.writeText`，失败降级`document.execCommand('copy')`）；成功与失败均以自制轻提示展示结果，不引入任何新依赖。样式延续既有大厂暗色视觉（渐变胶囊按钮）。
- **心跳与在线人数**：
  - `sessionKey`：首次进入生成随机键（优先`crypto.randomUUID`，降级时间戳+随机串）并存入`sessionStorage`（键`linapro-live-view-session`），同一浏览器会话刷新页面不重复计数。
  - 节律：仅在展示可播内容（进行中/回放）时启动独立 30 秒心跳定时器（与既有状态轮询定时器分离，避免与预告/异常态轮询语义纠缠）；切态或销毁播放器时停止。
  - 展示：心跳成功后，进行中状态在状态胶囊附近展示"N 人在看"徽标；回放态不展示实时人数；心跳失败静默保留上一次数值，不影响观播。

### 管理端（`frontend/pages/`）

- **直播间二维码**：直播间列表操作列新增"二维码"动作，弹窗内展示二维码图片、H5 链接文本与"复制链接"按钮。二维码由前端生成（新增轻量依赖`qrcode`，无框架绑定），链接按`${location.origin}/x/linapro-live-manage/h5?room=<roomCode>`拼装；房间行`tenantId`大于 0 时追加`&tenant=<tenantId>`。实施时需确认既有列表 DTO 是否已返回`tenantId`，若无则在既有列表投影中补充该字段，禁止为取租户号逐行补查。
- **统计列**：直播内容列表新增"在线人数""累计观看"两列，读取列表响应新字段，无额外请求。

## i18n 影响

- 插件`i18n.enabled: true`：新增心跳接口的`apidoc``zh-CN`翻译必须同步进`manifest/i18n/<locale>/apidoc/`；管理端新增列与按钮文案按插件前端既有文案机制维护，实施时确认键命名。
- H5 静态页为中文独立页（既有先例，不经宿主语言包），新增文案跟随既有做法，不引入宿主`i18n`依赖。

## 枚举与启停

- 无新增枚举语义值（心跳与统计均为数值口径，无状态/类型字典需求）。
- 插件禁用：公开路由与页面随插件整体下线（宿主机制），无需额外降级；租户能力禁用时`tenantId`忽略并固定`0`，统计口径不变。
- 统计表缺失或心跳链路异常：H5 静默降级（不显示人数徽标、不影响播放），管理端统计列显示`0`，不产生 500 或空白。

## 复杂度判断

- 新增的`view`服务组件承载真实新职责（统计写读聚合），局部直接实现即可，无新增抽象层、适配器或配置模型。
- H5 分享条、管理端二维码弹窗均为既有页面局部扩展，不引入状态管理或路由变化。
