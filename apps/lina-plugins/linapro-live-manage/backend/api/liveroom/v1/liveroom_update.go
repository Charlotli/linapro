// This file declares the update-live-room request/response DTOs used by the
// linapro-live-manage source plugin.

package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// Live Room Update API

// UpdateReq defines the request for updating a live room.
type UpdateReq struct {
	g.Meta      `path:"/liveroom/{id}" method:"put" tags:"LiveRooms" summary:"Update live room" dc:"Update the specified live room. All fields are optional; omitted fields keep their current values." permission:"live:room:edit"`
	Id          int64   `json:"id" v:"required" dc:"Live room ID" eg:"1"`
	RoomCode    *string `json:"roomCode" v:"length:1,64#gf.gvalid.rule.length" dc:"Room code exposed to external systems, unique per tenant" eg:"ROOM-MAIN"`
	RoomName    *string `json:"roomName" v:"length:1,128#gf.gvalid.rule.length" dc:"Room name" eg:"Main hall room"`
	RoomType    *int    `json:"roomType" dc:"Room type: 1=gathering, 2=event, 3=other" eg:"1"`
	Status      *int    `json:"status" dc:"Room status: 0=idle, 1=live, 2=disabled" eg:"0"`
	Description *string `json:"description" dc:"Room description" eg:"Updated description"`
}

// UpdateRes Live room update response
type UpdateRes struct{}
