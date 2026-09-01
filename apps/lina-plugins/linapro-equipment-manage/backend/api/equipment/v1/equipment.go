// This file defines shared equipment response DTOs for the linapro-equipment-manage API.
package v1

// EquipmentItem exposes equipment fields visible to callers.
type EquipmentItem struct {
	Id            int64   `json:"id" dc:"Equipment ID" eg:"1"`
	EquipmentCode string  `json:"equipmentCode" dc:"Equipment code, unique per tenant" eg:"EQ-PROJECTOR-01"`
	EquipmentName string  `json:"equipmentName" dc:"Equipment name" eg:"Meeting room projector"`
	EquipmentType int     `json:"equipmentType" dc:"Equipment type: 1=office, 2=it, 3=production, 4=other" eg:"2"`
	BrandModel    string  `json:"brandModel" dc:"Brand and model" eg:"Epson CB-X06"`
	PurchaseDate  string  `json:"purchaseDate" dc:"Purchase date in YYYY-MM-DD format with date-only semantics; empty means unset" eg:"2025-06-15"`
	PurchasePrice float64 `json:"purchasePrice" dc:"Purchase price with two decimal places" eg:"3299.00"`
	Location      string  `json:"location" dc:"Storage location" eg:"Main hall meeting room"`
	Owner         string  `json:"owner" dc:"Responsible person" eg:"John"`
	Status        int     `json:"status" dc:"Equipment status: 1=in use, 2=idle, 3=under repair, 4=scrapped" eg:"1"`
	Remark        string  `json:"remark" dc:"Remark" eg:""`
	CreatedBy     int64   `json:"createdBy" dc:"Creator user ID" eg:"1"`
	CreatedByName string  `json:"createdByName" dc:"Creator username resolved in batch" eg:"admin"`
	CreatedAt     *int64  `json:"createdAt" dc:"Creation time as Unix timestamp in milliseconds" eg:"1776756000000"`
	UpdatedAt     *int64  `json:"updatedAt" dc:"Last updated time as Unix timestamp in milliseconds" eg:"1776757800000"`
}
