## Context

用户需要`OA`审批功能：流程可配置（请款、核销、报销三类），支持配置审批人、回复内容、添加审批人、发送通知；附件以地址形式存储（服务器空间小，不做文件上传）。本设计将需求落地为官方源码插件`linapro-oa-approval`，复用宿主通知能力发送站内信。

## Goal

1. `流程配置`页：按类型维护审批流程与有序审批人节点（增删、排序、启停）。
2. `审批中心`页：提交请款/核销/报销申请（含金额、说明、附件地址），查看我的申请与待我审批，执行通过/驳回/回复/加签/撤回，查看完整审批时间线，关键进展收到站内通知。

## 非目标

- 不实现或签/会签（每节点单个审批人）、条件分支、转办与委托；作为后续迭代。
- 不实现附件文件上传与存储（附件仅存`URL`地址）。
- 不做审批模板打印与电子签章。

## 设计决策

### 插件形态与归属

- 源码插件`apps/lina-plugins/linapro-oa-approval/`，`type: source`、`distribution: managed`、`tenant_aware`、`default_install_mode: tenant_scoped`。
- 插件`ID`遵循`<author>-<domain>-<capability>`：`linapro-oa-approval`。

### 数据模型与快照语义

四张表：`plugin_oa_approval_flow`（流程）、`plugin_oa_approval_flow_node`（节点）、`plugin_oa_approval_request`（审批单）、`plugin_oa_approval_record`（时间线）。

关键决策：

- **节点快照**：审批单提交时将流程节点列表冻结为`node_snapshot`（`JSON`数组：`order/approverId/approverName`）。运行时动作全部基于快照，流程配置变更只影响未来提交，保证审批链完整性；加签在事务内修改快照并同步`current_approver_id`冗余列。
- **当前审批人冗余列**：`current_approver_id`随每次流转更新，使"待我审批"走数据库侧索引过滤（`idx(tenant_id, current_approver_id, status)`），避免快照`JSON`运行时过滤破坏分页。
- **节点唯一性**：节点表用部分唯一索引`uk(tenant_id, flow_id, node_order) WHERE deleted_at IS NULL`，软删除行不阻塞同序号新节点（配置频繁增删排序的场景）。
- **附件地址**：`attachments TEXT`存`JSON`字符串数组（如`["https://..."]`），后端校验`JSON`数组合法性且单条为`URL`格式，数量上限`10`；不上传文件。
- **金额**：`amount NUMERIC(14,2)`，三类审批共用。

### 状态机与动作

审批单状态：`1=审批中`、`2=已通过`、`3=已驳回`、`4=已撤回`（字典`plugin_oa_approval_status`）。时间线动作：`1=发起`、`2=通过`、`3=驳回`、`4=回复`、`5=加签`、`6=撤回`（字典`plugin_oa_approval_action`）。

| 动作 | 前置条件 | 副作用 |
| --- | --- | --- |
| 提交 | 流程存在且启用、节点非空 | 状态=审批中、冻结快照、当前节点=1、记录发起、通知首个审批人 |
| 通过 | 当前节点审批人、状态=审批中 | 写通过记录（可附回复）；末节点则状态=已通过并通知申请人，否则当前节点+1并通知下一审批人 |
| 驳回 | 当前节点审批人、状态=审批中 | 状态=已驳回、写驳回记录（回复内容建议填写）、通知申请人 |
| 回复 | 申请人或快照内审批人、状态=审批中 | 写回复记录，不做状态变更 |
| 加签 | 当前节点审批人、状态=审批中 | 快照当前节点后插入新审批人节点（后续序号+1）、写加签记录、更新当前审批人、通知新审批人 |
| 撤回 | 申请人、状态=审批中 | 状态=已撤回、写撤回记录、通知当前审批人 |

全部动作在同一事务内完成状态、记录与快照变更；通知在事务提交成功后发送，发送失败仅记录日志不回滚（与通知公告插件的降级语义一致）。

### 通知

复用宿主`notifycap.Service.Send`（`SourceType=plugin`、`SourceID=审批单ID`、`Recipients`为指定用户域`ID`、`Category=other`）。收件人使用用户域`ID`字符串（与`usercap`投影一致）。通知失败不阻断审批主流程。

### 可见性与数据权限

- 流程配置：租户隔离，持有流程查询权限的角色全量可见（配置面是管理数据）。
- 审批单：租户过滤 + 参与人可见性——申请人、快照内审批人可见详情与时间线，非参与人（含持有查询权限者）表现为不存在，不泄露单据存在性；列表接口以`scope`切换"我的申请"（`applicant_id`过滤）与"待我审批"（`current_approver_id+status`过滤）。
- 组织维度例外：审批人由流程显式指定，天然是人工授权边界，无组织归属推导诉求，不接入组织数据权限过滤；拒绝策略为参与人之外不可见。

### API 契约

`RESTful`，挂载插件`API`前缀`/api/v1`：

| 接口 | 方法与路径 | 权限 |
| --- | --- | --- |
| 流程列表 | `GET /approval/flow` | `oa:flow:query` |
| 流程详情 | `GET /approval/flow/{id}` | `oa:flow:query` |
| 创建流程 | `POST /approval/flow` | `oa:flow:add` |
| 更新流程 | `PUT /approval/flow/{id}` | `oa:flow:edit` |
| 删除流程 | `DELETE /approval/flow` | `oa:flow:remove` |
| 审批人候选 | `GET /approval/user/options` | `oa:flow:query` |
| 审批单列表 | `GET /approval/request` | `oa:request:query` |
| 审批单详情 | `GET /approval/request/{id}` | `oa:request:query` |
| 提交审批 | `POST /approval/request` | `oa:request:submit` |
| 通过 | `PUT /approval/request/{id}/approve` | `oa:request:approve` |
| 驳回 | `PUT /approval/request/{id}/reject` | `oa:request:approve` |
| 回复 | `PUT /approval/request/{id}/comment` | `oa:request:comment` |
| 加签 | `PUT /approval/request/{id}/appender` | `oa:request:approve` |
| 撤回 | `PUT /approval/request/{id}/withdraw` | `oa:request:submit` |

时间点响应字段为`Unix`毫秒时间戳；列表分页上限`100`；审批人候选上限`200`；流程节点数量上限`20`；回复内容长度上限`2000`。

### 复杂度判断

服务层按`flow`与`request`两个组件拆分：流程配置是普通`CRUD`+节点排序校验；审批单是状态机+快照+参与人可见性，均在单插件限界上下文内直接实现，不引入`workflow engine`抽象。快照解析/序列化收敛为`snapshot.go`职责文件。

## Risks / Trade-offs

- 单审批人节点：或签/会签后续迭代，表结构预留动作与时间线扩展。
- 快照使"改流程立即生效"只对新单据成立：这是有意的`OA`语义，避免在途单据被配置变更影响。
- 通知依赖宿主通知管道可用性：失败降级为日志，审批动作本身不失败。

## Migration Plan

1. 插件安装时执行`manifest/sql/001-linapro-oa-approval-schema.sql`（幂等建表+索引+字典 seed）。
2. 可选加载`manifest/sql/mock-data/001-linapro-oa-approval-mock-data.sql`演示流程与演示单据。
3. 卸载时执行`manifest/sql/uninstall/001-linapro-oa-approval-schema.sql`清理。

## Open Questions

- 无。
