// This file declares the create-announcement request/response DTOs used by
// the linapro-live-manage source plugin.

package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// Announcement Create API

// CreateReq defines the request for creating an announcement.
type CreateReq struct {
	g.Meta  `path:"/announcement" method:"post" tags:"LiveAnnouncements" summary:"Create announcement" dc:"Create one announcement for a live room owned by the current tenant. Enabled defaults to true when omitted." permission:"live:room:announce"`
	RoomId  int64  `json:"roomId" v:"required|min:1#gf.gvalid.rule.required|gf.gvalid.rule.min" dc:"Owning live room ID within the current tenant" eg:"1"`
	Title   string `json:"title" v:"required|length:1,256#gf.gvalid.rule.required|gf.gvalid.rule.length" dc:"Announcement title" eg:"Sunday service notice"`
	Content string `json:"content" v:"length:0,4000#gf.gvalid.rule.length" dc:"Announcement body, plain text rendered with line breaks" eg:"The service starts at 9:30."`
	Enabled *bool  `json:"enabled" dc:"Viewer visibility; omitted means enabled" eg:"true"`
	Sort    int    `json:"sort" dc:"Display order, smaller values first" eg:"0"`
}

// CreateRes Announcement create response
type CreateRes struct {
	Id int64 `json:"id" dc:"Created announcement ID" eg:"5"`
}
