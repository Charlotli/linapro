// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"time"
)

// Equipment is the golang structure for table equipment.
type Equipment struct {
	Id            int64      `json:"id"            orm:"id"             description:"Equipment ID"`
	TenantId      int        `json:"tenantId"      orm:"tenant_id"      description:"Owning tenant ID, 0 means PLATFORM"`
	EquipmentCode string     `json:"equipmentCode" orm:"equipment_code" description:"Equipment code, unique per tenant"`
	EquipmentName string     `json:"equipmentName" orm:"equipment_name" description:"Equipment name"`
	EquipmentType int        `json:"equipmentType" orm:"equipment_type" description:"Equipment type: 1=office, 2=it, 3=production, 4=other"`
	BrandModel    string     `json:"brandModel"    orm:"brand_model"    description:"Brand and model"`
	PurchaseDate  time.Time  `json:"purchaseDate"  orm:"purchase_date"  description:"Purchase date with date-only semantics"`
	PurchasePrice float64    `json:"purchasePrice" orm:"purchase_price" description:"Purchase price with two decimal places"`
	Location      string     `json:"location"      orm:"location"       description:"Storage location"`
	Owner         string     `json:"owner"         orm:"owner"          description:"Responsible person"`
	Status        int        `json:"status"        orm:"status"         description:"Equipment status: 1=in use, 2=idle, 3=under repair, 4=scrapped"`
	Remark        string     `json:"remark"        orm:"remark"         description:"Remark"`
	CreatedBy     int64      `json:"createdBy"     orm:"created_by"     description:"Creator"`
	UpdatedBy     int64      `json:"updatedBy"     orm:"updated_by"     description:"Updater"`
	CreatedAt     *time.Time `json:"createdAt"     orm:"created_at"     description:"Creation time"`
	UpdatedAt     *time.Time `json:"updatedAt"     orm:"updated_at"     description:"Update time"`
	DeletedAt     *time.Time `json:"deletedAt"     orm:"deleted_at"     description:"Deletion time"`
}
