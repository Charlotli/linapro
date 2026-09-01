// This file declares the create-equipment request/response DTOs used by the
// linapro-equipment-manage source plugin.

package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// Equipment Create API

// CreateReq defines the request for creating one equipment record.
type CreateReq struct {
	g.Meta        `path:"/equipment" method:"post" tags:"Equipment" summary:"Create equipment" dc:"Create one equipment record owned by the current tenant. The equipment code must be unique within the tenant." permission:"equipment:add"`
	EquipmentCode string  `json:"equipmentCode" v:"required|length:1,64#gf.gvalid.rule.required|gf.gvalid.rule.length" dc:"Equipment code, unique per tenant" eg:"EQ-PROJECTOR-01"`
	EquipmentName string  `json:"equipmentName" v:"required|length:1,128#gf.gvalid.rule.required|gf.gvalid.rule.length" dc:"Equipment name" eg:"Meeting room projector"`
	EquipmentType int     `json:"equipmentType" v:"required|in:1,2,3,4#gf.gvalid.rule.required|gf.gvalid.rule.in" dc:"Equipment type: 1=office, 2=it, 3=production, 4=other" eg:"2"`
	BrandModel    string  `json:"brandModel" dc:"Brand and model" eg:"Epson CB-X06"`
	PurchaseDate  string  `json:"purchaseDate" v:"date#gf.gvalid.rule.date" dc:"Purchase date in YYYY-MM-DD format with date-only semantics; empty means unset" eg:"2025-06-15"`
	PurchasePrice float64 `json:"purchasePrice" v:"min:0#gf.gvalid.rule.min" dc:"Purchase price with two decimal places" eg:"3299.00"`
	Location      string  `json:"location" dc:"Storage location" eg:"Main hall meeting room"`
	Owner         string  `json:"owner" dc:"Responsible person" eg:"John"`
	Status        *int    `json:"status" d:"1" dc:"Equipment status: 1=in use, 2=idle, 3=under repair, 4=scrapped" eg:"1"`
	Remark        string  `json:"remark" dc:"Remark" eg:""`
}

// CreateRes Equipment create response
type CreateRes struct {
	Id int64 `json:"id" dc:"Equipment ID" eg:"1"`
}
