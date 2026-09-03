# Tasks: add-live-calendar-subscription

## 1. 后端

- [x] 1.1 新建`calendar`服务组件（主文件`Service`接口 + `serviceImpl` + `New()`显式注入依赖）：实现窗口内公开直播日程查询（Go 侧计算时间窗口参数，投影必要字段，单条 SQL 有界返回）与 ICS 文本组装（`UID`稳定键、全天事件、UTC 转换、RFC 5545 转义为纯函数）；复用房间解析与公开可见性共享查询
  - 实施记录：房间解析与租户校验按插件既有先例（`play` 包声明独立于管理端服务）以包内私有函数实现，与`add-live-share-and-view-stats`设计中的"同包内部函数共享、不新造窄接口"口径一致，未新建共享包；租户校验失败、房间不存在、窗口内无公开直播统一降级为空日历（HTTP 200），不泄露存在性；`X-WR-CALNAME`仅在存在条目时输出，保证各种空日历逐字节一致；schema 无结束时间列，DTEND 按设计补默认时长 2 小时（常量注释说明）
- [x] 1.2 新增订阅路由：`GET /api/v1/subscribe`原始输出（`text/calendar; charset=utf-8` + `Content-Disposition`文件名，不走`HandlerResponse`包装，与公开分组既有约定一致），DTO 完整`dc`/`eg`文档标签；实施时确认直播表按房间查询的索引覆盖，不足则在该迭代 SQL 文件中补幂等索引
  - 实施记录：订阅控制器绑定在公开`/api/v1`组内、`HandlerResponse`包装之外的独立绑定（与设计一致）；窗口查询命中既有索引`idx_plugin_linapro_live_manage_live_tenant_room`（租户+房间等值前缀，房间内行数有界），无需补索引，未新增迭代 SQL；`gf` CLI 不可用，`api/subscribe/subscribe.go`接口文件沿用`add-live-h5-player`先例手工维护生成格式
- [x] 1.3 插件`manifest/i18n/<locale>/apidoc/`同步新增订阅接口`zh-CN`翻译（英文源文本）
  - 实施记录：`zh-CN/apidoc/plugin-api-main.json`新增`subscribe.v1`块；`en-US`按治理规则保持空对象占位；JSON 校验通过

## 2. 前端

- [x] 2.1 H5 预告态新增"添加到日历"入口：订阅链接打开 + 复制订阅链接 + 既有轻提示，样式延续大厂暗色视觉；非预告态不渲染
  - 实施记录：设计引用的`.share-bar`与轻提示属于未实施的`add-live-share-and-view-stats`变更，本变更按设计意图独立实现：信息面板内新增日历操作条（"添加到日历"打开订阅链接 + "复制订阅链接"次级行为）并补充轻量 toast 反馈；`i18n:check`与`vue-tsc`通过
- [x] 2.2 管理端直播间列表操作列新增"日历订阅"动作：弹窗展示订阅链接全文、复制按钮与订阅方式说明，复用二维码弹窗交互形态，无新增依赖
  - 实施记录：二维码弹窗属于未实施的`add-live-share-and-view-stats`变更，本变更按既有`useVbenModal`弹窗形态独立实现（`live-room-calendar-modal.vue`），订阅链接与公开接口同源同参（租户启用时携带当前租户 ID）；文案键同步维护`zh-CN`/`en-US`插件语言包

## 3. 测试

- [x] 3.1 后端单元测试：ICS 组装纯函数（窗口过滤、全天事件、时间戳条目、UID 稳定、标题转义、空日历）、窗口查询参数注入；`go test`编译门禁与`make lint`
  - 实施记录：`calendar_ics_test.go`（窗口计算、文档结构、时间戳/全天条目、转义、空日历一致性、UID 稳定）与`calendar_impl_test.go`（租户解析全拒绝路径降级空日历、feed 投影）全部通过；`GOWORK=off go test ./backend/...`全绿；`lint.go`对插件模块仅报`live_impl.go Update`圈复杂度一项预先存在问题（独立反馈，不属本变更）
- [x] 3.2 契约/E2E：新增`TC006-live-calendar.ts`覆盖订阅接口（`text/calendar`响应头、条目与稳定 UID、空日历不泄露存在性）、H5 预告态入口与复制行为、管理端弹窗；执行前设置`E2E_PSQL_DOCKER`并在 Docker 演示栈全量通过
  - 实施记录：按 E2E 编号连续递增规则落为`TC005-live-calendar.ts`（原 TC006 跳号，与既有 TC001-TC004 冲突）；api 契约测试新增订阅响应原始日历形态断言；Docker 演示栈（`E2E_PSQL_DOCKER=linapro-demo-postgres`）TC005 六用例与 TC004 六用例（回归）12/12 通过；多模态截图审查确认弹窗文案为翻译文本且链接完整渲染；测试数据通过 SQL 将界面未填日期的直播锚定到当天（界面新建直播日期留空会落库零值 2006-01-02，属插件既有数据质量问题，已记录为独立反馈）

## 4. 验证与审查

- [x] 4.1 运行`openspec validate add-live-calendar-subscription --strict`（工具不可用时记录阻断原因并以静态检索与结构自查替代）
  - 实施记录：`openspec` CLI 未安装（`command not found`），已记录阻断原因；静态自查通过：4 产物齐全、3 个 Requirement、6 个 Scenario、9 个任务项、spec 结构符合统一格式
- [x] 4.2 影响分析记录：`i18n`（接口翻译 + H5 中文独立页无影响判断）、缓存一致性（无缓存引入）、数据权限（匿名公开只读例外边界与`play`一致）、开发工具跨平台（无工具变更则记录无影响）、测试策略；调用`lina-review`完成审查
  - 实施记录：`i18n`——订阅接口英文源文本 + `zh-CN` apidoc 翻译已补，H5 为中文独立页（既有先例，无宿主语言包依赖），管理端文案双语已补，`i18n:check`通过；缓存一致性——订阅接口实时读库、无缓存/快照引入，无影响；数据权限——匿名公开只读例外，权威边界（房间存在 + `is_public=1` + 显式租户参数）与`play`一致，无效输入统一空日历不泄露存在性，由 TC005e 覆盖；开发工具跨平台——无工具脚本变更，无影响（`validate-e2e.mjs`内部`spawnSync('pnpm')`在 Windows 需 shell 解析的既有限制已记录）；测试策略——纯函数单测 + 契约断言 + E2E 全链路；`DI`来源检查——`calendar`服务仅注入`tenantSvc`，由路由注册处`services.Tenant()`与`play`共享同一启动期实例，未新增运行期依赖；`lina-review`审查通过（见审查结论）
