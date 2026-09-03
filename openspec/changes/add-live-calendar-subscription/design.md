# 设计：直播开播日程日历订阅（ICS）

## 模块归属

- 全部变更落在插件`linapro-live-manage`内部：公开 ICS 只读接口 + H5 静态页入口 + 管理端弹窗。不触碰`apps/lina-core`，不新增插件间契约。
- 后端实现以`calendar`服务组件承载（`Service`接口 + `serviceImpl` + `New()`显式注入依赖），与`play`、`view`并列；房间解析与公开可见性判定复用既有内部共享查询（与`add-live-share-and-view-stats`中抽取方式一致），不新造窄接口。无新增抽象层或配置模型。

## 接口设计

`GET /x/linapro-live-manage/api/v1/subscribe?roomCode=&tenantId=`，挂既有公开路由分组（无`Auth`/`Tenancy`/`Permission`）。

- 响应为日历流而非 JSON：`Content-Type: text/calendar; charset=utf-8`，`Content-Disposition`携带`filename="linapro-live-<roomCode>.ics"`。这是既定 RFC 5545 文本格式的合法只读响应形态；`g.Meta`仍需完整`dc`标签参与 apidoc 治理，响应体不走`HandlerResponse`包装（日历客户端不能解析 JSON 信封），在路由绑定处与既有公开分组约定一致地按原始输出处理。
- 边界语义与`play`一致：租户插件启用时`tenantId`必填且校验存在并启用；房间不存在、私密直播、停用租户一律返回空日历（`BEGIN:VCALENDAR/END:VCALENDAR`无条目，HTTP 200），不区分原因、不泄露房间存在性——空日历对订阅客户端是合法稳态，日历应用会正常保留已同步的历史条目。
- 时间窗口：默认`now-7d`到`now+90d`（已结束直播保留 7 天供"回看日历条目直达回放"，未来 90 天覆盖季度聚会表）；窗口由 Go 侧计算后作为查询参数注入，不使用数据库方言函数。
- 数据库访问：单条投影查询（`liveDate/startTime/endTime/title/state/is_public`必要字段，`BETWEEN`窗口过滤 + 租户与房间过滤 + `deleted_at`自动过滤），访问次数恒为 1，与直播场次数线性但有界（窗口上限），无需分页。查询按`liveDate`升序稳定排序。

## ICS 输出规则

- `VCALENDAR`头：`PRODID:-//LinaPro//linapro-live-manage//CN`、`VERSION:2.0`、`CALSCALE:GREGORIAN`、`X-WR-CALNAME:<房间名称>`、`X-WR-TIMEZONE:Asia/Shanghai`。
- 每场公开直播一个`VEVENT`：
  - `UID:linapro-live-<liveId>@linapro-live-manage`——以直播记录主键为稳定键，重新订阅/多次刷新不产生重复条目。
  - `DTSTAMP`为当前生成时刻（UTC，`YYYYMMDDTHHMMSSZ`）。
  - 未开始直播：有`startTime`时`DTSTART`取其 Unix 毫秒转 UTC，`DTEND`取`endTime`；无`endTime`时按插件既有展示约定补默认时长 2 小时（常量注释说明）；无`startTime`时输出全天事件`DTSTART;VALUE=DATE:<YYYYMMDD>`（`liveDate`），不编造时间。
  - 已结束直播：取实际起止时间戳。
  - `SUMMARY`为直播标题，`DESCRIPTION`含观播页链接（`<origin>/x/linapro-live-manage/h5?room=<roomCode>`，租户启用时含`&tenant=`）与房间名称。
  - 文本转义遵守 RFC 5545：`\`、`;`、`,`与换行转义，标题/房间名经统一转义函数处理。
- 全部时间输出为 UTC（`Z`后缀），由 Go 侧完成本地时间到 UTC 转换，不依赖客户端时区。

## 前端设计

### H5 预告态入口

- 未开始（预告）态信息面板新增"添加到日历"按钮（`.share-bar`同级，样式延续大厂暗色渐变胶囊）：点击优先尝试`window.open`订阅链接（系统日历应用可识别`webcal`场景由用户环境决定），并提供"复制订阅链接"的次级行为（复制`subscribe`接口 URL）；成功失败均走既有轻提示。进行中/回放/无直播态不展示。

### 管理端订阅链接入口

- 直播间列表操作列新增"日历订阅"动作（与二维码动作并列），弹窗复用二维码弹窗形态：展示订阅链接全文、"复制链接"按钮与一句说明（"在手机日历中选择'订阅日历'并粘贴此链接"）。无新增依赖。

## i18n 影响

- 插件`i18n.enabled: true`：订阅接口的`summary`/`dc`等文档元数据以英文源文本编写，`manifest/i18n/<locale>/apidoc/`同步`zh-CN`翻译。
- H5 静态页为中文独立页（既有先例），新增文案跟随既有做法；管理端新增文案按插件既有前端文案机制维护。

## 数据权限例外说明

- 公开订阅接口为匿名读取的合理例外：权威边界与`play`接口完全一致（房间存在 + 公开直播 + 显式租户参数），输出仅为日程投影（标题/时间/房间名/观播链接），不包含推流地址、管理端字段或非公开直播条目；无效输入统一空日历。对应验证由契约测试覆盖。

## 复杂度判断

- 单一职责只读组件 + 既有页面局部扩展，局部直接实现即可；ICS 转义/时间格式化为纯函数便于单元测试，不引入第三方 ICS 库（输出面小、规则明确，自实现约百余行且无运行时依赖）。
