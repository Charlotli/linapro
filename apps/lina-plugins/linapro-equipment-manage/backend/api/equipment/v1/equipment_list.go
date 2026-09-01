// This file declares the list-equipment request/response DTOs used by the
// linapro-equipment-manage source plugin.

package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// Equipment List API

// ListReq defines the request for listing equipment.
type ListReq struct {
	g.Meta   `path:"/equipment" method:"get" tags:"Equipment" summary:"Get equipment list" dc:"Query equipment by page, with filtering by name, code, type, and status." permission:"equipment:query"`
	PageNum  int    `json:"pageNum" d:"1" v:"min:1" dc:"Page number" eg:"1"`
	PageSize int    `json:"pageSize" d:"10" v:"min:1|max:100" dc:"Number of items per page" eg:"10"`
	Name     string `json:"name" dc:"Filter by equipment name or code (fuzzy match)" eg:"Projector"`
	Type     int    `json:"type" dc:"Filter by equipment type: 1=office, 2=it, 3=production, 4=other; 0 means no filter" eg:"2"`
	Status   *int   `json:"status" dc:"Filter by equipment status: 1=in use, 2=idle, 3=under repair, 4=scrapped; omitted means no filter" eg:"1"`
}

// ListRes Equipment list response
type ListRes struct {
	List  []*EquipmentItem `json:"list" dc:"Equipment list" eg:"[]"`
	Total int              `json:"total" dc:"Total number of items" eg:"12"`
}
