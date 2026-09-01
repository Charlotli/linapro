// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"time"
)

// Maintenance is the golang structure for table maintenance.
type Maintenance struct {
	Id          int64      `json:"id"          orm:"id"           description:"Maintenance record ID"`
	TenantId    int        `json:"tenantId"    orm:"tenant_id"    description:"Owning tenant ID, 0 means PLATFORM"`
	EquipmentId int64      `json:"equipmentId" orm:"equipment_id" description:"Maintained equipment ID within the same tenant"`
	MaintType   int        `json:"maintType"   orm:"maint_type"   description:"Maintenance type: 1=maintenance, 2=repair, 3=inspection, 4=other"`
	MaintDate   time.Time  `json:"maintDate"   orm:"maint_date"   description:"Maintenance date with date-only semantics"`
	Maintainer  string     `json:"maintainer"  orm:"maintainer"   description:"Maintainer name"`
	Cost        float64    `json:"cost"        orm:"cost"         description:"Maintenance cost with two decimal places"`
	Content     string     `json:"content"     orm:"content"      description:"Maintenance content description"`
	Result      string     `json:"result"      orm:"result"       description:"Maintenance result"`
	Remark      string     `json:"remark"      orm:"remark"       description:"Remark"`
	CreatedBy   int64      `json:"createdBy"   orm:"created_by"   description:"Creator"`
	UpdatedBy   int64      `json:"updatedBy"   orm:"updated_by"   description:"Updater"`
	CreatedAt   *time.Time `json:"createdAt"   orm:"created_at"   description:"Creation time"`
	UpdatedAt   *time.Time `json:"updatedAt"   orm:"updated_at"   description:"Update time"`
	DeletedAt   *time.Time `json:"deletedAt"   orm:"deleted_at"   description:"Deletion time"`
}
