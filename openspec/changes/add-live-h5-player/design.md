## Context

`add-live-manage-plugin`变更交付了直播间与直播内容的管理能力，`design.md`将"播放器与观众端页面"明确列为非目标。本次迭代补齐观众端：管理端录入的`live_url`为`.m3u8`（HLS）格式，需要提供免登录的 H5 播放页。观众无用户身份、无 JWT，因此不能复用挂载`Auth`/`Tenancy`/`Permission`的管理端路由；直播数据按`tenant_id`隔离，匿名请求需要显式的租户定位手段。

## Goal

观众在手机浏览器打开分享链接，即可观看指定直播间当前进行中的公开直播（HLS），或在直播未开始/已结束时看到相应的节目预告与回放信息。

## 非目标

- 不做公开直播列表页、点赞、评论、弹幕等互动能力（原需求确认仅播放页）。
- 不做房间访问码、观看鉴权（原需求确认仅公开直播，私密直播由`is_public=0`天然隔离）。
- 不做推流端能力与转码，HLS 切片由外部流媒体服务（如 SRS、云厂商直播）负责。
- 不接入宿主管理端语言包运行时，H5 页面文案为单一中文（见`i18n`设计决策）。

## 设计决策

### 归属与边界

- 归属：源码插件能力（`apps/lina-plugins/linapro-live-manage/`），业务前后端严格闭环在插件内部，不修改`lina-core`。宿主已发布的`pluginhost.RouteMiddlewares`中间件目录与`RouteRegistrar`路由注册器足以支撑公开路由，无需扩展宿主契约。
- 新增组件：`api/play`（DTO 与路由契约）、`internal/controller/play`（HTTP 控制器）、`internal/service/play`（播放查询服务）、`frontend/h5/`（静态 H5 资产）。

### 路由与中间件

| 分组 | 中间件 | 路由 | 用途 |
| --- | --- | --- | --- |
| 公开 JSON | `NeverDoneCtx` + `HandlerResponse` + `CORS` + `RequestBodyLimit` + `Ctx` | `GET /api/v1/play` | 播放信息查询，统一响应封装 |
| 公开静态 | `NeverDoneCtx` + `CORS` + `Ctx` | `GET /h5`、`GET /h5/*any` | H5 页面与`hls.js`等静态资产 |

- 公开分组不挂`Auth`/`Tenancy`/`Permission`：观众无 JWT；`Permission`依赖用户身份无意义；`Tenancy`依赖请求元数据解析租户，而观众分享链接无法稳定携带租户域名或请求头，租户上下文改为由接口参数显式传递（见下文）。
- 静态分组不挂`HandlerResponse`（避免 JSON 封装 HTML 输出）与`RequestBodyLimit`（只读 GET）。
- 未采用宿主`/x-assets/{plugin}/{version}/`公开资产通道：该 URL 携带插件版本号，插件升级会使已分享链接（如印刷二维码）失效；插件自注册路由提供版本无关的稳定 URL。未采用`plugin.yaml`的`public_assets`声明：本插件自服务静态文件，无需占用宿主公开资产治理面。

### 租户上下文

- `room_code`唯一性为`(tenant_id, room_code)`组合，跨租户可能重码，因此公开查询必须携带租户定位：`tenantId`查询参数。
- 租户插件未启用（`tenantSvc.Available()`为`false`）：全部数据`tenant_id=0`，参数被忽略，固定按`tenant_id=0`查询。
- 租户插件启用：`tenantId`必填；`tenantId=0`表示平台租户数据；`tenantId>0`时通过`tenantSvc.Directory().Get()`校验租户存在且状态为启用，校验失败返回`CodePlayTenantInvalid`。
- 服务层不依赖请求`bizctx`租户，显式以`Where(tenant_id, ?)`下推到数据库查询，保证数据权限过滤在查询阶段完成。

### 公开播放信息接口

`GET /api/v1/play`，`tags: LivePlay`，无`permission`标签（公开接口，公开分组不挂权限中间件；插件`api_contract_test.go`不强制插件 DTO 的权限标签，宿主静态权限审计仅覆盖宿主 API 包，已确认）。

请求参数：

| 参数 | 必填 | 说明 |
| --- | --- | --- |
| `roomCode` | 是 | 直播间编码，同租户内唯一 |
| `tenantId` | 见租户上下文 | 租户 ID；租户插件启用时必填 |

候选选择逻辑（同租户、同房间，全部在数据库侧过滤）：

1. `state=1 AND is_public=1`的直播（进行中，取`id`最新一场）→ 返回完整播放信息含`liveUrl`。
2. 无进行中时，回退`state=2 AND is_public=1`的最近一场（回放，HLS 切片通常可继续作为点播播放）→ 返回含`liveUrl`。
3. 无进行中且无已结束公开直播，但存在`state=0 AND is_public=1`的最近一场（预告）→ 返回节目信息，`liveUrl`置空，避免提前暴露播放地址。
4. 都不存在、直播间不存在、或直播间已被禁用且无任何可展示记录 → `CodePlayNotFound`（HTTP 404）。

实现为定点查询序列：先查房间（`tenant_id+room_code`命中`uk_plugin_linapro_live_manage_room_tenant_code`），再按优先级顺序（进行中 → 回放 → 预告）逐态查询直播（`tenant_id+room_id+state`命中`idx_plugin_linapro_live_manage_live_tenant_room`与`idx_plugin_linapro_live_manage_live_tenant_state`），命中即停。数据库访问次数以候选状态数为上界（固定 3 次），无`N+1`。

进行中态额外校验直播间状态必须为`直播中`（`status=1`，`Start`/`Stop`状态机保证），直播间被禁用时进行中直播不可播。该状态校验在已查出的房间行上内存判断，不增加查询。

响应投影（最小必要）：`roomId`、`roomCode`、`roomName`、`liveId`、`title`、`coverUrl`、`liveUrl`、`liveDate`（`YYYY-MM-DD`，`date-only`语义）、`startTime`（Unix 毫秒）、`state`、`host`、`songName`、`songList`、`leadSinger`、`accompaniment`、`sermonTitle`、`preacher`、`preacherIdentity`、`scriptureRef`、`scriptureContent`、`outline`。明确排除`pushUrl`（推流地址绝不外泄）、`pageUrl`、`deviceInfo`、`reception`、审计字段与`isPublic`（恒为 1 无信息量）。

### 错误码

在`service/play/play_code.go`集中定义：

| 错误码 | 场景 |
| --- | --- |
| `LIVE_MANAGE_PLAY_NOT_FOUND` | 直播间不存在或无可播放/预告/回放的公开直播 |
| `LIVE_MANAGE_PLAY_TENANT_REQUIRED` | 租户插件启用时未携带`tenantId` |
| `LIVE_MANAGE_PLAY_TENANT_INVALID` | `tenantId`对应的租户不存在或已停用 |

### H5 播放页

`frontend/h5/`静态资产（`index.html`、`hls.light.min.js`、`app.js`、`app.css`），由插件公开静态路由直接从`livemanage.EmbeddedFiles`读取（`path.Join("frontend/h5", …)`），`mime.TypeByExtension`判定`Content-Type`（兜底`http.DetectContentType`），路径经清理防止目录穿越。`frontend/`整体已在`plugin_embed.go`的`//go:embed`范围内，无需修改嵌入声明。

- 播放器：优先`hls.js`（`Hls.isSupported()`，采用`hls.light.min.js`轻量发行版），iOS Safari 等原生 HLS 环境直接`video.src`，两者都不可用时展示不支持提示。`hls.js`为内嵌资产，不依赖外网 CDN。
- 状态呈现：`state=1`播放器 + "直播中"标识；`state=0`封面预告 + 直播日期；`state=2`回放播放器 + "已结束"标识；404 或无数据显示"暂无直播"，网络异常与无效链接分别提示。
- 轮询：`state=0`或"无直播/网络异常"时每 30 秒自动刷新一次（管理端开启直播后观众页自动进入播放），播放中不轮询；"开启直播"后由下一次轮询自然接入，E2E 用刷新验证即时态。
- 节目信息：歌单（`songList`JSON 数组解析，解析失败按空处理并忽略）、讲道（主题/讲员/身份/经文/大纲）、主持/领唱/伴奏，仅有值时渲染。
- API 地址推导：以`location.pathname`中`/h5`之前的部分为插件前缀拼接`/api/v1/play`，兼容反向代理子路径部署。
- 目录位置说明：`frontend/`下现有约定为`pages/`（管理端 Vue 页）与`slots/`；`frontend/h5/`为独立的观众端静态页，不进入管理端动态页壳，属于对前端目录的受控扩展，在 design 记录即可。

### i18n

- 插件`i18n.enabled: true`（`zh-CN`/`en-US`）：新增`api/play/v1`英文源文本，并在`manifest/i18n/zh-CN/apidoc/plugin-api-main.json`补齐翻译；`en-US/apidoc`按规则保持空对象占位。
- 错误语言包：`manifest/i18n/{zh-CN,en-US}/error.json`新增 3 个播放错误码的翻译。
- H5 页面文案：观众端独立静态页，不接入管理端`$t`运行时，页面文案为中文（插件默认语言）。国际化播放页需要语言协商与语言包加载链路，超出本次范围，记录为后续迭代可选项。

### 数据权限例外说明

- 权威边界：仅`is_public=1`的直播内容可被匿名访问；租户隔离由显式`tenantId`参数约束；进行中直播要求房间状态为`直播中`。
- 拒绝策略：不满足边界的记录（私密、未开始的播放地址、其他租户、不存在的房间）统一返回`LIVE_MANAGE_PLAY_NOT_FOUND`，不泄露存在性差异。
- 为什么可以绕过常规数据权限：观众无用户身份，角色数据权限不适用；该接口属于`data-permission.md`所述"公开资源"合理例外，暴露的每一行都是管理端显式标记为公开的内容。
- 对应测试覆盖：单元测试覆盖候选选择与拒绝路径；E2E 覆盖"公开直播可播、私密直播 404、跨租户不可见"。

### 依赖注入

- `service/play`构造函数显式注入`tenantcap.Service`（租户存在性/启用状态校验）；不引入其他运行期依赖，无新增共享实例或缓存协调诉求。
- H5 静态处理器无状态，直接闭包引用插件`livemanage.EmbeddedFiles`包级变量（`fs.ReadFile`按需读取），不注入`fs.FS`参数；来源为插件嵌入资源。

## 影响分析

- 缓存一致性：无缓存、快照或跨实例协调状态，实时读库（定点索引查询），无影响。
- 数据库：无表结构变更、无 SQL 文件变更；查询命中既有唯一索引与组合索引，无全表扫描。
- 枚举治理：复用既有`plugin_live_state`字典语义，无新增枚举。
- 模块启停：插件禁用后路由与公开页随插件整体下线（宿主插件生命周期统一治理），无独立降级逻辑需求。
