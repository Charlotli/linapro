## Context

设备管理需求：设备登记（台账）与维护记录两个模块。台账记录设备基础信息与状态；维护记录跟踪每次保养/维修/巡检，与设备关联。

## Goal

1. `设备登记`页：维护设备台账（编号、名称、类型、品牌型号、购置信息、存放位置、负责人、状态）。
2. `维护记录`页：登记每台设备的维护历史（类型、日期、维护人、费用、内容与结果），支持按设备查看。

## 设计决策

- 表：`plugin_linapro_oa` 模式一致的前缀命名；设备表 `uk(tenant_id, equipment_code)`；维护表 `idx(tenant_id, equipment_id)`。
- 引用保护：删除有维护记录的设备返回业务错误；维护记录创建时校验设备存在（同租户）。
- 枚举走字典：`plugin_equipment_type`（办公设备/IT设备/生产设备/其他）、`plugin_equipment_status`（在用/闲置/维修中/报废）、`plugin_equipment_maint_type`（保养/维修/巡检/其他）。
- 价格字段 `NUMERIC(12,2)`；API 时间点字段毫秒时间戳；购置/维护日期为 `date-only`。
- 依赖注入：`bizctxcap/tenantcap/usercap` 来自宿主 `HostServices`，复用启动期共享实例。
- 组织维度例外：负责人为展示字段，不接入组织数据权限过滤，权威边界为租户隔离与权限点。

## Risks / Trade-offs

- 无维护计划/预警能力（后续迭代可加定时提醒）。
- 单存放位置文本字段，不做库房维度管理。

## Migration Plan

安装执行幂等建表+字典 seed；卸载清理；`mock-data` 提供演示数据。

## Open Questions

- 无。
