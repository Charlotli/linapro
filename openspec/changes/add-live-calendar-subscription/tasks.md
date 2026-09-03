# Tasks: add-live-calendar-subscription

## 1. 后端

- [ ] 1.1 新建`calendar`服务组件（主文件`Service`接口 + `serviceImpl` + `New()`显式注入依赖）：实现窗口内公开直播日程查询（Go 侧计算时间窗口参数，投影必要字段，单条 SQL 有界返回）与 ICS 文本组装（`UID`稳定键、全天事件、UTC 转换、RFC 5545 转义为纯函数）；复用房间解析与公开可见性共享查询
- [ ] 1.2 新增订阅路由：`GET /api/v1/subscribe`原始输出（`text/calendar; charset=utf-8` + `Content-Disposition`文件名，不走`HandlerResponse`包装，与公开分组既有约定一致），DTO 完整`dc`/`eg`文档标签；实施时确认直播表按房间查询的索引覆盖，不足则在该迭代 SQL 文件中补幂等索引
- [ ] 1.3 插件`manifest/i18n/<locale>/apidoc/`同步新增订阅接口`zh-CN`翻译（英文源文本）

## 2. 前端

- [ ] 2.1 H5 预告态新增"添加到日历"入口：订阅链接打开 + 复制订阅链接 + 既有轻提示，样式延续大厂暗色视觉；非预告态不渲染
- [ ] 2.2 管理端直播间列表操作列新增"日历订阅"动作：弹窗展示订阅链接全文、复制按钮与订阅方式说明，复用二维码弹窗交互形态，无新增依赖

## 3. 测试

- [ ] 3.1 后端单元测试：ICS 组装纯函数（窗口过滤、全天事件、时间戳条目、UID 稳定、标题转义、空日历）、窗口查询参数注入；`go test`编译门禁与`make lint`
- [ ] 3.2 契约/E2E：新增`TC006-live-calendar.ts`覆盖订阅接口（`text/calendar`响应头、条目与稳定 UID、空日历不泄露存在性）、H5 预告态入口与复制行为、管理端弹窗；执行前设置`E2E_PSQL_DOCKER`并在 Docker 演示栈全量通过

## 4. 验证与审查

- [ ] 4.1 运行`openspec validate add-live-calendar-subscription --strict`（工具不可用时记录阻断原因并以静态检索与结构自查替代）
- [ ] 4.2 影响分析记录：`i18n`（接口翻译 + H5 中文独立页无影响判断）、缓存一致性（无缓存引入）、数据权限（匿名公开只读例外边界与`play`一致）、开发工具跨平台（无工具变更则记录无影响）、测试策略；调用`lina-review`完成审查
