## Why

需要为教会直播场景提供`直播后台管理`能力：维护直播间（房间编码、名称、类型、状态）以及每场直播的内容（主题、推流/播放地址、歌单、讲道信息、经文、状态等）。业务参考既有系统的`liveroom`和`live`两张表结构，按当前框架规范落地为官方源码插件，闭环在`apps/lina-plugins/<plugin-id>/`内，不修改`lina-core`核心领域契约。

## What Changes

- 新增源码插件`linapro-live-manage`（直播后台管理），随宿主编译交付，`distribution: managed`，多租户`tenant_aware`。
- 新增两张插件业务表：
  - `plugin_linapro_live_manage_room`：直播间，含`room_code`、`room_name`、`room_type`、`status`、`description`与租户、审计、软删除字段。
  - `plugin_linapro_live_manage_live`：直播内容，含`room_id`、`title`、`content`、推流/播放/封面/页面地址、歌单、领唱、伴奏、主持人、讲道主题、讲员、身份、经文引用与内容、大纲、设备说明、接待安排、`state`、`is_public`、`live_date`、`start_time`与租户、审计、软删除字段。
- 新增直播间管理`RESTful`接口：分页列表、详情、创建、更新、批量删除、有界下拉候选；删除直播间时若仍存在直播内容引用则拒绝。
- 新增直播内容管理`RESTful`接口：分页列表、详情、创建、更新、批量删除；创建和更新时校验直播间存在且未禁用；列表按`room_id`批量装配直播间名称，避免`N+1`查询。
- 新增直播开启/关闭能力：`PUT /live/{id}/start`与`PUT /live/{id}/stop`状态动作接口（独立权限点`live:live:start`/`live:live:stop`）。开启记录开播时间、将直播间置为`直播中`，要求直播间可用且同房间无进行中的直播；关闭将直播置为`已结束`并将直播间恢复`空闲`；状态迁移与房间联动在同一事务内完成。
- 新增插件字典 seed：`plugin_live_room_type`、`plugin_live_room_status`、`plugin_live_state`、`plugin_live_public`，枚举标签由字典模块维护。
- 新增插件菜单（挂在宿主`content`内容管理目录下）：直播间管理、直播内容，以及查询、新增、修改、删除、开启直播、关闭直播按钮权限点。
- 新增插件前端页面：直播间管理页、直播内容管理页，遵循`useVbenVxeGrid`+`Page`+弹窗表单交互。
- 启用插件`i18n`（`zh-CN`/`en-US`），维护菜单、错误、插件文案与`apidoc`翻译资源。

## Capabilities

### New Capabilities

- `live-manage`：直播后台管理能力，覆盖直播间维护与直播内容维护的完整`CRUD`、下拉候选、字典枚举、租户隔离与权限点。

### Modified Capabilities

- 无。

## Impact

- 影响范围限定在`apps/lina-plugins/linapro-live-manage/`新增目录与`openspec/changes/add-live-manage-plugin/`变更文档；不修改`lina-core`与既有插件。
- `i18n`影响判断：插件启用`i18n`，新增`manifest/i18n`菜单、错误、插件文案资源与`apidoc`翻译资源；`API DTO`源文本使用英文。
- 缓存一致性影响判断：无缓存、快照或跨实例协调状态，全部数据实时读取数据库，无影响。
- 数据权限影响判断：列表、详情、候选与写操作均注入租户表过滤；组织维度数据权限例外说明见`design.md`。
- 开发工具跨平台影响判断：插件`Makefile`复用仓库共享`hack/makefiles/plugin.codegen.mk`，插件`hack/config.yaml`只含`dao`代码生成配置，无新增跨平台脚本。
- 测试策略影响判断：新增插件`API`契约测试、服务层单元测试与`hack/tests/e2e/live/`下的`E2E`用例；`openspec`工具未安装，`openspec validate`门禁记录为不可用。
