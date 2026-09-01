// This file declares the get-live-room request/response DTOs used by the
// linapro-live-manage source plugin.

package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// Live Room Get API

// GetReq defines the request for retrieving live-room details.
type GetReq struct {
	g.Meta `path:"/liveroom/{id}" method:"get" tags:"LiveRooms" summary:"Get live room details" dc:"Get live room details by ID, including creator information. Rooms outside the current tenant are reported as not found." permission:"live:room:query"`
	Id     int64 `json:"id" v:"required" dc:"Live room ID" eg:"1"`
}

// GetRes Live room detail response
type GetRes struct {
	RoomItem
	CreatedByName string `json:"createdByName" dc:"Creator username" eg:"admin"`
}
