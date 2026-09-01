# OA Approval Plugin

`linapro-oa-approval` is the official source plugin for configurable OA approvals covering fund requests, write-offs, and expense reimbursements.

## Capabilities

| Capability | Description |
| --- | --- |
| Approval flow configuration | Ordered approver nodes per flow type with tenant-unique names and enable/disable governance. |
| Approval requests | Submit fund/write-off/reimbursement requests with amount, statement, and attachment URL addresses. |
| Approval actions | Approve, reject, reply, append approver, and withdraw with a frozen per-request node snapshot. |
| Inbox notifications | Submissions, advances, rejections, appends, and withdrawals notify the related users through the host notify capability. |
| Attachment URLs | Attachments are stored as URL addresses (JSON array, max 10); no files are uploaded. |

## Routes

| Menu | Path | Permission |
| --- | --- | --- |
| Approval flows | `/oa/approval-flow` | `oa:flow:list` |
| Approval center | `/oa/approval` | `oa:request:list` |

Button permissions: `oa:flow:query`, `oa:flow:add`, `oa:flow:edit`, `oa:flow:remove`, `oa:request:query`, `oa:request:submit`, `oa:request:approve`, `oa:request:comment`.

## Data Tables

| Table | Purpose |
| --- | --- |
| `plugin_linapro_oa_approval_flow` | Tenant-scoped flow configuration. |
| `plugin_linapro_oa_approval_flow_node` | Ordered approver nodes with a partial unique index per flow order. |
| `plugin_linapro_oa_approval_request` | Requests with frozen node snapshots and current approver columns. |
| `plugin_linapro_oa_approval_record` | Full approval timeline with reply content. |

## Dictionary Types

| Type | Values |
| --- | --- |
| `plugin_oa_approval_flow_type` | `1` fund request, `2` write-off, `3` reimbursement |
| `plugin_oa_approval_status` | `1` pending, `2` approved, `3` rejected, `4` withdrawn |
| `plugin_oa_approval_action` | `1` submit, `2` approve, `3` reject, `4` reply, `5` append, `6` withdraw |
