# Tasks: add-live-replay-library

## 1. 数据库

- [ ] 1.1 插件迭代 SQL：`manifest/sql/{序号}-live-replay-library.sql`为直播表新增`replay_enabled BOOLEAN NOT NULL DEFAULT TRUE`（幂等`ADD COLUMN IF NOT EXISTS`），`manifest/sql/uninstall/`同步`DROP COLUMN IF EXISTS`；执行`make db.init`与插件`make dao`重生成工件

## 2. 后端

- [ ] 2.1 `play`接口已结束回退分支增加`replay_enabled=true`过滤（进行中分支不受影响），与回放列表口径一致
- [ ] 2.2 新增公开往期回放列表接口`GET /api/v1/replays`：DTO（`roomCode`/`tenantId`/`page`/`pageSize`，默认 10 上限 50，完整`dc`/`eg`标签）+ 分页最小投影查询（`state=2`+公开+开启回放，开播日期倒序，数据库侧分页），空列表不泄露存在性；路由挂既有公开分组
- [ ] 2.3 管理端：直播内容列表响应新增`replayEnabled`；更新接口请求 DTO 新增可选`replayEnabled`（不传不改）
- [ ] 2.4 插件`manifest/i18n/<locale>/apidoc/`同步新增`replays`接口`zh-CN`翻译（英文源文本）

## 3. 前端

- [ ] 3.1 管理端直播内容列表新增"回放"列（行内开关，仅已结束行可操作）与新建/编辑弹窗"结束后提供公开回放"开关
- [ ] 3.2 H5 回放态新增"往期回放"区块：`replays`首页数据渲染（封面缩略/标题/日期），点击切换播放流并更新标题徽标，当前场高亮；空列表或失败静默不渲染；进行中/预告不请求

## 4. 测试

- [ ] 4.1 后端单元测试：`play`回退受开关过滤、开关默认值、`replays`分页与空列表语义；编译门禁与`make lint`
- [ ] 4.2 契约/E2E：新增`TC007-live-replay-library.ts`覆盖开关生效（关闭后`play`回退与列表均不可见、开启恢复）、H5 往期回放跳转与高亮；执行前设置`E2E_PSQL_DOCKER`并在 Docker 演示栈全量通过

## 5. 验证与审查

- [ ] 5.1 运行`openspec validate add-live-replay-library --strict`（工具不可用时记录阻断原因并以静态检索与结构自查替代）
- [ ] 5.2 影响分析记录：`i18n`（接口翻译 + H5 中文独立页无影响判断）、缓存一致性（无缓存引入）、数据权限（匿名公开只读例外与`play`一致 + 管理端走既有更新接口权限）、开发工具跨平台（仅运行既有构建命令，记录无工具变更）、测试策略；调用`lina-review`完成审查
