// This file declares the list-approval-request request/response DTOs used by
// the linapro-oa-approval source plugin.

package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// Approval Request List API

// ListReq defines the request for listing approval requests.
type ListReq struct {
	g.Meta   `path:"/approval/request" method:"get" tags:"ApprovalRequests" summary:"Get approval request list" dc:"Query approval requests by page within the participant visibility boundary. The scope switches between my submitted requests and requests waiting for my approval." permission:"oa:request:query"`
	PageNum  int    `json:"pageNum" d:"1" v:"min:1" dc:"Page number" eg:"1"`
	PageSize int    `json:"pageSize" d:"10" v:"min:1|max:100" dc:"Number of items per page" eg:"10"`
	Scope    string `json:"scope" d:"mine" v:"in:mine,pending#gf.gvalid.rule.in" dc:"List scope: mine=my submitted requests, pending=requests waiting for my approval" eg:"mine"`
	Title    string `json:"title" dc:"Filter by title (fuzzy match)" eg:"Youth camp"`
	FlowType int    `json:"flowType" dc:"Filter by flow type: 1=fund request, 2=write-off, 3=expense reimbursement; 0 means no filter" eg:"1"`
	Status   *int   `json:"status" dc:"Filter by request status: 1=pending, 2=approved, 3=rejected, 4=withdrawn; omitted means no filter" eg:"1"`
}

// ListRes Approval request list response
type ListRes struct {
	List  []*ListItem `json:"list" dc:"Approval request list" eg:"[]"`
	Total int         `json:"total" dc:"Total number of items" eg:"8"`
}

// ListItem Approval request list item
type ListItem struct {
	RequestItem
}
