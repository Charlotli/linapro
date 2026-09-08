# 设计：直播公告与圣经阅读器

## 背景与数据来源

用户提供的 MySQL 导出（`grace` 库）：

- `bibleid`（66 行）：`sn` 主键、`kindsn` 分类、`chapternumber` 章节总数、`neworold`（数据实际取值 `0`=旧约、`1`=新约，与建表注释"1=旧约 2=新约"不符，以数据为准并在 PG 注释中修正语义）、`pinyin`、`shortname`、`fullname`。
- `bible`（31103 行）：`volumesn`/`chaptersn`/`versesn` 唯一键、`lection` 经文（≤500 字符）、`soundbegin`/`soundend` 音频毫秒位置（本期只入库备用，阅读器不消费）。

## 数据模型

### 公告表 `plugin_linapro_live_manage_announcement`

| 列 | 类型 | 说明 |
|----|------|------|
| id | BIGINT IDENTITY PK | |
| tenant_id | INT NOT NULL DEFAULT 0 | 归属租户，与直播间一致 |
| room_id | BIGINT NOT NULL | 关联直播间（软删直播间保留公告数据） |
| title | VARCHAR(256) NOT NULL DEFAULT '' | 公告标题 |
| content | VARCHAR(4000) NOT NULL DEFAULT '' | 公告内容（纯文本，H5 按换行渲染） |
| enabled | BOOLEAN NOT NULL DEFAULT TRUE | 启用开关，关闭后观众端不可见 |
| sort | INT NOT NULL DEFAULT 0 | 排序值，小者在前 |
| created_by/updated_by | BIGINT NOT NULL DEFAULT 0 | 审计 |
| created_at/updated_at | TIMESTAMPTZ | 宿主惯例 |
| deleted_at | TIMESTAMPTZ NULL | 软删（与 room/live 一致） |

索引：`(tenant_id, room_id)` 组合索引覆盖管理端按房间列表与观众端按房间查询；观众端查询附加 `enabled=true` 过滤。

### 圣经书卷表 `plugin_linapro_live_manage_bible_book`

对齐导出源结构：`sn INT PRIMARY KEY`（稳定业务键，非自增，Seed 显式写入合法）、`kind_sn INT`、`chapter_number INT`、`testament SMALLINT`（`1`=旧约、`2`=新约——转换时把源数据 `0/1` 归一为 `1/2`，消除源注释与数据的语义矛盾）、`pinyin`、`short_name`、`full_name`。

### 经文表 `plugin_linapro_live_manage_bible_verse`

`id BIGINT IDENTITY PK`（Seed 不写 id 列）、`volume_sn INT NOT NULL`、`chapter_sn INT NOT NULL`、`verse_sn INT NOT NULL`、`lection VARCHAR(500)`、`sound_begin INT NULL`、`sound_end INT NULL`；唯一键 `(volume_sn, chapter_sn, verse_sn)`；索引 `(volume_sn, chapter_sn)` 支撑单章查询。外键约束 `volume_sn → bible_book.sn`（与源库一致）。

### Seed DML 策略

- 31103 行经文按 500 行/条合并为多值 `INSERT ... SELECT ... WHERE NOT EXISTS` 等价幂等形式不可行（多值行无行级判重），采用 `INSERT ... ON CONFLICT (volume_sn, chapter_sn, verse_sn) DO NOTHING`（唯一键为真实业务键，符合 database.md 幂等要求）；书卷同理 `ON CONFLICT (sn) DO NOTHING`。
- 演示栈每次启动 `init --rebuild=true` 会重建 schema 后装载，多值批量 INSERT 保证装载耗时有界（约 63 条语句）。
- 转换以一次性脚本完成，脚本本身不入库（属生成交付物的过程动作）；转换中校验经文无反斜杠转义残留。

## 后端设计

### 接口（api/announcement/v1、api/bible/v1）

管理端（Auth+Tenancy+Permission 组）：

- `GET /announcement`（`permission:"live:room:announce"`）：`roomId` 必填，返回该直播间全部公告（含禁用），按 `sort asc, id asc` 排序，`pageSize` 上限 100；公告量天然有界，分页仅作防御边界。
- `POST /announcement`：创建（`roomId`、`title` 必填、`content`、`sort`、`enabled`），校验直播间归属当前租户。
- `PUT /announcement`：更新（`id` 必填，其余为可变字段）；软删记录视为不存在。
- `DELETE /announcement`：软删。

观众端（匿名 viewer 组）：

- `GET /announcements?roomCode=&tenantId=`：启用中公告列表（`sort asc, id asc`，上限 50 条），返回 `id/title/content/updatedAt(ms)` 最小投影；房间不存在或租户校验失败返回空列表不泄露存在性（与 `play` 契约一致：租户错误逐字透传 `play.CodePlayTenantRequired/Invalid`）。
- `GET /bible/books`：66 卷书卷（`sn/shortName/fullName/chapterCount/testament`），全局静态数据。
- `GET /bible/chapter?volumeSn=&chapter=`：单章经文（`verseSn/lection`），卷号越界或章号超出该书卷章节数返回空列表（不报错，防止探测）。

### 服务（internal/service/announcement、internal/service/bible）

- `announcement.Service`：依赖 `bizctx.Service`、`tenantcap.Service`、`usercap.Service`（管理端写路径走租户过滤与用户上下文）+ `play.Service`（观众端复用 `ResolveRoom` 房间解析，房间可见性与 `play` 单点对齐，不重复实现房间查询）。管理端 List/Create/Update/Delete + 观众端 ForRoom。
- `play.Service` 新增 `ResolveRoom(ctx, ResolveRoomInput) (*RoomRef, error)`：内部复用既有 `loadRoom`（租户 + roomCode，软删过滤），未命中返回 `nil, nil`；`RoomRef` 为契约自有投影（`Id/TenantId/RoomCode/RoomName`），不泄漏 entity。
- `bible.Service`：`Books(ctx)`、`Chapter(ctx, in)`。纯静态数据查询，无租户语义（经文内容无租户属性，属数据权限规则下的公开内容例外，测试覆盖）。
- 错误码（`announcement_code.go`）：`CodeAnnouncementNotFound`、`CodeAnnouncementRoomNotFound`（管理端：直播间不存在或不属于当前租户）。观众端不抛业务错误（空列表语义）。
- 分页/上限常量集中定义；DAO/DO 承载查询，禁止 `g.Map`。

### 路由绑定（plugin.go）

- viewer 组 JSON 子组追加绑定 announcement viewer 控制器与 bible 控制器（HandlerResponse 组内）。
- admin 组追加绑定 announcement 管理控制器。
- 服务实例均在 `registerRoutes` 启动装配一次，构造注入。

## 管理端前端

- `live-room-management.vue` 行操作新增"公告"（GhostButton，权限 `live:room:announce` 由宿主权限指令控制显示）。
- 新增 `live-room-announcement-modal.vue`（useVbenModal）：上部列表（标题、内容摘要、排序、启用开关、更新时间、编辑/删除），底部新增按钮打开内嵌编辑表单（title/content/sort/enabled）；删除走 Popconfirm；复用 qrcode/calendar modal 的既有样式与交互模式。
- `live-client.ts` 新增类型与五个方法（list/create/update/delete + roomId 查询）。

## H5 设计（纯静态 ES5，延续大厂暗色风格）

- 信息面板顶部新增 `tool-row`：两个等宽按钮"公告"（铃铛图标）"圣经"（书卷图标），点击各自打开底部滑出面板（overlay），再次点击或关闭按钮收回；两态（播放态/预告态）均展示，提示态（无直播/错误）隐藏。
- 公告面板：标题栏 + 关闭按钮 + 公告卡片列表（标题、内容按 `\n` 分行、更新时间）；空态"暂无公告"；加载失败静默空态。
- 圣经面板三视图：
  1. 书卷视图：旧约/新约分组标题 + 卷片名网格（4 列），点击进入章节视图；
  2. 章节视图：卷名标题 + 返回按钮 + 章节号网格，点击进入阅读视图；
  3. 阅读视图：`卷名 第N章` 标题 + 返回章节 + 逐节列表（节号高亮 + 经文）；底部"上一章/下一章"（跨卷连续导航，按 sn 顺序）。
- 数据流：书卷列表首次打开时请求一次并缓存于内存 state；章节切换按需请求；公告在页面每次进入可展示态时请求一次（从面板打开也强制刷新一次）。

## i18n

- apidoc：`zh-CN/apidoc/plugin-api-main.json` 补 announcement/bible 全部接口与字段翻译；`en-US` 保持空占位。
- 管理端：`manifest/i18n/{zh-CN,en-US}/plugin.json` 补公告弹窗、按钮、权限点名称等键；`menu.json` 补 `room:announce` 按钮键。
- 错误消息：`error.json` 双语补 `announcement.notFound`、`announcement.room.notFound`。
- H5 中文独立页（既有决策），无语言包依赖。

## 数据权限与例外说明

- 管理端公告 CRUD：走 Auth+Tenancy+Permission 全链路，`live:room:announce` 权限点，租户过滤由宿主 tenant filter 注入，直播间归属校验防跨租户写入。
- 观众端公告/圣经：匿名公开只读。公告以"直播间编码 + 显式租户参数"为权威边界且仅返回启用中数据，与 `play` 相同的公开例外；圣经为全球公共静态内容（无租户/组织属性），属合理公开例外；两者均不泄露房间/数据存在性（空列表语义）。

## 缓存一致性

无缓存、快照或派生状态引入；公告实时读库，圣经静态 Seed 数据随插件 schema 生命周期重建。无影响。

## 测试策略

- 单元测试：announcement 输入校验/分页规整纯函数；bible 章/卷号规整纯函数。
- E2E `TC008`（Docker 演示栈，stamp=450000000000）：管理端公告 CRUD + H5 公告面板联动（含禁用开关不可见）、圣经 books/chapter 接口契约、H5 阅读器三级导航与跨章翻页、房间不存在空列表语义、清理。
- 既有 TC001~TC007 全量回归。
