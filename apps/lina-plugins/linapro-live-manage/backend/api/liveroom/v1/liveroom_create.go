// This file declares the create-live-room request/response DTOs used by the
// linapro-live-manage source plugin.

package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// Live Room Create API

// CreateReq defines the request for creating a live room.
type CreateReq struct {
	g.Meta      `path:"/liveroom" method:"post" tags:"LiveRooms" summary:"Create live room" dc:"Create a live room owned by the current tenant. The room code must be unique within the tenant." permission:"live:room:add"`
	RoomCode    string `json:"roomCode" v:"required|length:1,64#gf.gvalid.rule.required|gf.gvalid.rule.length" dc:"Room code exposed to external systems, unique per tenant" eg:"ROOM-MAIN"`
	RoomName    string `json:"roomName" v:"required|length:1,128#gf.gvalid.rule.required|gf.gvalid.rule.length" dc:"Room name" eg:"Main hall room"`
	RoomType    int    `json:"roomType" v:"required|in:1,2,3#gf.gvalid.rule.required|gf.gvalid.rule.in" dc:"Room type: 1=gathering, 2=event, 3=other" eg:"1"`
	Status      *int   `json:"status" d:"0" dc:"Room status: 0=idle, 1=live, 2=disabled" eg:"0"`
	Description string `json:"description" dc:"Room description" eg:"Main hall gathering room"`
}

// CreateRes Live room create response
type CreateRes struct {
	Id int64 `json:"id" dc:"Live room ID" eg:"1"`
}
