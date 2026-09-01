// This file declares the live-room options request/response DTOs used by the
// linapro-live-manage source plugin.

package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// Live Room Options API

// OptionsReq defines the request for listing bounded live-room candidates.
type OptionsReq struct {
	g.Meta          `path:"/liveroom/options" method:"get" tags:"LiveRooms" summary:"Get live room options" dc:"Return the minimal live-room projection (id, code, name, status) for select controls, filtered by keyword. At most 100 items are returned by default; pass limit to raise the cap up to 200 items, and results beyond the cap are not paged." permission:"live:room:query"`
	Keyword         string `json:"keyword" dc:"Filter by room name or room code (fuzzy match); empty means all rooms" eg:"Main"`
	IncludeDisabled *bool  `json:"includeDisabled" dc:"Whether to include disabled rooms; omitted or false means disabled rooms are excluded" eg:"false"`
	Limit           int    `json:"limit" dc:"Maximum number of options to return; omitted or non-positive means 100, and values above 200 are clamped to 200" eg:"100"`
}

// OptionsRes Live room options response
type OptionsRes struct {
	List []*OptionItem `json:"list" dc:"Live room option list" eg:"[]"`
}

// OptionItem Live room minimal option projection
type OptionItem struct {
	Id       int64  `json:"id" dc:"Live room ID" eg:"1"`
	RoomCode string `json:"roomCode" dc:"Room code" eg:"ROOM-MAIN"`
	RoomName string `json:"roomName" dc:"Room name" eg:"Main hall room"`
	Status   int    `json:"status" dc:"Room status: 0=idle, 1=live, 2=disabled" eg:"0"`
}
