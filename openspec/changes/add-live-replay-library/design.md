# 设计：直播往期回放库

## 模块归属

- 全部落在插件`linapro-live-manage`内部：直播表一列、`play`回退过滤、公开`replays`列表接口、H5 静态页区块、管理端既有页面扩展。不触碰`apps/lina-core`，不新增菜单与权限点（管理端复用既有直播内容页与权限）。
- 后端实现放在既有`play`服务组件同包内（回放列表与播放选择是同一"观众端公开投影"职责域，不新建组件）；房间解析、租户上下文、公开可见性判定复用既有共享查询。无新增抽象层。

## 数据库设计

- `plugin_linapro_live_manage_live`新增列`replay_enabled BOOLEAN NOT NULL DEFAULT TRUE`：场级回放开关，语义"该场直播结束后是否提供公开回放"；`DEFAULT TRUE`保证既有数据与既有行为不变。
- 插件`manifest/sql/`新增迭代 SQL（`ALTER TABLE ... ADD COLUMN IF NOT EXISTS`幂等），`manifest/sql/uninstall/`同步`ALTER TABLE ... DROP COLUMN IF EXISTS`；随后`make db.init`、`make dao`重生成。
- 索引评估：`replays`列表查询路径为"房间 + `state=2` + 公开 + `replay_enabled` + 按`live_date`倒序分页"，与既有`play`查询同形（房间维度过滤），沿用既有直播表房间查询索引；数据规模为单房间历史场次（教会场景每年量级），无需新索引，设计阶段记录该判断，实施时用`EXPLAIN`抽查确认。

## 接口设计

### 公开往期回放列表

`GET /x/linapro-live-manage/api/v1/replays?roomCode=&tenantId=&page=&pageSize=`，挂既有公开路由分组（免登录、免权限点）。

- 过滤：房间匹配 + `state=2` + `is_public` + `replay_enabled=true` + 租户上下文（语义与`play`一致：租户启用时`tenantId`必填并校验存在启用，未启用忽略固定 0）。
- 响应最小投影：`liveId`、`title`、`coverUrl`、`liveDate`、`startTime`/`endTime`（Unix 毫秒，`*int64`，遵循时间戳契约）、`liveUrl`（已结束公开回放返回播放地址，与`play`回放态同边界）。
- 分页：`page`默认 1、`pageSize`默认 10 上限 50（超限按上限截断并在`dc`说明）；单条分页 SQL（数据库侧过滤排序分页），访问次数恒为 1，无 N+1。
- 无匹配内容返回空列表与零总数（HTTP 200），与房间不存在同形态，不泄露存在性。

### `play`接口收紧

- 已结束回退分支增加`replay_enabled=true`条件：管理端关闭某场回放后，该场既不出现在`replays`列表，也不再作为"最近回放"被`play`返回，两处口径一致。
- 进行中直播（`state=1`）不受该开关影响，直播中无"关闭直播"语义。

### 管理端

- 直播内容列表响应新增`replayEnabled`；`PUT /live/{id}`请求 DTO 新增可选`replayEnabled`字段（不传不改，避免既有编辑路径回写覆盖）。
- 行内开关直接调用既有更新接口，成功后刷新行数据；批量操作不涉及。

## 前端设计

### H5 观播页

- 回放态在信息面板下方新增`#replay-list`区块：调`replays`接口首页（`pageSize=10`），条目为封面缩略图（无封面用既有光环占位）+ 标题 + 日期；点击条目用其`liveUrl`切换播放器流（复用`attachStream`）并更新标题、日期与徽标，当前播放场高亮描边。
- 区块仅在回放态渲染；列表为空或接口失败静默不渲染，不影响播放主流程。进行中/预告态不请求该接口。

### 管理端

- 直播内容列表新增"回放"列（`Switch`行内开关，仅`state=2`行可操作，未结束行显示"—"）；新建/编辑弹窗在播放地址下方新增"结束后提供公开回放"开关。全部走既有页面与权限。

## i18n 影响

- 插件`i18n.enabled: true`：`replays`接口英文源文本 + `manifest/i18n/<locale>/apidoc/`同步`zh-CN`翻译；管理端新增列/开关文案按插件既有文案机制维护。
- H5 静态页中文独立页（既有先例），无宿主语言包依赖。

## 数据权限例外说明

- 公开`replays`接口为匿名读取的合理例外：权威边界与`play`一致（房间存在 + 公开 + `replay_enabled` + 显式租户参数），输出仅日程与回放地址投影；无效输入返回空列表，不泄露存在性。管理端`replayEnabled`修改走既有直播内容更新接口，权限边界不变。

## 复杂度判断

- 一列、一个查询接口、两处前端扩展，局部直接实现；不新建组件、不引入配置模型。`play`回退与列表口径一致性由单元测试锁定。
