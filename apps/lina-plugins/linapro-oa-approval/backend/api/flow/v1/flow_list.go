// This file declares the list-approval-flow request/response DTOs used by the
// linapro-oa-approval source plugin.

package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// Approval Flow List API

// ListReq defines the request for listing approval flows.
type ListReq struct {
	g.Meta   `path:"/approval/flow" method:"get" tags:"ApprovalFlows" summary:"Get approval flow list" dc:"Query approval flows by page, with filtering by flow name, flow type, and status." permission:"oa:flow:query"`
	PageNum  int    `json:"pageNum" d:"1" v:"min:1" dc:"Page number" eg:"1"`
	PageSize int    `json:"pageSize" d:"10" v:"min:1|max:100" dc:"Number of items per page" eg:"10"`
	FlowName string `json:"flowName" dc:"Filter by flow name (fuzzy match)" eg:"Standard"`
	FlowType int    `json:"flowType" dc:"Filter by flow type: 1=fund request, 2=write-off, 3=expense reimbursement; 0 means no filter" eg:"1"`
	Status   *int   `json:"status" dc:"Filter by flow status: 1=enabled, 0=disabled; omitted means no filter" eg:"1"`
}

// ListRes Approval flow list response
type ListRes struct {
	List  []*FlowItem `json:"list" dc:"Approval flow list" eg:"[]"`
	Total int         `json:"total" dc:"Total number of items" eg:"5"`
}
