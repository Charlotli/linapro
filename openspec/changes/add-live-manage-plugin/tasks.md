## 1. 插件骨架与数据契约

- [x] 1.1 创建`apps/lina-plugins/linapro-live-manage/`插件骨架：`plugin.yaml`（菜单、权限点、`i18n`、租户与分发治理声明）、`plugin_embed.go`、`go.mod`、`Makefile`、`hack/config.yaml`、`backend/plugin.go`路由注册。
- [x] 1.2 编写`manifest/sql/001-linapro-live-manage-schema.sql`（幂等建表、索引、字典 seed）、`manifest/sql/uninstall/001-linapro-live-manage-schema.sql`与`manifest/sql/mock-data/001-linapro-live-manage-mock-data.sql`。
- [x] 1.3 编写`backend/api/liveroom`与`backend/api/live`的`RESTful DTO`（`GET`查询、`POST`创建、`PUT`更新、`DELETE`删除；时间点字段`Unix`毫秒、`dc/eg/permission`标签完整、英文源文本）。
- [x] 1.4 创建插件数据表并执行`make dao`生成`DAO/DO/Entity`，执行`make ctrl`生成控制器骨架。

## 2. 后端实现

- [x] 2.1 实现直播间服务层：租户过滤列表/详情/候选/创建/更新/删除，删除引用校验，`bizerr`错误码，创建人名称批量装配。
- [x] 2.2 实现直播内容服务层：租户过滤`CRUD`、直播间存在与未禁用校验、歌单`JSON`校验、列表直播间名称批量装配、`bizerr`错误码。
- [x] 2.3 实现控制器层并完成路由绑定，运行插件启动绑定与`API`契约测试。

## 3. 前端实现

- [x] 3.1 实现直播间管理页（查询表单、表格、字典标签、新增/编辑弹窗、删除确认、候选接口）。
- [x] 3.2 实现直播内容管理页（查询表单、表格、直播间候选联动、歌单动态行编辑、新增/编辑弹窗、删除确认）。

## 4. 资源与文档

- [x] 4.1 维护`manifest/i18n/zh-CN`与`manifest/i18n/en-US`的菜单、错误、插件文案资源与`apidoc`翻译资源（`en-US`空占位）。
- [x] 4.2 维护插件`README.md`与`README.zh-CN.md`双语镜像及`manifest/docs`文档。

## 5. 测试与验证

- [x] 5.1 新增`hack/tests/e2e/live/`下直播间与直播内容`CRUD`的`E2E`用例及页面对象（使用`lina-e2e`技能规范）。
- [x] 5.2 运行覆盖变更包的`go test`编译门禁与`make plugins.check`治理检查。
- [x] 5.3 运行`openspec validate add-live-manage-plugin --strict`；工具未安装时记录阻断原因。

## 6. 直播开启/关闭迭代

- [x] 6.1 新增`PUT /live/{id}/start`与`PUT /live/{id}/stop`DTO、服务层状态机（`未开始`→`进行中`→`已结束`）、事务内直播间状态联动、同房间占用校验与新 bizerr 错误码。
  - 验证记录：`liveStateTransition`纯函数状态机（仅`未开始`可开启、仅`进行中`可关闭）由`TestLiveStateTransition`六分支覆盖；开启/关闭在`dao.Live.Transaction`闭包内同步直播与直播间状态；开启前置校验直播间可用（复用`ensureRoomAvailable`）与同房间无其他进行中直播（`ensureRoomFreeForLive`，新错误码`LIVE_MANAGE_ROOM_BUSY`）；非法迁移返回`LIVE_MANAGE_LIVE_STATE_TRANSITION`；全部查询注入租户过滤，跨租户`ID`表现为不存在。
- [x] 6.2 `plugin.yaml`新增`live:live:start`/`live:live:stop`按钮权限点，同步双语菜单、错误与`apidoc`翻译资源。
- [x] 6.3 直播内容列表新增条件动作按钮（`未开始`显示开启、`进行中`显示关闭）与确认交互。
- [x] 6.4 新增状态迁移单测与`TC003`直播开启/关闭`E2E`用例，运行编译门禁与治理检查。
  - 验证记录：`GOWORK=off go build/go test`通过；`gofmt`干净；SFC 编译校验通过；`linactl plugins.check`484 个文件 0 发现；`linactl i18n.check`三项通过。`TC003-live-start-stop.ts`（TC-3a~d）与`LiveContentPage.startLive/stopLive`页面对象已创建，端到端执行依赖运行环境（同任务 5.1 记录）。

## 7. 优化加固迭代

- [x] 7.1 直播间绑定不可变：更新接口拒绝换绑直播间（新错误码`LIVE_MANAGE_ROOM_IMMUTABLE`），消除换绑绕过占用校验、房间状态联动失联的一致性漏洞；前端编辑弹窗与`apidoc`契约同步。
- [x] 7.2 数据库层并发防护：新增部分唯一索引`uk_plugin_linapro_live_manage_live_tenant_room_ongoing`（`tenant_id+room_id`，`state=1`且未软删除），并发开启无法再绕过服务端占用校验；SQL 幂等并已应用至运行库。
- [x] 7.3 死代码清理：移除 live 服务未使用的`timeToMillis`/`formatLiveDate`与两个服务未使用的`pluginID`常量。
- [x] 7.4 构建工具修复：`linactl build`新增插件前端内容戳（`internal/frontend/frontend_stamp.go`），插件前端/清单变更时通过`TURBO_FORCE`强制前端重建，修复插件页面改动被 turbo 缓存吞掉的缺陷；遵循`dev-tooling.md`以 Go 实现、跨平台。
  - 验证记录：`linactl`编译通过；带戳构建触发`forcing turbo rebuild`且产物包含开启/关闭按钮逻辑（`bootstrap-DgFfUiuq.js`含`liveStart`/`startConfirm`）；stamp 文件成功写入`temp/frontend-plugin-pages.stamp`。
  - 阻断记录：`linactl`自带测试`TestRunBuildDirBuildsSelectedPluginOnly`与`TestResolvePluginConfigBuildStepsSkipsPluginsWithoutCommands`在`Windows`下失败（路径断言），经验证在不含本次改动的基线上同样失败，属既有平台问题，不在本变更处理。
- [x] 7.5 遗留记录：`linapro-demo-dynamic`测试文件引用已被重构移除的能力类型（`SetSysConfigValueOptions`、`DirectPutInput`等，存储分片上传协议与`hostconfig`选项已重设计），需按新能力`SDK`独立移植后方可解锁`lint.go plugins=1`；不阻塞构建与运行。

## 8. 影响记录与审查

- [x] 6.1 记录影响判断：`i18n`（插件多语言资源）、缓存一致性（无缓存状态）、数据权限（租户过滤+组织维度例外说明）、开发工具跨平台（复用共享`codegen`入口）、测试策略（契约测试/单元测试/`E2E`）。
- [x] 6.2 完成实现和验证后调用`lina-review`进行代码和规范审查，审查通过后再标记任务完成。

## 任务执行记录

### 验证证据

- 插件独立模块构建`GOWORK=off go build ./...`通过；`GOWORK=off go test ./... -count=1`全部通过（含`backend/api`契约测试、`live`与`liveroom`服务测试、候选条件构造测试）。
- `gofmt -l`格式化后0文件残留；`go vet ./backend/...`通过。
- `make plugins.check`通过：扫描503个文件，0项治理发现。
- `make lint`阻断：`linactl lint.go plugins=auto`因`linapro-ai-core`插件与`lina-core`模块不同步（`capregistry`包无法解析）在`ai-core`处失败，属与本变更无关的既有问题；已用`gofmt`+`go vet`+全量`go test`替代覆盖本插件模块，`golangci-lint`（v2.12.2）专项检查待环境修复后补跑。
- `openspec validate add-live-manage-plugin --strict`阻断：`openspec CLI`在本环境未安装（全局与本地`npm`均无，`npx`无法解析该包，`linactl`无对应子命令），按规则记录为阻断原因，待工具可用后补跑。
- `E2E`用例（`TC001`、`TC002`）及页面对象已按`lina-e2e`规范落位插件`hack/tests/e2e/`与`hack/tests/pages/`目录；本环境缺少数据库配置（`manifest/config/config.yaml`缺失）与前端依赖，未执行端到端运行。

### 实现期修复记录

- 修复直播间候选接口租户隔离缺陷：`Options`原实现将关键词`WhereOrLike`直接链接在租户过滤之后，关键词非空时`OR`条件逃逸租户谓词，导致跨租户直播间泄漏且禁用过滤失效。现改为`buildOptionsCondition`以`WhereBuilder`子构建器分组关键词条件后再挂载模型，并新增`liveroom_impl_test.go`以`Build()`纯渲染断言分组语义（不依赖数据库）。
- 对齐设计契约：候选接口按设计补充`limit`参数与默认上限100条、硬上限200条（`clampOptionsLimit`）的语义，`OptionsReq`、控制器映射、服务层输入同步更新。
- `zh-CN apidoc`翻译资源修正：新增`limit`字段翻译并更新候选接口元数据描述；两个删除接口的`ids`字段描述由"多个用逗号分隔"修正为查询数组契约（`ids[]=1&ids[]=2`）。
- `TC002`补充创建直播间的清理步骤（`deleteRoomIfExists`），避免测试残留数据。
- `gofmt -w`统一格式化5个不符合格式规范的Go源码文件（含本次编辑文件与既有生成文件）。

### 审查结论（6.2）

`lina-review`审查通过：严重问题0项，警告4项（详见审查报告）；实现期发现的候选接口租户隔离缺陷已修复并以单元测试固化，任务6.2据此标记完成。

### 影响判断（6.1）

- `i18n`：插件`plugin.yaml`显式启用`i18n`，`manifest/i18n/zh-CN`与`en-US`的菜单、错误、插件文案资源`key`一一对应，`en-US apidoc`保持空占位；`zh-CN apidoc`翻译随本次接口变更同步更新，无遗漏翻译项。
- 缓存一致性：本插件不持有缓存、派生状态或快照，无影响。
- 数据权限：读取与写入均在数据库查询阶段经`tenantspi.ApplyPluginTableFilter`注入租户过滤，写操作按`tenant_id`+`id`双条件更新；候选接口已在实现期修复租户隔离缺陷并以单元测试固化。组织维度例外沿用设计说明：直播间与直播内容无组织/部门归属语义，权威边界为租户隔离，不接入组织数据权限过滤，同租户内按权限点控制读写。
- 开发工具跨平台：复用共享`codegen`与治理入口（`make`/`make.cmd`），插件内`Makefile`与宿主同构，未新增平台专属脚本。
- 测试策略：静态结构契约测试（`api_contract_test.go`）、服务层单元测试（错误码元数据、枚举校验、候选条件分组）、`E2E`用例（`TC001`/`TC002`）三层覆盖；`E2E`未在本会话执行运行验证。
- `DI`来源检查：无新增运行期依赖；服务构造函数逐项注入`bizctxcap.Service`、`tenantcap.Service`、`usercap.Service`，来源为插件注册回调的宿主`HostServices`（`registrar.Services()`），复用启动期共享实例，与设计"依赖注入"一节一致。
