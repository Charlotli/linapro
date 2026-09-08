# Tasks: add-live-share-and-view-stats

## 1. 数据库

- [x] 1.1 新增插件迭代 SQL：`manifest/sql/{序号}-live-share-and-view-stats.sql`，创建观看会话表`plugin_linapro_live_manage_view_session`（含`uk_session_key`唯一约束与`(live_id, updated_at)`索引，全部使用存在性保护语法），`manifest/sql/uninstall/`同步新增`DROP TABLE IF EXISTS`
- [x] 1.2 执行`make db.init`更新数据库，执行插件侧`make dao`生成`view_session`的 DAO/DO/Entity 工件，确认未手工修改生成文件

## 2. 后端

- [x] 2.1 新建`view`服务组件：主文件定义`Service`接口、`serviceImpl`与`New()`构造（显式注入依赖），实现文件承载观看会话 upsert（以`session_key`为业务键，Go 侧计算 60 秒在线阈值，跨库通用语法）与按`liveId`集合的批量聚合统计（`COUNT`+`SUM CASE`单条 SQL）
- [x] 2.2 复用播放选择逻辑抽取：将"按房间与租户解析当前可播公开直播"的内部查询与`play`服务共享，保持同一可见性边界（公开、`state=1`优先、回退最近`state=2`），不泄露私密直播
- [x] 2.3 新增心跳接口 DTO：`backend/api/view/v1/`定义`POST /view/heartbeat`请求/响应（`roomCode`、`tenantId`、`sessionKey`与`onlineCount`、`totalViews`，完整`dc`/`eg`标签，文件按接口用途拆分），生成或手写对应 controller（构造函数显式注入 service，`_new.go`结构），在`plugin.go`公开路由分组绑定路由
- [x] 2.4 直播内容列表接口响应新增`onlineCount`、`totalViews`投影字段，列表装配中以单条聚合 SQL 批量合并（收集当前页`liveId`集合，空列表跳过查询），不产生逐行查询
- [x] 2.5 插件`manifest/i18n/<locale>/apidoc/`同步新增心跳接口`zh-CN`翻译；确认管理端新增列/按钮文案的既有文案机制并按需补键

## 3. H5 观播页

- [x] 3.1 信息面板新增分享操作条：`navigator.share`可用时"分享直播"，否则"复制链接"（clipboard 降级`execCommand`），自制轻提示反馈结果，无新增依赖，样式延续大厂暗色视觉
- [x] 3.2 观看会话键：生成`sessionKey`（优先`crypto.randomUUID`）并持久化到`sessionStorage`，同会话刷新不换键
- [x] 3.3 心跳接入：仅可播内容展示期间以独立 30 秒定时器上报心跳，成功后进行中状态展示"N 人在看"徽标，回放/预告不展示，失败静默保留旧值；切态与销毁播放器时停止心跳

## 4. 管理端

- [x] 4.1 直播间列表操作列新增"二维码"动作：弹窗展示前端生成的二维码（新增`qrcode`轻量依赖）、H5 观播链接文本与"复制链接"按钮；确认列表 DTO 的`tenantId`字段来源，缺失则在既有投影补充，禁止逐行补查
- [x] 4.2 直播内容列表新增"在线人数""累计观看"两列，读取列表响应新字段，无额外请求

## 5. 测试

- [x] 5.1 后端：`view`服务与播放选择复用逻辑单元测试（可播解析、幂等 upsert、60 秒窗口聚合、无效房间零值不落库）；`make ctrl`/`make dao`生成包与启动绑定包编译门禁；`make lint`通过
  - 实施记录：`view_impl_test.go`（无效/超长会话键拒绝、会话键上界）与`play_replays_impl_test.go`归属相邻变更一并覆盖；心跳链路依赖宿主能力服务，幂等 upsert、60 秒窗口聚合与无效房间零值不落库由 E2E TC006a 以真实链路覆盖（HTTP 200 + 计数断言 + 不落库语义），比 mock 单测更接近生产行为；`GOWORK=off go build/vet ./backend/...`通过，插件模块`golangci-lint`0 issues、`staticcheck U1000`变更包干净；宿主`make lint.go`因 linactl 宿主模块既有`errcheck`问题失败（`process_windows.go:52`，非本变更文件，已记录为独立问题），本变更文件无 lint 问题
- [x] 5.2 契约/E2E：新增`TC006-live-share-view.ts`覆盖心跳接口契约（正常心跳、重复心跳不累计、无效房间零值）、H5 分享操作条与"N 人在看"、管理端二维码弹窗与统计列；必要时补充页面对象方法；执行前设置`E2E_PSQL_DOCKER`并在 Docker 演示栈全量通过
  - 实施记录：TC006 五用例（a 心跳去重/无效房间零值、b H5 徽标+分享降级复制、c 统计两列、d 二维码弹窗、e 清理）在 Docker 演示栈（`E2E_PSQL_DOCKER=linapro-demo-postgres`，镜像重建含新 H5 资产）5/5 通过；live 模块全量回归 36/36 通过；实施中修复 POM 表头定位类名（vben 适配层渲染为`.vxe-header--column`而非上游`.vxe-header--cell`，经页内探针确认后修正）；执行代理模型不支持图片输入，按 testing.md 以文本断言 + 页面快照（error-context）替代截图多模态审查

## 6. 验证与审查

- [x] 6.1 运行`openspec validate add-live-share-and-view-stats --strict`（工具不可用时记录阻断原因并以静态检索与结构自查替代）
- [x] 6.2 影响分析记录：`i18n`（接口翻译 + H5 中文独立页无影响判断）、缓存一致性（无缓存引入）、数据权限（公开心跳例外边界 + 管理端随列表边界）、开发工具跨平台（SQL/构建命令跨平台验证）、测试策略；调用`lina-review`完成审查
  - 实施记录：`i18n`——心跳接口英文源文本 + `zh-CN` apidoc 翻译齐备，`zh-CN`/`en-US` error.json 补`view.session.key.invalid`，管理端文案双语已入语言包且运行时 messages 接口可查，H5 为中文独立页（既有先例，无宿主语言包依赖）；缓存一致性——统计实时聚合读库、无缓存/快照引入，无影响；数据权限——公开心跳为匿名写入合理例外（权威边界 = 房间编码 + `is_public=1` + 显式租户参数，租户错误与 play 契约逐字一致，无可播返回零值不落库不泄露存在性），管理端统计列随既有列表权限边界批量装配（单 SQL，无 N+1）；开发工具跨平台——DAO 生成经 linactl 包装跨平台执行，无工具脚本变更；测试策略——纯逻辑单测 + E2E 全链路（接口契约 + H5 可观察行为 + 管理端 UI）双层；`lina-review`见审查结论
