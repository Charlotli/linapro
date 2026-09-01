# 设备管理插件

`linapro-equipment-manage`是设备登记与维护记录管理的官方源码插件。

## 能力清单

| 能力 | 说明 |
| --- | --- |
| 设备台账 | 编号同租户唯一，类型/状态字典治理，含购置信息、存放位置与负责人。 |
| 维护记录 | 按设备登记保养/维修/巡检历史，含费用、内容与结果。 |
| 引用保护 | 仍被维护记录引用的设备禁止删除。 |

## 路由与权限

| 菜单 | 路径 | 权限 |
| --- | --- | --- |
| 设备登记 | `/equipment/register` | `equipment:list` |
| 维护记录 | `/equipment/maintenance` | `maintenance:list` |

## 数据表

`plugin_linapro_equipment_manage_equipment`、`plugin_linapro_equipment_manage_maintenance`。

## 字典类型

`plugin_equipment_type`、`plugin_equipment_status`、`plugin_equipment_maint_type`。
