// This file declares the reject-approval-request request/response DTOs used
// by the linapro-oa-approval source plugin.

package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// Approval Request Reject API

// RejectReq defines the request for rejecting one approval request.
type RejectReq struct {
	g.Meta  `path:"/approval/request/{id}/reject" method:"put" tags:"ApprovalRequests" summary:"Reject request" dc:"Reject one pending request as the current node approver. The request status becomes rejected and the applicant is notified. A reply content explaining the rejection is required." permission:"oa:request:approve"`
	Id      int64  `json:"id" v:"required" dc:"Approval request ID" eg:"1"`
	Comment string `json:"comment" v:"required|max-length:2000#gf.gvalid.rule.required|gf.gvalid.rule.max-length" dc:"Required rejection reason stored on the timeline" eg:"Budget exceeds the quarterly quota"`
}

// RejectRes Approval request reject response
type RejectRes struct{}
