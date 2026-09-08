// This file declares the public replay-library request/response DTOs used by
// the linapro-live-manage viewer H5 page. The endpoint is anonymous: no JWT
// and no permission tag are required.

package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// Replays API

// ReplaysReq defines the anonymous paginated replay-library request.
type ReplaysReq struct {
	g.Meta   `path:"/replays" method:"get" tags:"LivePlay" summary:"List public replays of one live room" dc:"List the paginated public replay library of one live room for anonymous viewers: finished public lives with the replay switch enabled, newest first. Tenant scoping matches the public play query. Page defaults to 1 and page size defaults to 10 with a hard cap of 50; oversized values clamp to the cap. Unknown rooms and empty libraries return an empty list with a zero total so room existence is not leaked."`
	RoomCode string `json:"roomCode" v:"required" dc:"Live room code, unique within the tenant" eg:"MAIN-HALL"`
	TenantId *int   `json:"tenantId" dc:"Tenant ID used to scope the anonymous query; required when the tenant plugin is enabled, where 0 selects the platform tenant; ignored when the tenant plugin is disabled" eg:"1"`
	Page     int    `json:"page" dc:"Page number starting from 1; zero or negative defaults to 1" eg:"1"`
	PageSize int    `json:"pageSize" dc:"Page size; zero defaults to 10 and values above 50 clamp to 50" eg:"10"`
}

// ReplaysRes defines the paginated replay-library response.
type ReplaysRes struct {
	List  []ReplayItem `json:"list" dc:"Replay items of the requested page, newest live date first; empty when the room has no visible replays" eg:"[]"`
	Total int          `json:"total" dc:"Total number of visible replays for the room across all pages" eg:"12"`
}

// ReplayItem exposes one replay-library entry with the minimal viewer-facing
// projection. Administrative fields such as the push URL are intentionally
// excluded.
type ReplayItem struct {
	LiveId    int64  `json:"liveId" dc:"Live content ID of the replay entry" eg:"7"`
	Title     string `json:"title" dc:"Live title shown on the replay card" eg:"Sunday service live"`
	CoverUrl  string `json:"coverUrl" dc:"Cover image URL for the replay card; empty means the viewer renders a placeholder" eg:"https://example.com/assets/cover.jpg"`
	LiveUrl   string `json:"liveUrl" dc:"HLS replay URL for viewers; finished replays always expose it" eg:"https://example.com/hls/sunday.m3u8"`
	LiveDate  string `json:"liveDate" dc:"Calendar date of the live event in YYYY-MM-DD format with date-only semantics" eg:"2026-04-26"`
	StartTime *int64 `json:"startTime" dc:"Actual start time of the live as Unix timestamp in milliseconds; null means the live has no recorded start time" eg:"1776756000000"`
}
