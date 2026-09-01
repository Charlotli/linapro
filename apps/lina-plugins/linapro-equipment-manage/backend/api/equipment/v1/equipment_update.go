// This file declares the update-equipment request/response DTOs used by the
// linapro-equipment-manage source plugin.

package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// Equipment Update API

// UpdateReq defines the request for updating one equipment record.
type UpdateReq struct {
	g.Meta        `path:"/equipment/{id}" method:"put" tags:"Equipment" summary:"Update equipment" dc:"Update the specified equipment. All fields are optional; omitted fields keep their current values." permission:"equipment:edit"`
	Id            int64    `json:"id" v:"required" dc:"Equipment ID" eg:"1"`
	EquipmentCode *string  `json:"equipmentCode" v:"length:1,64#gf.gvalid.rule.length" dc:"Equipment code, unique per tenant" eg:"EQ-PROJECTOR-01"`
	EquipmentName *string  `json:"equipmentName" v:"length:1,128#gf.gvalid.rule.length" dc:"Equipment name" eg:"Meeting room projector"`
	EquipmentType *int     `json:"equipmentType" dc:"Equipment type: 1=office, 2=it, 3=production, 4=other" eg:"2"`
	BrandModel    *string  `json:"brandModel" dc:"Brand and model" eg:"Epson CB-X06"`
	PurchaseDate  *string  `json:"purchaseDate" v:"date#gf.gvalid.rule.date" dc:"Purchase date in YYYY-MM-DD format with date-only semantics" eg:"2025-06-15"`
	PurchasePrice *float64 `json:"purchasePrice" v:"min:0#gf.gvalid.rule.min" dc:"Purchase price with two decimal places" eg:"3299.00"`
	Location      *string  `json:"location" dc:"Storage location" eg:"Main hall meeting room"`
	Owner         *string  `json:"owner" dc:"Responsible person" eg:"John"`
	Status        *int     `json:"status" dc:"Equipment status: 1=in use, 2=idle, 3=under repair, 4=scrapped" eg:"1"`
	Remark        *string  `json:"remark" dc:"Remark" eg:""`
}

// UpdateRes Equipment update response
type UpdateRes struct{}
