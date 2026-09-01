# Equipment Management Plugin

`linapro-equipment-manage` is the official source plugin for equipment registry and maintenance record management.

## Capabilities

| Capability | Description |
| --- | --- |
| Equipment registry | Tenant-unique codes, dictionary-driven types/statuses, purchase info, location, and owner. |
| Maintenance records | Upkeep/repair/inspection history per equipment with cost, content, and result. |
| Reference protection | Equipment still referenced by maintenance records cannot be deleted. |

## Routes

| Menu | Path | Permission |
| --- | --- | --- |
| Equipment | `/equipment/register` | `equipment:list` |
| Maintenance records | `/equipment/maintenance` | `maintenance:list` |

## Data Tables

`plugin_linapro_equipment_manage_equipment`, `plugin_linapro_equipment_manage_maintenance`.

## Dictionary Types

`plugin_equipment_type`, `plugin_equipment_status`, `plugin_equipment_maint_type`.
