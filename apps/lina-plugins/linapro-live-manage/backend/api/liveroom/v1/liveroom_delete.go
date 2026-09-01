// This file declares the delete-live-room request/response DTOs used by the
// linapro-live-manage source plugin.

package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// Live Room Delete API

// DeleteReq defines the request for deleting live rooms.
type DeleteReq struct {
	g.Meta `path:"/liveroom" method:"delete" tags:"LiveRooms" summary:"Delete live rooms" dc:"Soft-delete one or more live rooms by query array ids[]. Deleting a room referenced by live content is rejected." permission:"live:room:remove"`
	Ids    []int64 `json:"ids" v:"required|min-length:1" dc:"Live room ID list as a query array, e.g. ids[]=1&ids[]=2&ids[]=3" eg:"[1,2,3]"`
}

// DeleteRes Live room delete response
type DeleteRes struct{}
