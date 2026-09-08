// This file declares the anonymous viewer announcement list DTOs used by the
// linapro-live-manage source plugin. It intentionally carries no permission
// tag: the endpoint binds inside the anonymous viewer route group.

package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// Announcement Viewer List API

// ViewerListReq defines the anonymous request for the viewer announcement list.
type ViewerListReq struct {
	g.Meta   `path:"/announcements" method:"get" tags:"LiveViewer" summary:"Get viewer announcement list" dc:"Return the enabled announcements of one live room for the anonymous H5 viewer, ordered by sort then creation and capped at 50. Unknown rooms and empty rooms answer with an empty list so room existence is never leaked."`
	RoomCode string `json:"roomCode" v:"required|length:1,64#gf.gvalid.rule.required|gf.gvalid.rule.length" dc:"Live room code exposed to viewers" eg:"ROOM-MAIN"`
	TenantId *int   `json:"tenantId" dc:"Tenant ID scoping the query; required when the tenant capability is enabled, matching the play endpoint contract" eg:"0"`
}

// ViewerListRes Announcement viewer list response
type ViewerListRes struct {
	List []*ViewerAnnouncementItem `json:"list" dc:"Enabled announcements ordered for display" eg:"[]"`
}

// ViewerAnnouncementItem defines one viewer-side announcement projection.
type ViewerAnnouncementItem struct {
	Id        int64  `json:"id" dc:"Announcement ID" eg:"5"`
	Title     string `json:"title" dc:"Announcement title" eg:"Sunday service notice"`
	Content   string `json:"content" dc:"Announcement body, plain text" eg:"The service starts at 9:30."`
	UpdatedAt int64  `json:"updatedAt" dc:"Last update time, Unix timestamp in milliseconds" eg:"1776756000000"`
}
