// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
)

// Equipment is the golang structure of table plugin_linapro_equipment_manage_equipment for DAO operations like Where/Data.
type Equipment struct {
	g.Meta        `orm:"table:plugin_linapro_equipment_manage_equipment, do:true"`
	Id            any        // Equipment ID
	TenantId      any        // Owning tenant ID, 0 means PLATFORM
	EquipmentCode any        // Equipment code, unique per tenant
	EquipmentName any        // Equipment name
	EquipmentType any        // Equipment type: 1=office, 2=it, 3=production, 4=other
	BrandModel    any        // Brand and model
	PurchaseDate  any        // Purchase date with date-only semantics
	PurchasePrice any        // Purchase price with two decimal places
	Location      any        // Storage location
	Owner         any        // Responsible person
	Status        any        // Equipment status: 1=in use, 2=idle, 3=under repair, 4=scrapped
	Remark        any        // Remark
	CreatedBy     any        // Creator
	UpdatedBy     any        // Updater
	CreatedAt     *time.Time // Creation time
	UpdatedAt     *time.Time // Update time
	DeletedAt     *time.Time // Deletion time
}
