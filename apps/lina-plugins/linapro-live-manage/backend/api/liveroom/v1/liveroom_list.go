// This file declares the list-live-room request/response DTOs used by the
// linapro-live-manage source plugin.

package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// Live Room List API

// ListReq defines the request for listing live rooms.
type ListReq struct {
	g.Meta   `path:"/liveroom" method:"get" tags:"LiveRooms" summary:"Get live room list" dc:"Query live rooms by page, with filtering by room name, room type, and status." permission:"live:room:query"`
	PageNum  int    `json:"pageNum" d:"1" v:"min:1" dc:"Page number" eg:"1"`
	PageSize int    `json:"pageSize" d:"10" v:"min:1|max:100" dc:"Number of items per page" eg:"10"`
	RoomName string `json:"roomName" dc:"Filter by room name (fuzzy match)" eg:"Main hall"`
	RoomType int    `json:"roomType" dc:"Filter by room type: 1=gathering, 2=event, 3=other; 0 means no filter" eg:"1"`
	Status   *int   `json:"status" dc:"Filter by room status: 0=idle, 1=live, 2=disabled; omitted means no filter" eg:"0"`
}

// ListRes Live room list response
type ListRes struct {
	List  []*ListItem `json:"list" dc:"Live room list" eg:"[]"`
	Total int         `json:"total" dc:"Total number of items" eg:"20"`
}

// ListItem Live room list item
type ListItem struct {
	RoomItem
	CreatedByName string `json:"createdByName" dc:"Creator username" eg:"admin"`
}
