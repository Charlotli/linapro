// This file declares the list-announcement request/response DTOs used by the
// linapro-live-manage source plugin.

package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// Announcement List API

// ListReq defines the request for listing announcements of one live room.
type ListReq struct {
	g.Meta   `path:"/announcement" method:"get" tags:"LiveAnnouncements" summary:"Get announcement list of one live room" dc:"Query the announcements of one live room by page for the admin console, disabled rows included, ordered by sort then creation." permission:"live:room:announce"`
	RoomId   int64 `json:"roomId" v:"required|min:1#gf.gvalid.rule.required|gf.gvalid.rule.min" dc:"Owning live room ID within the current tenant" eg:"1"`
	PageNum  int   `json:"pageNum" d:"1" v:"min:1" dc:"Page number" eg:"1"`
	PageSize int   `json:"pageSize" d:"10" v:"min:1|max:100" dc:"Number of items per page, capped at 100" eg:"10"`
}

// ListRes Announcement list response
type ListRes struct {
	List  []*AnnouncementItem `json:"list" dc:"Announcement list of the requested page" eg:"[]"`
	Total int                 `json:"total" dc:"Total number of announcements of the room" eg:"3"`
}

// AnnouncementItem defines one admin-side announcement projection.
type AnnouncementItem struct {
	Id        int64  `json:"id" dc:"Announcement ID" eg:"5"`
	RoomId    int64  `json:"roomId" dc:"Owning live room ID" eg:"1"`
	Title     string `json:"title" dc:"Announcement title" eg:"Sunday service notice"`
	Content   string `json:"content" dc:"Announcement body, plain text" eg:"The service starts at 9:30."`
	Enabled   bool   `json:"enabled" dc:"Whether viewers can see the announcement" eg:"true"`
	Sort      int    `json:"sort" dc:"Display order, smaller values first" eg:"0"`
	CreatedAt int64  `json:"createdAt" dc:"Creation time, Unix timestamp in milliseconds" eg:"1776756000000"`
	UpdatedAt int64  `json:"updatedAt" dc:"Last update time, Unix timestamp in milliseconds" eg:"1776756000000"`
}
