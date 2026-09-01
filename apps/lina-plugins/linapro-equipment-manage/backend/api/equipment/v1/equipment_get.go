// This file declares the get-equipment request/response DTOs used by the
// linapro-equipment-manage source plugin.

package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// Equipment Get API

// GetReq defines the request for retrieving equipment details.
type GetReq struct {
	g.Meta `path:"/equipment/{id}" method:"get" tags:"Equipment" summary:"Get equipment details" dc:"Get equipment details by ID. Equipment outside the current tenant is reported as not found." permission:"equipment:query"`
	Id     int64 `json:"id" v:"required" dc:"Equipment ID" eg:"1"`
}

// GetRes Equipment detail response
type GetRes struct {
	EquipmentItem
}
