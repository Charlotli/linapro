// This file declares the maintenance request/response DTOs used by the
// linapro-equipment-manage source plugin.

package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// Maintenance List API

// ListReq defines the request for listing maintenance records.
type ListReq struct {
	g.Meta      `path:"/maintenance" method:"get" tags:"MaintenanceRecords" summary:"Get maintenance record list" dc:"Query maintenance records by page, with filtering by equipment, maintenance type, and date range." permission:"maintenance:query"`
	PageNum     int    `json:"pageNum" d:"1" v:"min:1" dc:"Page number" eg:"1"`
	PageSize    int    `json:"pageSize" d:"10" v:"min:1|max:100" dc:"Number of items per page" eg:"10"`
	EquipmentId int64  `json:"equipmentId" dc:"Filter by equipment ID; 0 means no filter" eg:"1"`
	MaintType   int    `json:"maintType" dc:"Filter by maintenance type: 1=maintenance, 2=repair, 3=inspection, 4=other; 0 means no filter" eg:"3"`
	DateStart   string `json:"dateStart" dc:"Maintenance date range start in YYYY-MM-DD format, inclusive; empty means unbounded" eg:"2026-01-01"`
	DateEnd     string `json:"dateEnd" dc:"Maintenance date range end in YYYY-MM-DD format, inclusive; empty means unbounded" eg:"2026-12-31"`
}

// ListRes Maintenance record list response
type ListRes struct {
	List  []*MaintenanceItem `json:"list" dc:"Maintenance record list" eg:"[]"`
	Total int                `json:"total" dc:"Total number of items" eg:"5"`
}
