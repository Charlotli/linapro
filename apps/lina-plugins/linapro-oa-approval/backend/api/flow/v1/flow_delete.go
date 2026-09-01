// This file declares the delete-approval-flow request/response DTOs used by
// the linapro-oa-approval source plugin.

package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// Approval Flow Delete API

// DeleteReq defines the request for deleting approval flows.
type DeleteReq struct {
	g.Meta `path:"/approval/flow" method:"delete" tags:"ApprovalFlows" summary:"Delete approval flows" dc:"Soft-delete one or more approval flows by query array ids[]. Submitted approval requests keep their frozen node snapshots and are unaffected." permission:"oa:flow:remove"`
	Ids    []int64 `json:"ids" v:"required|min-length:1" dc:"Approval flow ID list as a query array, e.g. ids[]=1&ids[]=2&ids[]=3" eg:"[1,2,3]"`
}

// DeleteRes Approval flow delete response
type DeleteRes struct{}
