// This file declares the public calendar subscription request/response DTOs
// used by viewer calendar clients. The endpoint is anonymous: no JWT and no
// permission tag are required, and the route group intentionally omits the
// host auth, tenancy, and permission middlewares.

package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// SubscribeReq defines the anonymous request for the ICS calendar stream of
// one live room.
type SubscribeReq struct {
	g.Meta   `path:"/subscribe" method:"get" tags:"LiveSubscribe" summary:"Subscribe to the live calendar" dc:"Return the RFC 5545 ICS calendar stream of one live room for anonymous calendar clients. The stream covers public lives from the last 7 days through the next 90 days ordered by live date. Every entry carries a stable UID keyed by the live record, the live title, start and end times in UTC, and the viewer page link; lives without a recorded start time are emitted as all-day entries on the live date. The response body is a raw text/calendar stream instead of the JSON envelope, and it always answers HTTP 200: missing rooms, private schedules, and failed tenant validation all produce a valid empty calendar so the room existence stays undisclosed."`
	RoomCode string `json:"roomCode" v:"required" dc:"Live room code, unique within the tenant" eg:"MAIN-HALL"`
	TenantId *int   `json:"tenantId" dc:"Tenant ID used to scope the anonymous query; required when the tenant plugin is enabled, where 0 selects the platform tenant; ignored when the tenant plugin is disabled; invalid values yield an empty calendar instead of an error" eg:"1"`
}

// SubscribeRes is the raw calendar response contract. The handler writes the
// ICS document directly to the HTTP response with the text/calendar content
// type because calendar clients cannot parse the JSON envelope, so this
// struct carries no JSON body fields by design.
type SubscribeRes struct {
	g.Meta `mime:"text/calendar" dc:"RFC 5545 ICS calendar stream of the requested live room" eg:"BEGIN:VCALENDAR"`
}
