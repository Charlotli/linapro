// This file declares the update-announcement request/response DTOs used by
// the linapro-live-manage source plugin.

package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// Announcement Update API

// UpdateReq defines the request for updating an announcement. Nil fields keep
// their stored values.
type UpdateReq struct {
	g.Meta  `path:"/announcement" method:"put" tags:"LiveAnnouncements" summary:"Update announcement" dc:"Update the mutable fields of one announcement owned by the current tenant; omitted fields keep their stored values." permission:"live:room:announce"`
	Id      int64   `json:"id" v:"required|min:1#gf.gvalid.rule.required|gf.gvalid.rule.min" dc:"Announcement ID" eg:"5"`
	Title   *string `json:"title" v:"length:1,256#gf.gvalid.rule.length" dc:"Announcement title; omitted keeps the stored title" eg:"Sunday service notice"`
	Content *string `json:"content" v:"length:0,4000#gf.gvalid.rule.length" dc:"Announcement body; omitted keeps the stored body" eg:"The service starts at 9:30."`
	Enabled *bool   `json:"enabled" dc:"Viewer visibility; omitted keeps the stored state" eg:"true"`
	Sort    *int    `json:"sort" dc:"Display order; omitted keeps the stored order" eg:"0"`
}

// UpdateRes Announcement update response
type UpdateRes struct{}
