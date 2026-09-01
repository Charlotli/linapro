// This file exposes the plugin-local generated entities through the service
// packages so controller code can keep stable type names.

package equipment

import entitymodel "lina-plugin-linapro-equipment-manage/backend/internal/model/entity"

// EquipmentEntity mirrors the generated plugin_linapro_equipment_manage_equipment entity owned by this plugin.
type EquipmentEntity = entitymodel.Equipment

// MaintenanceEntity mirrors the generated plugin_linapro_equipment_manage_maintenance entity owned by this plugin.
type MaintenanceEntity = entitymodel.Maintenance
