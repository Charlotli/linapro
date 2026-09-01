// This file declares the pending-count request/response DTOs used by the
// linapro-oa-approval source plugin.

package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// Approval Pending Count API

// PendingCountReq defines the request for querying pending approval counts.
type PendingCountReq struct {
	g.Meta `path:"/approval/pending-count" method:"get" tags:"ApprovalRequests" summary:"Count pending approvals" dc:"Return the number of pending approval requests waiting for the current user within the tenant, intended for badge display." permission:"oa:request:query"`
}

// PendingCountRes Approval pending count response
type PendingCountRes struct {
	Count int64 `json:"count" dc:"Number of pending approval requests waiting for the current user" eg:"3"`
}
