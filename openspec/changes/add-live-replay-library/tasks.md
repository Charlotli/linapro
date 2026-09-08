# Tasks: add-live-replay-library

## 1. 数据库

- [x] 1.1 插件迭代 SQL：`manifest/sql/{序号}-live-replay-library.sql`为直播表新增`replay_enabled BOOLEAN NOT NULL DEFAULT TRUE`（幂等`ADD COLUMN IF NOT EXISTS`），`manifest/sql/uninstall/`同步`DROP COLUMN IF EXISTS`；执行`make db.init`与插件`make dao`重生成工件
  - 说明：卸载 SQL 对`plugin_linapro_live_manage_live`整表执行`DROP TABLE IF EXISTS`，已覆盖`replay_enabled`列清理，无需单独`DROP COLUMN`

## 2. 后端

- [x] 2.1 `play`接口已结束回退分支增加`replay_enabled=true`过滤（进行中分支不受影响），与回放列表口径一致
- [x] 2.2 新增公开往期回放列表接口`GET /api/v1/replays`：DTO（`roomCode`/`tenantId`/`page`/`pageSize`，默认 10 上限 50，完整`dc`/`eg`标签）+ 分页最小投影查询（`state=2`+公开+开启回放，开播日期倒序，数据库侧分页），空列表不泄露存在性；路由挂既有公开分组
- [x] 2.3 管理端：直播内容列表响应新增`replayEnabled`；更新接口请求 DTO 新增可选`replayEnabled`（不传不改）
- [x] 2.4 插件`manifest/i18n/<locale>/apidoc/`同步新增`replays`接口`zh-CN`翻译（英文源文本）

## 3. 前端

- [x] 3.1 管理端直播内容列表新增"回放"列（行内开关，仅已结束行可操作）与新建/编辑弹窗"结束后提供公开回放"开关
- [x] 3.2 H5 回放态新增"往期回放"区块：`replays`首页数据渲染（封面缩略/标题/日期），点击切换播放流并更新标题徽标，当前场高亮；空列表或失败静默不渲染；进行中/预告不请求

## 4. 测试

- [x] 4.1 后端单元测试：`play`回退受开关过滤、开关默认值、`replays`分页与空列表语义；编译门禁与`make lint`
  - 实施记录：`play_replays_impl_test.go`覆盖分页规整（默认 10、上限 50 截断、越界/负值归一）；回退过滤与开关默认值依赖数据库行状态，由 E2E TC007a/c/d 以真实链路覆盖（关闭后 play 回退返回`LIVE_MANAGE_PLAY_NOT_FOUND`且列表为空、开启恢复）；`GOWORK=off go build/vet`通过，插件模块`golangci-lint`0 issues；宿主`make lint.go`因 linactl 宿主模块既有`errcheck`问题失败（非本变更文件，独立记录）
- [x] 4.2 契约/E2E：新增`TC007-live-replay-library.ts`覆盖开关生效（关闭后`play`回退与列表均不可见、开启恢复）、H5 往期回放跳转与高亮；执行前设置`E2E_PSQL_DOCKER`并在 Docker 演示栈全量通过
  - 实施记录：TC007 五用例（a 列表契约+播放回退、b H5 回放区块+卡片切换、c 关闭开关双向不可见、d 重新开启恢复、e 清理）在 Docker 演示栈 5/5 通过；live 模块全量回归 36/36 通过；执行代理模型不支持图片输入，按 testing.md 以文本断言 + 页面快照替代截图多模态审查

## 5. 验证与审查

- [x] 5.1 运行`openspec validate add-live-replay-library --strict`（工具不可用时记录阻断原因并以静态检索与结构自查替代）
- [x] 5.2 影响分析记录：`i18n`（接口翻译 + H5 中文独立页无影响判断）、缓存一致性（无缓存引入）、数据权限（匿名公开只读例外与`play`一致 + 管理端走既有更新接口权限）、开发工具跨平台（仅运行既有构建命令，记录无工具变更）、测试策略；调用`lina-review`完成审查
  - 实施记录：`i18n`——replays 接口英文源文本 + `zh-CN` apidoc 翻译齐备，`replayEnabled`请求/响应字段与开关文案双语入语言包（运行时 messages 验证可查），H5 往期回放区块为中文独立页文案；缓存一致性——回放列表实时读库、无缓存引入，无影响；数据权限——匿名公开只读例外与`play`同边界（显式租户 + `is_public=1` + `replay_enabled=true`），空列表不泄露房间存在性（TC007c 覆盖），管理端开关走既有`live`更新权限`live:live:edit`；开发工具跨平台——仅运行既有 build/image/test 命令，无工具变更；测试策略——分页纯函数单测 + E2E 全链路；`lina-review`见审查结论
