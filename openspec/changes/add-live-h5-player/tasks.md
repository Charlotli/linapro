# add-live-h5-player 实施任务

## 1. 后端：公开播放信息接口

- [x] 1.1 新增`backend/api/play/v1/play.go`：定义`PlayItem`共享响应 DTO（最小投影，含 Unix 毫秒时间字段与`date-only`日期说明）。
- [x] 1.2 新增`backend/api/play/v1/play_get.go`：定义`GET /play`请求/响应 DTO（`roomCode`、可空`tenantId`参数与完整`dc`/`eg`标签）。
- [x] 1.3 新增`backend/api/play/play.go`接口定义文件（`gf` CLI 不可用，按生成格式手工维护，文件头已注明；task 记录该例外）。
- [x] 1.4 新增`backend/internal/service/play/play.go`：组件主文件（`Service`接口、`serviceImpl`、`New()`构造、输入/输出类型），显式注入`tenantcap.Service`。
- [x] 1.5 新增`backend/internal/service/play/play_code.go`：定义`LIVE_MANAGE_PLAY_NOT_FOUND`、`LIVE_MANAGE_PLAY_TENANT_REQUIRED`、`LIVE_MANAGE_PLAY_TENANT_INVALID`错误码。
- [x] 1.6 新增`backend/internal/service/play/play_impl.go`：实现租户校验（`Directory().Get`）、房间定点查询、逐态定点探测（进行中 > 回放 > 预告，至多 3 次）与最小投影装配，进行中态校验直播间为`直播中`。
- [x] 1.7 新增`backend/internal/controller/play/play_new.go`与`play_v1_get.go`：控制器构造（显式注入服务）与`Get`处理方法。

## 2. 路由注册

- [x] 2.1 修改`backend/plugin.go`：新增公开 JSON 分组（`NeverDoneCtx`+`HandlerResponse`+`CORS`+`RequestBodyLimit`+`Ctx`）绑定播放控制器；新增公开静态分组（`NeverDoneCtx`+`CORS`+`Ctx`）提供`GET /h5`与`GET /h5/*any`嵌入式静态文件服务（路径清理 + `mime.TypeByExtension`判定 Content-Type + 非缺失错误日志）。
- [x] 2.2 记录 DI 来源检查：播放服务依赖`tenantcap.Service`来自注册回调的宿主`HostServices`目录；静态资产由处理器闭包直接读取插件`EmbeddedFiles`；无新增共享实例。

## 3. 前端：H5 播放页

- [x] 3.1 新增`frontend/h5/index.html`：移动端优先的播放页骨架（播放器容器、封面区、信息区、状态提示区）。
- [x] 3.2 新增`frontend/h5/hls.light.min.js`：内嵌`hls.js`v1.5.20 轻量发行版（jsdelivr 下载，约 297KB，不依赖外网 CDN）。
- [x] 3.3 新增`frontend/h5/app.js`：参数解析、播放信息拉取、`hls.js`/原生 HLS 播放降级、直播中/回放/预告/无直播状态呈现、30 秒轮询、节目信息渲染（歌单 JSON 解析失败按空处理）、推流地址零暴露。
- [x] 3.4 新增`frontend/h5/app.css`：移动端样式（16:9 播放器、暗色主题、安全区适配、全离线系统字体栈）。

## 4. i18n 资源

- [x] 4.1 更新`manifest/i18n/zh-CN/apidoc/plugin-api-main.json`：补齐`play/v1`接口翻译（`GetReq`meta/fields、`PlayItem`fields）。
- [x] 4.2 确认`manifest/i18n/en-US/apidoc/plugin-api-main.json`占位合规（保持空对象占位，未新增映射）。
- [x] 4.3 更新`manifest/i18n/zh-CN/error.json`与`manifest/i18n/en-US/error.json`：新增 3 个播放错误码翻译（`error.live.manage.play.*`）。

## 5. 测试

- [x] 5.1 新增`backend/internal/service/play/play_impl_test.go`：覆盖错误码元数据、播放地址可见性、房间状态门控、状态优先级、投影排除字段与租户解析（stub 租户服务）的单元测试（自包含、无 DB）。
- [x] 5.2 扩展`backend/api/api_contract_test.go`：新增`TestPlayResponseDTOsHideAdministrativeFields`（响应 DTO 不含`pushUrl`/`pageUrl`等管理端字段）。
- [x] 5.3 新增`apps/lina-plugins/linapro-live-manage/hack/tests/e2e/live/TC004-live-h5-play.ts`与`hack/tests/pages/LivePlayH5Page.ts`：覆盖未开始预告、开启直播后播放器呈现、接口不泄露管理端字段、回放态、不存在直播间统一`NOT_FOUND`提示与清理。
- [x] 5.4 Go 编译门禁：`go build ./...`（插件模块，`GOWORK=off`）与`go test ./... -count=1`全部通过。
- [x] 5.5 Docker 演示栈（`hack/deploy/docker-compose.yaml` + 本地构建镜像）端到端冒烟：`/h5`、`/h5/`、`/h5/*` 全部 200；播放接口未开始（无播放地址）/进行中（含播放地址）/已结束（回放）/不存在房间统一`NOT_FOUND`/缺租户参数`TENANT_REQUIRED`全部符合预期；H5 页面标题与脚本资源加载正常。（Playwright TC004 已在 Docker 演示栈执行通过，见 FB-1/FB-2 验证记录）

## 6. 文档

- [x] 6.1 更新插件`README.md`与`README.zh-CN.md`：新增观众端播放页能力说明与分享链接格式（中英文镜像同步）。
- [x] 6.2 治理验证：`openspec validate`不可用已记录（CLI 未安装）；静态检索确认无`pushUrl`外泄（契约测试 + 单测断言）；公开路由分组中间件核对完成；文档镜像一致性检查通过。

## 7. 审查

- [x] 7.1 完成`lina-review`代码与规范审查（数据权限例外、`i18n`影响、性能边界、DI 来源、跨平台影响判断齐备）。
  - 审查记录：范围=`git status --short` + `git ls-files --others --exclude-standard`（插件子仓库已展开），覆盖反馈修复 FB-1~FB-4 文件与变更产物；已读取规则文件=`AGENTS.md`、`backend-go.md`、`frontend-ui.md`、`testing.md`、`i18n.md`、`dev-tooling.md`、`api-contract.md`、`plugin.md`、`openspec.md`（architecture/data-permission/cache-consistency/database 未命中，无模块边界、schema、缓存与数据权限路径变更）。插件根无本地`AGENTS.md`（已记录不存在判断）。发现 0 严重、2 警告（TC004 链式数据依赖为模块既有 TC 断面模式，记录改进建议；TC004 基线依赖 psql/docker 助手可用性，已在任务记录说明环境要求）。多模态截图审查因当前模型无图像输入按 testing.md 记录跳过，以文本断言替代。验证证据：插件`go build ./...`、`go test ./backend/internal/service/... -count=1`、E2E 单文件`6 passed`与全量`20 passed (2.8m)`；`openspec validate --strict`因 CLI 未安装不可用（已记录，静态核对 tasks/specs 格式合规）。

## Feedback

- [x] **FB-1**: 裸`/h5`入口下`index.html`相对路径资源（`app.css`/`app.js`）解析到`/x/linapro-live-manage/*`导致 404，页面无样式且状态文案堆叠；改用宿主强制命名空间约定下的绝对路径引用，并补充 E2E 子断言防回归。
  - 验证：`index.html`资源引用改为`/x/linapro-live-manage/h5/app.css`、`/x/linapro-live-manage/h5/hls.light.min.js`、`/x/linapro-live-manage/h5/app.js`；同时修复`plugin.go`裸`/h5`（空`any`通配参数）回退到`h5IndexFile`；E2E 新增`expectStylesApplied()`（计算样式证明 CSS+JS 生效）与`gotoBare()`裸入口回归块；TC004 全量通过。
- [x] **FB-2**: H5 播放页视觉质量不达标；按"礼堂夜幕"暗色视觉方向重做页面结构、样式与状态呈现，并同步更新 E2E 页面对象选择器。
  - 验证：`index.html`重构为 loading 骨架/播放器/封面/信息/提示/页脚分区（各状态显式分区 id），`app.css`重做为墨色分层暗色主题（金色强调、衬线展示字体、骨架微光/脉冲动画、安全区适配、680px 阅读宽度），`app.js`改为状态机渲染（直播中→播放器+徽章、已结束→回放徽章、未开始→预告+30s 轮询、无直播/租户异常/网络异常分态提示）；`LivePlayH5Page`选择器同步（`#state-chip`等）；TC004 在 Docker 演示栈全部通过（`20 passed`）。
  - 二次迭代（用户反馈"大厂风格"）：按头部直播平台观播页语言重做为沉浸式深色视觉——全宽贴顶播放器舞台（仅底部大圆角+底缘渐变过渡）、上浮叠压毛玻璃信息面板（负 margin 叠压+backdrop blur+顶部高光描边）、实色渐变状态胶囊（直播中红橙/回放蓝紫）、渐变竖条节目卡头、内联 SVG 提示图标替代 emoji、三层呼吸光环品牌封面、错峰入场动画、环境光晕装饰；修复光晕导致的移动端横向溢出（`html overflow-x: hidden`）。全部 E2E 选择器 id 保持兼容，Docker 演示栈重建后视觉冒烟（资产 200、无横向滚动、状态渲染正确）+ 全量 E2E `20 passed (2.7m)` 通过。
  - 三次迭代（用户反馈"继续优化"）：播放健壮性升级——(1) hls.js 分级恢复（网络错误 `startLoad` 续流、媒体错误 `recoverMediaError`，上限 2 次）；(2) 修复恢复盲区：`startLoad` 对"源彻底不可达"静默失败不再抛 fatal，新增 10s 恢复监视（`FRAG_LOADED` 清除，`readyState` 无进展降级为错误胶囊+重试按钮，本地坏流源实测命中）；(3) 自动播放被拦截时静音自动播放 + "点击开启声音"提示按钮；(4) 播放失败胶囊内嵌"重试"按钮重新拉流；(5) fetch 10s 硬超时（AbortController）防弱网挂起；(6) 回前台时轮询态立即刷新、播放态不打断。信息呈现升级——直播间名称独立眉标（渐变光点）、预告封面新增开播日期胶囊、视频补 X5/微信内核属性、移除死代码 `hideAll`。TC004a 增加 `#live-room`/`#cover-date` 断言；Docker 演示栈重建后坏流实测（错误胶囊+重试可见）+ 全量 E2E `20 passed (2.6m)` 通过。
- [x] **FB-3**: E2E 数据库基线助手（`hack/tests/support/postgres.ts`）依赖宿主机本地`psql`且默认路径硬编码 macOS，Windows 宿主无法执行（`spawnSync psql ENOENT`）；新增`E2E_PSQL_DOCKER`模式经`docker exec`调用容器内 psql，保持跨平台可执行。
  - 验证：`runPsql`分支按`E2E_PSQL_DOCKER`路由`docker exec -i -e PG*… <container> psql`（Node 22 `spawnSync`拒绝`.cmd` shim，故直连 docker CLI）；`psqlEnv()`在 docker 模式固定`PGHOST=127.0.0.1`；全量 live 模块 E2E 在 Windows 宿主 + Docker 演示栈下运行通过（`E2E_PSQL_DOCKER=linapro-demo-postgres`，`20 passed`）。
- [x] **FB-4**: TC004 修复过程中连带发现并修复的两处问题：(1) 后端`liveroom_impl.go`的`firstReferencedRoom`对空结果集`Scan`到非指针结构体导致"删除未被引用直播间"必然返回`sql: no rows in result set`（TC001e 删除失败），改为切片投影容忍空结果；(2) E2E 页面对象适配 vxe-table 固定操作列（操作行与数据行按`rowid`关联）与管理后台 hash 路由（`/admin#/live/room`），以及 TC004 改为确定性 stamp + `beforeAll`硬删除基线，规避 Playwright 用例失败回收 worker 后重新 import 测试文件导致`Date.now()`常量漂移、共享测试数据失联的问题。
  - 验证：插件`go build ./...`与`go test ./backend/internal/service/... -count=1`通过；重建镜像（`make.cmd image plugins`）后 Docker 演示栈删除接口返回成功提示；全量 live 模块 E2E `20 passed (2.8m)`。
