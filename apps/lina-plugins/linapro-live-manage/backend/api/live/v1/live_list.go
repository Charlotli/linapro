// This file declares the list-live request/response DTOs used by the
// linapro-live-manage source plugin.

package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// Live Content List API

// ListReq defines the request for listing live content.
type ListReq struct {
	g.Meta        `path:"/live" method:"get" tags:"LiveContents" summary:"Get live content list" dc:"Query live content by page, with filtering by title, live room, state, and visibility." permission:"live:live:query"`
	PageNum       int    `json:"pageNum" d:"1" v:"min:1" dc:"Page number" eg:"1"`
	PageSize      int    `json:"pageSize" d:"10" v:"min:1|max:100" dc:"Number of items per page" eg:"10"`
	Title         string `json:"title" dc:"Filter by title (fuzzy match)" eg:"Sunday"`
	RoomId        int64  `json:"roomId" dc:"Filter by live room ID; 0 means no filter" eg:"1"`
	State         *int   `json:"state" dc:"Filter by live state: 0=not started, 1=ongoing, 2=finished; omitted means no filter" eg:"1"`
	IsPublic      *int   `json:"isPublic" dc:"Filter by visibility: 1=public, 0=private; omitted means no filter" eg:"1"`
	LiveDateStart string `json:"liveDateStart" dc:"Filter by live date range start in YYYY-MM-DD format, inclusive; empty means unbounded" eg:"2026-04-01"`
	LiveDateEnd   string `json:"liveDateEnd" dc:"Filter by live date range end in YYYY-MM-DD format, inclusive; empty means unbounded" eg:"2026-04-30"`
}

// ListRes Live content list response
type ListRes struct {
	List  []*ListItem `json:"list" dc:"Live content list" eg:"[]"`
	Total int         `json:"total" dc:"Total number of items" eg:"20"`
}

// ListItem Live content list item
type ListItem struct {
	LiveItem
	CreatedByName string `json:"createdByName" dc:"Creator username" eg:"admin"`
	OnlineCount   int64  `json:"onlineCount" dc:"Watch sessions with a heartbeat inside the last 60 seconds; zero when nobody is watching" eg:"3"`
	TotalViews    int64  `json:"totalViews" dc:"Deduplicated watch sessions recorded for the live across page refreshes" eg:"42"`
}
