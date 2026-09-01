# OA审批插件

`linapro-oa-approval`是流程可配置的`OA`审批官方源码插件，覆盖请款、核销、报销三类审批场景。

## 能力清单

| 能力 | 说明 |
| --- | --- |
| 流程配置 | 按类型维护有序审批人节点，支持同租户唯一命名与启停治理。 |
| 审批单 | 提交请款/核销/报销申请，包含金额、申请说明与附件地址。 |
| 审批动作 | 通过、驳回、回复、加签、撤回，配合提交时冻结的节点快照执行。 |
| 站内通知 | 提交、流转、驳回、加签、撤回通过宿主通知能力触达相关用户。 |
| 附件地址 | 附件以`URL`地址存储（`JSON`数组，最多 10 条），不上传文件。 |

## 路由与权限

| 菜单 | 路径 | 权限 |
| --- | --- | --- |
| 流程配置 | `/oa/approval-flow` | `oa:flow:list` |
| 审批中心 | `/oa/approval` | `oa:request:list` |

按钮权限：`oa:flow:query`、`oa:flow:add`、`oa:flow:edit`、`oa:flow:remove`、`oa:request:query`、`oa:request:submit`、`oa:request:approve`、`oa:request:comment`。

## 数据表

| 表 | 用途 |
| --- | --- |
| `plugin_linapro_oa_approval_flow` | 租户级流程配置。 |
| `plugin_linapro_oa_approval_flow_node` | 有序审批人节点，同流程序号部分唯一索引兜底。 |
| `plugin_linapro_oa_approval_request` | 审批单，含冻结节点快照与当前审批人冗余列。 |
| `plugin_linapro_oa_approval_record` | 完整审批时间线与回复内容。 |

## 字典类型

| 类型 | 枚举值 |
| --- | --- |
| `plugin_oa_approval_flow_type` | `1`请款、`2`核销、`3`报销 |
| `plugin_oa_approval_status` | `1`审批中、`2`已通过、`3`已驳回、`4`已撤回 |
| `plugin_oa_approval_action` | `1`发起、`2`通过、`3`驳回、`4`回复、`5`加签、`6`撤回 |
