## ADDED Requirements

### Requirement: 观众端公开播放信息查询

系统 SHALL 提供免登录的公开播放信息查询接口`GET /api/v1/play`，按`roomCode`与`tenantId`定位直播间，并按以下优先级返回该房间当前可展示的公开直播：进行中（`state=1`）优先，其次最近一场已结束（`state=2`）公开直播作为回放，再次最近一场未开始（`state=0`）公开直播作为预告。私有（`is_public=0`）直播内容 MUST NOT 通过该接口返回；进行中直播 MUST 在直播间状态为`直播中`时才返回播放地址；未开始直播 MUST NOT 返回播放地址。所有不满足边界的请求 SHALL 统一返回`LIVE_MANAGE_PLAY_NOT_FOUND`，不泄露记录存在性。

#### Scenario: 查询进行中的公开直播

- **WHEN** 观众以正确参数查询一个存在进行中公开直播的直播间
- **THEN** 系统返回该直播的最小投影（标题、直播间名称、封面、播放地址、状态、节目信息）
- **AND** 响应不包含`pushUrl`、`pageUrl`等管理端字段

#### Scenario: 无进行中直播时回放最近一场已结束公开直播

- **WHEN** 直播间没有进行中直播，但存在已结束的公开直播
- **THEN** 系统返回最近一场已结束公开直播的信息与播放地址，状态为已结束

#### Scenario: 未开始直播仅返回预告信息

- **WHEN** 直播间只有未开始的公开直播
- **THEN** 系统返回节目预告信息，播放地址为空
- **AND** 直播开始后查询返回进行中信息与播放地址

#### Scenario: 私密与其他租户直播不可见

- **WHEN** 观众查询私密直播所在直播间，或以错误`tenantId`查询，或查询不存在的直播间
- **THEN** 系统统一返回`LIVE_MANAGE_PLAY_NOT_FOUND`，不区分具体原因

### Requirement: 公开接口租户上下文

公开播放信息接口 SHALL 通过显式`tenantId`参数定位租户上下文。租户能力启用时`tenantId`必填且 MUST 校验租户存在并处于启用状态，校验失败返回`LIVE_MANAGE_PLAY_TENANT_INVALID`，缺失返回`LIVE_MANAGE_PLAY_TENANT_REQUIRED`；租户能力未启用时参数被忽略并固定`tenant_id=0`。租户过滤 MUST 在数据库查询阶段完成。

#### Scenario: 租户参数缺失被拒绝

- **WHEN** 租户插件启用且请求未携带`tenantId`
- **THEN** 系统返回`LIVE_MANAGE_PLAY_TENANT_REQUIRED`

#### Scenario: 停用租户内容不可访问

- **WHEN** 观众以已停用租户的`tenantId`查询
- **THEN** 系统返回`LIVE_MANAGE_PLAY_TENANT_INVALID`

#### Scenario: 单租户部署无需租户参数

- **WHEN** 租户插件未启用时观众仅携带`roomCode`查询
- **THEN** 系统按平台租户`tenant_id=0`完成查询

### Requirement: 观众端 H5 播放页

系统 SHALL 通过插件公开路由提供版本无关稳定 URL 的 H5 播放页（`GET /h5`），该页面 MUST 免登录访问，MUST 支持使用`hls.js`（内嵌资产，不依赖外网 CDN）播放`.m3u8`并在原生 HLS 环境降级为原生播放。页面 SHALL 呈现直播中、未开始（预告）、已结束（回放）与无直播四种状态，未开始或无直播时 SHALL 周期性自动刷新。页面 MUST NOT 展示推流地址等管理端字段。

#### Scenario: 打开播放页观看进行中直播

- **WHEN** 观众在手机浏览器打开携带`roomCode`与`tenantId`的播放页链接且该房间有进行中公开直播
- **THEN** 页面展示播放器并开始播放`.m3u8`流，展示"直播中"状态与节目信息

#### Scenario: 未开始直播显示预告并自动刷新

- **WHEN** 直播间仅有未开始的公开直播
- **THEN** 页面展示封面预告与直播日期，不出现播放器
- **AND** 页面按固定周期自动刷新，直播开始后自动进入播放状态

#### Scenario: 无直播时展示友好提示

- **WHEN** 直播间不存在或没有任何公开直播
- **THEN** 页面展示"暂无直播"提示并支持手动重试

### Requirement: 公开路由中间件边界

插件的公开路由分组 SHALL 仅复用宿主发布的与身份无关的中间件（`NeverDoneCtx`、`HandlerResponse`、`CORS`、`RequestBodyLimit`、`Ctx`），MUST NOT 挂载`Auth`、`Tenancy`或`Permission`中间件；静态资产分组 SHALL NOT 挂载`HandlerResponse`。公开路由能力 MUST 完全由插件自身组合，不修改宿主中间件目录与路由注册契约。

#### Scenario: 未认证请求可访问公开分组

- **WHEN** 未携带 JWT 的请求访问`GET /api/v1/play`或`GET /h5`
- **THEN** 请求被正常处理并返回业务响应，不要求登录

#### Scenario: 宿主契约零修改

- **WHEN** 审查本次变更的宿主代码
- **THEN** `lina-core`无任何修改，公开路由全部由`backend/plugin.go`基于既有`pluginhost`契约组合
