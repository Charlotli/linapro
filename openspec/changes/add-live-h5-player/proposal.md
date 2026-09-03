# 直播 H5 播放页

## Why

直播管理插件（`linapro-live-manage`）目前只提供后台维护能力，观众无法观看直播。管理端录入的播放地址（`live_url`）为`.m3u8`（HLS）格式，需要一个面向观众的 H5 播放页面，通过公开、免登录的方式消费这些数据，完成"后台维护 → 观众观看"的业务闭环。原`add-live-manage-plugin`变更的`design.md`已将"播放器与观众端页面"列为非目标，本变更作为其后续迭代补齐该能力。

## What Changes

- 新增观众端公开播放信息接口（免登录、免权限点）：
  - `GET /x/linapro-live-manage/api/v1/play?roomCode=<roomCode>&tenantId=<tenantId>`：按直播间编码查询该房间当前可播放的直播，优先返回`state=1`（进行中）的公开直播，其次回退最近一场`state=2`（已结束）的公开直播作为回放；`liveUrl`仅在进行中或回放态返回，未开始态不返回播放地址。
  - 匿名请求无法携带租户身份，接口通过显式`tenantId`查询参数解析租户上下文；租户插件启用时该参数必填并校验租户存在且未停用，租户插件未启用时忽略该参数并固定`tenant_id=0`。
- 插件后端新增公开路由分组：复用宿主发布的`NeverDoneCtx`、`CORS`、`Ctx`中间件（JSON 接口另含`HandlerResponse`、`RequestBodyLimit`），不挂`Auth`/`Tenancy`/`Permission`。
- 插件前端新增 H5 播放页（`frontend/h5/` 静态资产，随`go:embed`嵌入）：
  - 使用内嵌的`hls.js`库播放`.m3u8`（iOS Safari 走原生 HLS，其余浏览器走`hls.js`）。
  - 播放页 URL 形如`/x/linapro-live-manage/h5?room=<roomCode>&tenant=<tenantId>`，由插件注册的公开路由直接提供静态文件服务，URL 不含插件版本号，插件升级不破坏已分享链接。
  - 页面展示直播标题、直播间名称、封面、直播中/未开始/已结束状态，以及歌单、讲道主题、讲员、经文等节目信息（按有值显示）；未开始或不可播时按 30 秒间隔轮询刷新。
- 插件`apidoc`翻译资源同步新增观众端接口的`zh-CN`翻译；插件错误语言包新增播放业务错误码翻译（插件`i18n.enabled: true`）。
- 新增播放选择逻辑单元测试、API 契约测试断言与 E2E 用例（`TC004-live-h5-play.ts`）。

## Capabilities

### New Capabilities

- `live-manage` 能力扩展：观众端 H5 播放（公开播放信息查询 + 内嵌 H5 播放页）。

### Modified Capabilities

- 无修改既有能力；管理端接口、菜单、权限点均不变。

## Impact

- 影响范围限定在`apps/lina-plugins/linapro-live-manage/`内部与`openspec/changes/add-live-h5-player/`变更文档；不修改`lina-core`与既有插件，符合核心宿主边界要求。
- `i18n`影响判断：插件启用`i18n`，新增公开接口的`apidoc`英文源文本与`zh-CN`翻译资源；H5 播放页为独立轻量页面，不接入宿主前端语言包，页面文案使用中文（观众端目标语言），不引入`$t`依赖。
- 缓存一致性影响判断：公开播放信息接口为实时读库查询（单行、租户+状态+公开性条件命中既有索引`idx_plugin_linapro_live_manage_live_tenant_state`），无缓存、快照或跨实例协调状态，无影响。
- 数据权限影响判断：公开播放接口属于"公开资源"例外——匿名观众无用户身份，不存在角色数据权限语义；租户边界通过显式`tenantId`参数约束，且仅暴露`is_public=1`且`state=1`的记录，私有/未开始/已结束/跨租户记录一律按不存在处理，例外理由与拒绝策略见`design.md`。
- 开发工具跨平台影响判断：无新增脚本或构建工具变更；`hls.js`库文件直接放入插件前端资产目录，随`go:embed`嵌入，无跨平台影响。
- 测试策略影响判断：新增服务层单元测试（租户上下文构造、公开性/状态过滤、错误路径）、API 契约测试（最小投影断言）与 E2E 用例（公开页可访问、播放器初始化、私密直播不可见）。
- `openspec validate`门禁：`openspec`工具未安装（执行环境验证`openspec: command not found`），记录为不可用，变更文档由人工按规范维护。
