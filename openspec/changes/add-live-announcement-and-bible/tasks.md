# 任务：直播公告与圣经阅读器

## 1. 数据库与数据迁移

- [x] 1.1 新增 `manifest/sql/004-live-announcement-and-bible.sql`：公告/书卷/经文三表 DDL（幂等、含注释与索引）+ 66 卷书卷 Seed + 31103 节经文 Seed（多值 INSERT + `ON CONFLICT DO NOTHING`，不写自增 id）；由 `E:\my\mlztai\bibleid.sql`/`bible.sql` 转换；卸载 SQL 同步三表
- [x] 1.2 `hack/config.yaml` dao tables 增补三表；经 linactl 包装执行 DAO 生成（临时端口方案），生成物入 `internal/{dao,do,entity}`

## 2. 后端

- [x] 2.1 `play.Service` 增加 `ResolveRoom` 房间解析（复用 `loadRoom`，契约投影 `RoomRef`，未命中零值）；`GetInput` 租户语义不动
- [x] 2.2 `internal/service/announcement`：Service 接口与实现（管理端 CRUD 走租户过滤与直播间归属校验；观众端 ForRoom 走 ResolveRoom + enabled 过滤 + 上限 50）+ `announcement_code.go`（NotFound/RoomNotFound）
- [x] 2.3 `internal/service/bible`：Books/Chapter（卷号章号校验、越界空列表、单章集合查询）
- [x] 2.4 api DTO：`api/announcement/v1/`（list/create/update/delete/viewer）与 `api/bible/v1/`（books/chapter），`dc/eg` 标签齐备，时间字段 Unix 毫秒；控制器与 `plugin.go` 路由绑定（viewer 组 + admin 组）

## 3. 前端

- [x] 3.1 管理端：直播间行操作"公告"入口 + `live-room-announcement-modal.vue`（列表/新增/编辑/删除/启停）+ `live-client.ts` 方法 + i18n 键（zh-CN/en-US）+ `plugin.yaml` 公告按钮菜单声明（`live:room:announce`）
- [x] 3.2 H5：`tool-row` 公告/圣经入口 + 公告面板（列表/空态）+ 圣经阅读器（书卷→章节→经文三级视图、上一章/下一章跨卷导航）+ CSS；`hideAll` 纳入新元素

## 4. i18n

- [x] 4.1 `zh-CN/apidoc/plugin-api-main.json` 补 announcement/bible 接口翻译；`en-US` 保持空占位；`error.json` 双语错误键

## 5. 测试与验证

- [x] 5.1 单元测试（announcement 输入校验/分页规整、bible 卷章规整）；`go build/vet`、插件模块 `golangci-lint`、`make i18n.check` 通过
- [x] 5.2 E2E `TC008`（stamp=450000000000）：公告 CRUD→观众端可见→禁用隐藏、圣经接口契约（66 卷/单章/越界）、H5 公告面板与阅读器导航、清理；TC001~TC007 回归；Docker 演示栈全量通过（Seed 装载验证）

## 6. 验证与审查

- [x] 6.1 OpenSpec 严格校验（CLI 不可用时静态自查替代并记录）
- [x] 6.2 影响分析（i18n/缓存一致性/数据权限/开发工具跨平台/测试策略）+ `lina-review` 审查
