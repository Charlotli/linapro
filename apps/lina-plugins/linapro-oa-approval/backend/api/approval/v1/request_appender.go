// This file declares the append-approver request/response DTOs used by the
// linapro-oa-approval source plugin.

package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// Approval Request Append Approver API

// AppenderReq defines the request for appending one approver to a pending
// approval request.
type AppenderReq struct {
	g.Meta     `path:"/approval/request/{id}/appender" method:"put" tags:"ApprovalRequests" summary:"Append approver" dc:"Insert one new approver node right after the current node of a pending request. Later node orders shift back and the request stays on its current node. Only the current node approver can append." permission:"oa:request:approve"`
	Id         int64 `json:"id" v:"required" dc:"Approval request ID" eg:"1"`
	ApproverId int64 `json:"approverId" v:"required|min:1#gf.gvalid.rule.required|gf.gvalid.rule.min" dc:"New approver user ID inserted after the current node" eg:"3"`
}

// AppenderRes Approval request append approver response
type AppenderRes struct{}
