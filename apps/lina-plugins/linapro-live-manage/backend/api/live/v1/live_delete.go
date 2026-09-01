// This file declares the delete-live request/response DTOs used by the
// linapro-live-manage source plugin.

package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// Live Content Delete API

// DeleteReq defines the request for deleting live content.
type DeleteReq struct {
	g.Meta `path:"/live" method:"delete" tags:"LiveContents" summary:"Delete live content" dc:"Soft-delete one or more live content records by query array ids[]" permission:"live:live:remove"`
	Ids    []int64 `json:"ids" v:"required|min-length:1" dc:"Live content ID list as a query array, e.g. ids[]=1&ids[]=2&ids[]=3" eg:"[1,2,3]"`
}

// DeleteRes Live content delete response
type DeleteRes struct{}
