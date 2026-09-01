// This file declares the stop-live request/response DTOs used by the
// linapro-live-manage source plugin.

package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// Live Content Stop API

// StopReq defines the request for stopping one live content record.
type StopReq struct {
	g.Meta `path:"/live/{id}/stop" method:"put" tags:"LiveContents" summary:"Stop live content" dc:"Transition one ongoing live content record to finished within the current tenant. The live room returns to idle in the same transaction while the recorded start time is preserved." permission:"live:live:stop"`
	Id     int64 `json:"id" v:"required" dc:"Live content ID" eg:"1"`
}

// StopRes Live content stop response
type StopRes struct{}
