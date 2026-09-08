// This file declares the delete-announcement request/response DTOs used by
// the linapro-live-manage source plugin.

package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// Announcement Delete API

// DeleteReq defines the request for deleting an announcement.
type DeleteReq struct {
	g.Meta `path:"/announcement" method:"delete" tags:"LiveAnnouncements" summary:"Delete announcement" dc:"Soft-delete one announcement owned by the current tenant; repeated deletes succeed silently." permission:"live:room:announce"`
	Id     int64 `json:"id" v:"required|min:1#gf.gvalid.rule.required|gf.gvalid.rule.min" dc:"Announcement ID" eg:"5"`
}

// DeleteRes Announcement delete response
type DeleteRes struct{}
