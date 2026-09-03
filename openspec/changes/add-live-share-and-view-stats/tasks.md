# Tasks: add-live-share-and-view-stats

## 1. 数据库

- [ ] 1.1 新增插件迭代 SQL：`manifest/sql/{序号}-live-share-and-view-stats.sql`，创建观看会话表`plugin_linapro_live_manage_view_session`（含`uk_session_key`唯一约束与`(live_id, updated_at)`索引，全部使用存在性保护语法），`manifest/sql/uninstall/`同步新增`DROP TABLE IF EXISTS`
- [ ] 1.2 执行`make db.init`更新数据库，执行插件侧`make dao`生成`view_session`的 DAO/DO/Entity 工件，确认未手工修改生成文件

## 2. 后端

- [ ] 2.1 新建`view`服务组件：主文件定义`Service`接口、`serviceImpl`与`New()`构造（显式注入依赖），实现文件承载观看会话 upsert（以`session_key`为业务键，Go 侧计算 60 秒在线阈值，跨库通用语法）与按`liveId`集合的批量聚合统计（`COUNT`+`SUM CASE`单条 SQL）
- [ ] 2.2 复用播放选择逻辑抽取：将"按房间与租户解析当前可播公开直播"的内部查询与`play`服务共享，保持同一可见性边界（公开、`state=1`优先、回退最近`state=2`），不泄露私密直播
- [ ] 2.3 新增心跳接口 DTO：`backend/api/view/v1/`定义`POST /view/heartbeat`请求/响应（`roomCode`、`tenantId`、`sessionKey`与`onlineCount`、`totalViews`，完整`dc`/`eg`标签，文件按接口用途拆分），生成或手写对应 controller（构造函数显式注入 service，`_new.go`结构），在`plugin.go`公开路由分组绑定路由
- [ ] 2.4 直播内容列表接口响应新增`onlineCount`、`totalViews`投影字段，列表装配中以单条聚合 SQL 批量合并（收集当前页`liveId`集合，空列表跳过查询），不产生逐行查询
- [ ] 2.5 插件`manifest/i18n/<locale>/apidoc/`同步新增心跳接口`zh-CN`翻译；确认管理端新增列/按钮文案的既有文案机制并按需补键

## 3. H5 观播页

- [ ] 3.1 信息面板新增分享操作条：`navigator.share`可用时"分享直播"，否则"复制链接"（clipboard 降级`execCommand`），自制轻提示反馈结果，无新增依赖，样式延续大厂暗色视觉
- [ ] 3.2 观看会话键：生成`sessionKey`（优先`crypto.randomUUID`）并持久化到`sessionStorage`，同会话刷新不换键
- [ ] 3.3 心跳接入：仅可播内容展示期间以独立 30 秒定时器上报心跳，成功后进行中状态展示"N 人在看"徽标，回放/预告不展示，失败静默保留旧值；切态与销毁播放器时停止心跳

## 4. 管理端

- [ ] 4.1 直播间列表操作列新增"二维码"动作：弹窗展示前端生成的二维码（新增`qrcode`轻量依赖）、H5 观播链接文本与"复制链接"按钮；确认列表 DTO 的`tenantId`字段来源，缺失则在既有投影补充，禁止逐行补查
- [ ] 4.2 直播内容列表新增"在线人数""累计观看"两列，读取列表响应新字段，无额外请求

## 5. 测试

- [ ] 5.1 后端：`view`服务与播放选择复用逻辑单元测试（可播解析、幂等 upsert、60 秒窗口聚合、无效房间零值不落库）；`make ctrl`/`make dao`生成包与启动绑定包编译门禁；`make lint`通过
- [ ] 5.2 契约/E2E：新增`TC005-live-share-view.ts`覆盖心跳接口契约（正常心跳、重复心跳不累计、无效房间零值）、H5 分享操作条与"N 人在看"、管理端二维码弹窗与统计列；必要时补充页面对象方法；执行前设置`E2E_PSQL_DOCKER`并在 Docker 演示栈全量通过

## 6. 验证与审查

- [ ] 6.1 运行`openspec validate add-live-share-and-view-stats --strict`（工具不可用时记录阻断原因并以静态检索与结构自查替代）
- [ ] 6.2 影响分析记录：`i18n`（接口翻译 + H5 中文独立页无影响判断）、缓存一致性（无缓存引入）、数据权限（公开心跳例外边界 + 管理端随列表边界）、开发工具跨平台（SQL/构建命令跨平台验证）、测试策略；调用`lina-review`完成审查
