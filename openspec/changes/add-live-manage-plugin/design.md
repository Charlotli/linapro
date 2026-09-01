## Context

用户参考既有系统的`liveroom`（直播间）与`live`（直播内容）两张`MySQL`表，要求实现`直播后台管理`功能。本设计将参考表结构映射为符合当前项目规范的`PostgreSQL`表，并落地为官方源码插件。

## Goal

直播后台管理插件提供两个管理页面：

1. `直播间管理`：维护直播间基础信息（编码、名称、类型、状态、描述）。
2. `直播内容管理`：维护每场直播内容（主题、正文、播放/推流/封面/页面地址、歌单、诗歌与讲道人员、经文、大纲、设备与接待信息、状态、公开性、直播日期与开播时间）。

## 非目标

- 不实现推拉流协议、转码、播放器或观众端页面；本插件只做后台信息维护。
- 不实现直播实时状态上报（`status`由人工维护）。
- 不在`v1`集成宿主文件中心的封面图上传（封面`URL`手工维护），后续迭代可接入`Files()`能力。

## 设计决策

### 插件形态与归属

- 归属：源码插件能力（`apps/lina-plugins/linapro-live-manage/`），业务前后端严格闭环在插件内部，不改`lina-core`。
- `type: source`、`distribution: managed`：通过插件管理显式安装启用。
- `scope_nature: tenant_aware`、`supports_multi_tenant: true`、`default_install_mode: tenant_scoped`：直播数据按租户隔离。
- 插件`ID`命名遵循`<author>-<domain>-<capability>`：`linapro-live-manage`。

### 数据表设计

参考表到插件表的映射（`PostgreSQL`方言，幂等`DDL`）：

| 参考表 | 插件表 | 说明 |
| --- | --- | --- |
| `liveroom.roomid`(varchar UUID) | `plugin_linapro_live_manage_room.id`(BIGINT identity) | 主键统一为项目惯例的`BIGINT`自增，另设`room_code`唯一业务编码承接对外编码语义 |
| `liveroom.roomcode/roomname/roomtype/statuscode/description` | `room_code/room_name/room_type/status/description` | `status`改用`SMALLINT`枚举 |
| `live.id`(UUID) | `plugin_linapro_live_manage_live.id`(BIGINT identity) | 同上 |
| `live.roomid`(varchar) | `live.room_id`(BIGINT) | 逻辑关联直播间主键，同租户内校验存在性 |
| `live.date` | `live.live_date`(DATE) | 日历日期，`date-only`语义，`API`返回`YYYY-MM-DD`字符串 |
| `live.url/pushUrl/liveUrl/coverUrl` | `page_url/push_url/live_url/cover_url` | 列名规范化为蛇形 |
| `live.songList`(json) | `song_list`(TEXT) | 存`JSON`数组文本`[{"name","singer","order"}]`，后端校验`JSON`合法性 |
| `live.deviceInfo`(json) | `device_info`(TEXT) | 简化为自由文本设备说明，避免`JSON`编辑器复杂度 |
| `live.type` | 移除 | 参考表注释已声明该字段作废并移入直播间`room_type` |
| `live.isDeleted/deletedTime` | `deleted_at` | 采用`GoFrame`自动软删除，替代业务层软删标记 |
| 新增 | `tenant_id/created_by/updated_by/created_at/updated_at` | 租户隔离与审计字段，`GoFrame`自动维护时间 |

索引：直播间`uk(tenant_id, room_code)`、`idx(tenant_id, status)`、`idx(tenant_id, room_type)`；直播内容`idx(tenant_id, room_id)`、`idx(tenant_id, state)`、`idx(tenant_id, live_date)`、`idx(tenant_id, start_time)`。列表默认按`id`倒序走主键，筛选走`tenant_id`前导组合索引，无全表扫描路径。

### 枚举治理

所有枚举通过字典模块维护标签，后端仅定义命名常量：

| 字典类型 | 枚举值 |
| --- | --- |
| `plugin_live_room_type` | `1`=正常聚会，`2`=活动，`3`=其他 |
| `plugin_live_room_status` | `0`=空闲，`1`=直播中，`2`=禁用 |
| `plugin_live_state` | `0`=未开始，`1`=进行中，`2`=已结束 |
| `plugin_live_public` | `1`=公开，`0`=私有 |

`seed SQL`幂等写入宿主`sys_dict_type/sys_dict_data`（`ON CONFLICT DO NOTHING`），前端通过字典`store`渲染标签。

### API 契约

`RESTful`，路径挂载在插件`API`前缀`/api/v1`下：

| 接口 | 方法与路径 | 权限 |
| --- | --- | --- |
| 直播间列表 | `GET /liveroom` | `live:room:query` |
| 直播间下拉候选 | `GET /liveroom/options` | `live:room:query` |
| 直播间详情 | `GET /liveroom/{id}` | `live:room:query` |
| 创建直播间 | `POST /liveroom` | `live:room:add` |
| 更新直播间 | `PUT /liveroom/{id}` | `live:room:edit` |
| 批量删除直播间 | `DELETE /liveroom` | `live:room:remove` |
| 直播内容列表 | `GET /live` | `live:live:query` |
| 直播内容详情 | `GET /live/{id}` | `live:live:query` |
| 创建直播内容 | `POST /live` | `live:live:add` |
| 更新直播内容 | `PUT /live/{id}` | `live:live:edit` |
| 批量删除直播内容 | `DELETE /live` | `live:live:remove` |

- 响应时间点字段（`createdAt/updatedAt/startTime`）为`Unix`毫秒时间戳（`int64`）；`liveDate`为`date-only`字符串并在`dc`说明。
- 列表接口支持分页（`pageSize`上限`100`）、直播间名称模糊、类型/状态筛选；下拉候选默认上限`100`条且有`limit`上限`200`，文档说明超限行为。
- 直播内容列表通过收集当前页`room_id`后一次批量查询直播间名称完成装配，`N+1`规避路径明确。

### 业务规则

- 删除直播间前校验同租户内是否存在直播内容引用（含软删除行则忽略，仅统计未删除行），存在时返回业务错误`LIVE_MANAGE_ROOM_REFERENCED`，提示先处理直播内容。
- 创建/更新直播内容时校验`room_id`存在且直播间状态不为`禁用`；直播间存在性校验在同租户范围内执行。
- 禁用直播间不级联修改已有直播内容，历史数据保留，重新启用后恢复可用（符合模块禁用不破坏数据的原则）。

### 数据权限

- 读取与写入均在数据库查询阶段通过`tenantspi.ApplyPluginTableFilter`注入租户过滤，禁止先取全量再内存过滤；写操作按`tenant_id`+`id`双条件更新，跨租户`ID`表现为不存在。
- 组织维度数据权限例外说明：直播直播间与直播内容无组织/部门归属语义，数据规模为教会场景的小规模管理数据，权威边界为租户隔离；本插件不接入组织数据权限过滤，拒绝策略为同租户内按权限点控制读写。该例外在测试中以租户隔离断言覆盖。
- 创建后记录可见性边界：新记录属于当前用户所在租户，同租户具备对应查询权限点的角色可见。

### 依赖注入

- `service`构造函数逐项显式注入`bizctxcap.Service`、`tenantcap.Service`、`usercap.Service`，来源为插件注册回调的宿主`HostServices`目录；不引入其他运行期依赖，无新增共享实例诉求。
- 直播内容服务校验直播间时直接复用插件内部`dao`，属于同一插件限界上下文，不构成跨模块契约。

### 前端

- 菜单挂载宿主`content`目录：直播间管理`/live/room`、直播内容`/live/content`；页面组件通过`pluginPageMeta.routePath`与菜单路径匹配，由宿主动态页壳加载，不改宿主前端代码。
- 页面遵循`useVbenVxeGrid`+`Page`+`useVbenModal`+`ghost-button`/`Popconfirm`交互；枚举列与查询下拉通过字典`store`渲染。
- 歌单采用动态行编辑（歌名+演唱者），提交时序列化为`JSON`数组文本；插件禁用后菜单与页面由宿主统一隐藏。

## 复杂度判断

不引入新抽象层：两个资源各自走`api/controller/service/dao`标准分层，服务层直接实现；批量名称装配用一次批量查询+内存映射完成。无需`manager/provider/adapter`等抽象。

## Risks / Trade-offs

- `song_list`以`TEXT`承载`JSON`：牺牲数据库侧结构化查询能力，换取`ORM`兼容与实现简单；当前无按歌曲维度查询的诉求。
- 直播间状态由人工维护，可能出现`进行中`状态与真实推流不一致；`v1`接受该简化。
- 封面`URL`手工维护：接入文件中心属后续迭代。

## Migration Plan

1. 插件安装时执行`manifest/sql/001-linapro-live-manage-schema.sql`（幂等建表+索引+字典 seed）。
2. 可选加载`manifest/sql/mock-data/001-linapro-live-manage-mock-data.sql`演示数据。
3. 卸载时执行`manifest/sql/uninstall/001-linapro-live-manage-schema.sql`删除插件表与字典 seed。

## Open Questions

- 无。
