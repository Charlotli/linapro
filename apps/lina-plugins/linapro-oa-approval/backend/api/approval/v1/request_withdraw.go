// This file declares the withdraw-approval-request request/response DTOs used
// by the linapro-oa-approval source plugin.

package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// Approval Request Withdraw API

// WithdrawReq defines the request for withdrawing one approval request.
type WithdrawReq struct {
	g.Meta `path:"/approval/request/{id}/withdraw" method:"put" tags:"ApprovalRequests" summary:"Withdraw request" dc:"Withdraw one pending request as the applicant. The request status becomes withdrawn and the current approver is notified." permission:"oa:request:submit"`
	Id     int64 `json:"id" v:"required" dc:"Approval request ID" eg:"1"`
}

// WithdrawRes Approval request withdraw response
type WithdrawRes struct{}
