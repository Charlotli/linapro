// This file declares the create-approval-request request/response DTOs used
// by the linapro-oa-approval source plugin.

package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// Approval Request Create API

// CreateReq defines the request for submitting one approval request.
type CreateReq struct {
	g.Meta `path:"/approval/request" method:"post" tags:"ApprovalRequests" summary:"Submit approval request" dc:"Submit one approval request owned by the current user. The flow must exist, be enabled, and contain at least one node; the node list and form field definitions are frozen into the request at submit time. Form carries the dynamic form values keyed by the configured field keys and is validated against the frozen field snapshot." permission:"oa:request:submit"`
	FlowId int64          `json:"flowId" v:"required|min:1#gf.gvalid.rule.required|gf.gvalid.rule.min" dc:"Originating approval flow ID" eg:"1"`
	Title  string         `json:"title" v:"required|length:1,256#gf.gvalid.rule.required|gf.gvalid.rule.length" dc:"Request title" eg:"Youth camp fund request"`
	Form   map[string]any `json:"form" dc:"Dynamic form values keyed by the configured field keys; required fields must be present and values must match field types" eg:"{}"`
}

// CreateRes Approval request create response
type CreateRes struct {
	Id int64 `json:"id" dc:"Approval request ID" eg:"1"`
}
