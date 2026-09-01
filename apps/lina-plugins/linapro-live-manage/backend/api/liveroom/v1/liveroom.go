// This file defines shared live-room response DTOs for the linapro-live-manage API.
package v1

// RoomItem exposes live-room fields visible to callers.
type RoomItem struct {
	Id          int64  `json:"id" dc:"Live room ID" eg:"1"`
	RoomCode    string `json:"roomCode" dc:"Room code exposed to external systems, unique per tenant" eg:"ROOM-MAIN"`
	RoomName    string `json:"roomName" dc:"Room name" eg:"Main hall room"`
	RoomType    int    `json:"roomType" dc:"Room type: 1=gathering, 2=event, 3=other" eg:"1"`
	Status      int    `json:"status" dc:"Room status: 0=idle, 1=live, 2=disabled" eg:"0"`
	Description string `json:"description" dc:"Room description" eg:"Main hall gathering room"`
	CreatedBy   int64  `json:"createdBy" dc:"Creator user ID" eg:"1"`
	UpdatedBy   int64  `json:"updatedBy" dc:"Last updated user ID" eg:"1"`
	CreatedAt   *int64 `json:"createdAt" dc:"Creation time as Unix timestamp in milliseconds" eg:"1776756000000"`
	UpdatedAt   *int64 `json:"updatedAt" dc:"Last updated time as Unix timestamp in milliseconds" eg:"1776757800000"`
}
