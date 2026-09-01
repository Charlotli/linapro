// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
)

// Maintenance is the golang structure of table plugin_linapro_equipment_manage_maintenance for DAO operations like Where/Data.
type Maintenance struct {
	g.Meta      `orm:"table:plugin_linapro_equipment_manage_maintenance, do:true"`
	Id          any        // Maintenance record ID
	TenantId    any        // Owning tenant ID, 0 means PLATFORM
	EquipmentId any        // Maintained equipment ID within the same tenant
	MaintType   any        // Maintenance type: 1=maintenance, 2=repair, 3=inspection, 4=other
	MaintDate   any        // Maintenance date with date-only semantics
	Maintainer  any        // Maintainer name
	Cost        any        // Maintenance cost with two decimal places
	Content     any        // Maintenance content description
	Result      any        // Maintenance result
	Remark      any        // Remark
	CreatedBy   any        // Creator
	UpdatedBy   any        // Updater
	CreatedAt   *time.Time // Creation time
	UpdatedAt   *time.Time // Update time
	DeletedAt   *time.Time // Deletion time
}
