// This file defines shared maintenance response DTOs for the linapro-equipment-manage API.
package v1

// MaintenanceItem exposes maintenance record fields visible to callers.
type MaintenanceItem struct {
	Id            int64   `json:"id" dc:"Maintenance record ID" eg:"1"`
	EquipmentId   int64   `json:"equipmentId" dc:"Maintained equipment ID within the same tenant" eg:"1"`
	EquipmentName string  `json:"equipmentName" dc:"Equipment name resolved in batch" eg:"Meeting room projector"`
	MaintType     int     `json:"maintType" dc:"Maintenance type: 1=maintenance, 2=repair, 3=inspection, 4=other" eg:"3"`
	MaintDate     string  `json:"maintDate" dc:"Maintenance date in YYYY-MM-DD format with date-only semantics" eg:"2026-04-18"`
	Maintainer    string  `json:"maintainer" dc:"Maintainer name" eg:"John"`
	Cost          float64 `json:"cost" dc:"Maintenance cost with two decimal places" eg:"0"`
	Content       string  `json:"content" dc:"Maintenance content description" eg:"Quarterly inspection"`
	Result        string  `json:"result" dc:"Maintenance result" eg:"Normal"`
	Remark        string  `json:"remark" dc:"Remark" eg:""`
	CreatedAt     *int64  `json:"createdAt" dc:"Creation time as Unix timestamp in milliseconds" eg:"1776756000000"`
	UpdatedAt     *int64  `json:"updatedAt" dc:"Last updated time as Unix timestamp in milliseconds" eg:"1776757800000"`
}
