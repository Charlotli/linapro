// This file declares the approve-approval-request request/response DTOs used
// by the linapro-oa-approval source plugin.

package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// Approval Request Approve API

// ApproveReq defines the request for approving one approval request.
type ApproveReq struct {
	g.Meta  `path:"/approval/request/{id}/approve" method:"put" tags:"ApprovalRequests" summary:"Approve request" dc:"Approve one pending request as the current node approver. The record advances to the next node or finishes the request on the last node. An optional reply content is stored on the approval timeline." permission:"oa:request:approve"`
	Id      int64  `json:"id" v:"required" dc:"Approval request ID" eg:"1"`
	Comment string `json:"comment" v:"max-length:2000#gf.gvalid.rule.max-length" dc:"Optional reply content stored on the timeline" eg:"Approved, please proceed"`
}

// ApproveRes Approval request approve response
type ApproveRes struct{}
