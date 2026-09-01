// This file declares the resubmit-approval-request request/response DTOs used
// by the linapro-oa-approval source plugin.

package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// Approval Request Resubmit API

// ResubmitReq defines the request for re-submitting one rejected or withdrawn
// approval request.
type ResubmitReq struct {
	g.Meta `path:"/approval/request/{id}/resubmit" method:"put" tags:"ApprovalRequests" summary:"Resubmit request" dc:"Restart one rejected or withdrawn approval request as the applicant. The frozen node snapshot is preserved, the request restarts from the first node, and the first approver is notified." permission:"oa:request:submit"`
	Id     int64 `json:"id" v:"required" dc:"Approval request ID" eg:"1"`
}

// ResubmitRes Approval request resubmit response
type ResubmitRes struct{}
