// This file declares the public play request/response DTOs used by the
// linapro-live-manage viewer H5 page. The endpoint is anonymous: no JWT and no
// permission tag are required, and the route group intentionally omits the
// host auth, tenancy, and permission middlewares.

package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// Play Get API

// GetReq defines the anonymous request for resolving one playable live by
// room code and explicit tenant ID.
type GetReq struct {
	g.Meta   `path:"/play" method:"get" tags:"LivePlay" summary:"Get public play info" dc:"Resolve one publicly visible live for anonymous viewers by live-room code and explicit tenant ID. The endpoint prefers an ongoing public live, falls back to the latest finished public live that has the replay switch enabled, and finally to the latest not-started public live as a preview. Private lives never resolve, and the play URL stays empty for not-started lives. The push URL and other administrative fields are never exposed."`
	RoomCode string `json:"roomCode" v:"required" dc:"Live room code, unique within the tenant" eg:"MAIN-HALL"`
	TenantId *int   `json:"tenantId" dc:"Tenant ID used to scope the anonymous query; required when the tenant plugin is enabled, where 0 selects the platform tenant; ignored when the tenant plugin is disabled" eg:"1"`
}

// GetRes public play info response
type GetRes struct {
	PlayItem
}
