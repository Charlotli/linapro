// This file declares the maintenance mutation request/response DTOs used by
// the linapro-equipment-manage source plugin.

package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// Maintenance Create API

// CreateReq defines the request for creating one maintenance record.
type CreateReq struct {
	g.Meta      `path:"/maintenance" method:"post" tags:"MaintenanceRecords" summary:"Create maintenance record" dc:"Create one maintenance record owned by the current tenant. The equipment must exist within the same tenant." permission:"maintenance:add"`
	EquipmentId int64   `json:"equipmentId" v:"required|min:1#gf.gvalid.rule.required|gf.gvalid.rule.min" dc:"Maintained equipment ID within the same tenant" eg:"1"`
	MaintType   int     `json:"maintType" v:"required|in:1,2,3,4#gf.gvalid.rule.required|gf.gvalid.rule.in" dc:"Maintenance type: 1=maintenance, 2=repair, 3=inspection, 4=other" eg:"3"`
	MaintDate   string  `json:"maintDate" v:"date#gf.gvalid.rule.date" dc:"Maintenance date in YYYY-MM-DD format with date-only semantics; defaults to today when omitted" eg:"2026-04-18"`
	Maintainer  string  `json:"maintainer" dc:"Maintainer name" eg:"John"`
	Cost        float64 `json:"cost" v:"min:0#gf.gvalid.rule.min" dc:"Maintenance cost with two decimal places" eg:"0"`
	Content     string  `json:"content" dc:"Maintenance content description" eg:"Quarterly inspection"`
	Result      string  `json:"result" dc:"Maintenance result" eg:"Normal"`
	Remark      string  `json:"remark" dc:"Remark" eg:""`
}

// CreateRes Maintenance create response
type CreateRes struct {
	Id int64 `json:"id" dc:"Maintenance record ID" eg:"1"`
}

// Maintenance Update API

// UpdateReq defines the request for updating one maintenance record.
type UpdateReq struct {
	g.Meta      `path:"/maintenance/{id}" method:"put" tags:"MaintenanceRecords" summary:"Update maintenance record" dc:"Update the specified maintenance record. All fields are optional; omitted fields keep their current values. When EquipmentId is updated the equipment must exist within the same tenant." permission:"maintenance:edit"`
	Id          int64    `json:"id" v:"required" dc:"Maintenance record ID" eg:"1"`
	EquipmentId *int64   `json:"equipmentId" dc:"Maintained equipment ID within the same tenant" eg:"1"`
	MaintType   *int     `json:"maintType" dc:"Maintenance type: 1=maintenance, 2=repair, 3=inspection, 4=other" eg:"3"`
	MaintDate   *string  `json:"maintDate" v:"date#gf.gvalid.rule.date" dc:"Maintenance date in YYYY-MM-DD format with date-only semantics" eg:"2026-04-18"`
	Maintainer  *string  `json:"maintainer" dc:"Maintainer name" eg:"John"`
	Cost        *float64 `json:"cost" v:"min:0#gf.gvalid.rule.min" dc:"Maintenance cost with two decimal places" eg:"0"`
	Content     *string  `json:"content" dc:"Maintenance content description" eg:"Quarterly inspection"`
	Result      *string  `json:"result" dc:"Maintenance result" eg:"Normal"`
	Remark      *string  `json:"remark" dc:"Remark" eg:""`
}

// UpdateRes Maintenance update response
type UpdateRes struct{}

// Maintenance Delete API

// DeleteReq defines the request for deleting maintenance records.
type DeleteReq struct {
	g.Meta `path:"/maintenance" method:"delete" tags:"MaintenanceRecords" summary:"Delete maintenance records" dc:"Soft-delete one or more maintenance records by query array ids[]" permission:"maintenance:remove"`
	Ids    []int64 `json:"ids" v:"required|min-length:1" dc:"Maintenance record ID list as a query array" eg:"[1,2,3]"`
}

// DeleteRes Maintenance delete response
type DeleteRes struct{}

// Maintenance Get API

// GetReq defines the request for retrieving maintenance record details.
type GetReq struct {
	g.Meta `path:"/maintenance/{id}" method:"get" tags:"MaintenanceRecords" summary:"Get maintenance record details" dc:"Get maintenance record details by ID. Records outside the current tenant are reported as not found." permission:"maintenance:query"`
	Id     int64 `json:"id" v:"required" dc:"Maintenance record ID" eg:"1"`
}

// GetRes Maintenance record detail response
type GetRes struct {
	MaintenanceItem
}
