// This file declares the delete-equipment request/response DTOs used by the
// linapro-equipment-manage source plugin.

package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// Equipment Delete API

// DeleteReq defines the request for deleting equipment.
type DeleteReq struct {
	g.Meta `path:"/equipment" method:"delete" tags:"Equipment" summary:"Delete equipment" dc:"Soft-delete one or more equipment records by query array ids[]. Deleting equipment still referenced by maintenance records is rejected." permission:"equipment:remove"`
	Ids    []int64 `json:"ids" v:"required|min-length:1" dc:"Equipment ID list as a query array, e.g. ids[]=1&ids[]=2&ids[]=3" eg:"[1,2,3]"`
}

// DeleteRes Equipment delete response
type DeleteRes struct{}
