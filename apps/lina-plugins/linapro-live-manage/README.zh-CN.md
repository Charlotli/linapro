# 直播后台管理插件

`linapro-live-manage`是直播后台管理官方源码插件，提供直播间维护与直播内容维护能力，具备租户隔离、字典化枚举和插件菜单治理。

## 能力清单

| 能力 | 说明 |
| --- | --- |
| 直播间管理 | 创建、查询、更新、删除直播间，支持同租户唯一编码、类型与状态维护。 |
| 直播内容管理 | 维护直播活动，包括主题、推流/播放/封面/页面地址、歌单、讲道信息、经文与状态。 |
| 直播间可用性治理 | 直播内容必须关联存在且未禁用的直播间；仍被引用的直播间禁止删除。 |
| 有界直播间候选 | 直播内容表单消费最小投影、数量有界的直播间候选接口。 |
| 观众端 H5 播放页 | 免登录的移动端页面，播放指定直播间的公开 HLS 流，支持回放与预告状态。 |

## 路由与权限

| 菜单 | 路径 | 权限 |
| --- | --- | --- |
| 直播间管理 | `/live/room` | `live:room:list` |
| 直播内容 | `/live/content` | `live:live:list` |

按钮权限：`live:room:query`、`live:room:add`、`live:room:edit`、`live:room:remove`、`live:live:query`、`live:live:add`、`live:live:edit`、`live:live:remove`。

## 观众端 H5 播放页

插件通过自身公开路由提供免登录、版本无关稳定 URL 的观众端页面：

| 地址 | 用途 |
| --- | --- |
| `GET /x/linapro-live-manage/h5?room={roomCode}&tenant={tenantId}` | 移动端 H5 播放页；仅租户插件启用时需要`tenant`参数。 |
| `GET /x/linapro-live-manage/api/v1/play?roomCode={roomCode}&tenantId={tenantId}` | 公开播放信息接口，遵循相同租户规则。 |

- 选择优先级：进行中公开直播优先，其次最近一场已结束公开直播作为回放，最后最近一场未开始公开直播作为预告（不返回播放地址）。
- 私密直播、其他租户与不存在的直播间统一返回同一种不存在错误；推流地址等管理端字段永不暴露。
- 播放使用内嵌`hls.js`（不依赖外网 CDN），iOS Safari 自动回退原生 HLS；未开始或无可看直播时页面每 30 秒自动刷新。

## 数据表

| 表 | 用途 |
| --- | --- |
| `plugin_linapro_live_manage_room` | 租户级直播间，含审计与软删除字段。 |
| `plugin_linapro_live_manage_live` | 租户级直播内容，含审计与软删除字段。 |

## 字典类型

| 类型 | 枚举值 |
| --- | --- |
| `plugin_live_room_type` | `1`正常聚会、`2`活动、`3`其他 |
| `plugin_live_room_status` | `0`空闲、`1`直播中、`2`禁用 |
| `plugin_live_state` | `0`未开始、`1`进行中、`2`已结束 |
| `plugin_live_public` | `1`公开、`0`私有 |
