// This file declares the start-live request/response DTOs used by the
// linapro-live-manage source plugin.

package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// Live Content Start API

// StartReq defines the request for starting one live content record.
type StartReq struct {
	g.Meta `path:"/live/{id}/start" method:"put" tags:"LiveContents" summary:"Start live content" dc:"Transition one not-started live content record to ongoing within the current tenant. The live room must exist, be enabled, and have no other ongoing live. The start time is recorded and the room becomes live in the same transaction." permission:"live:live:start"`
	Id     int64 `json:"id" v:"required" dc:"Live content ID" eg:"1"`
}

// StartRes Live content start response
type StartRes struct{}
