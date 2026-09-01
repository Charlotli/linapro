// This file exposes the plugin-local generated maintenance entity through the
// service package so controller code can keep a stable type name.

package maintenance

import entitymodel "lina-plugin-linapro-equipment-manage/backend/internal/model/entity"

// MaintenanceEntity mirrors the generated plugin_linapro_equipment_manage_maintenance entity owned by this plugin.
type MaintenanceEntity = entitymodel.Maintenance
