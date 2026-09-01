// This file declares the approver user options request/response DTOs used by
// the linapro-oa-approval source plugin.

package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// Approver User Options API

// UserOptionsReq defines the request for listing bounded approver candidates.
type UserOptionsReq struct {
	g.Meta  `path:"/approval/user/options" method:"get" tags:"ApprovalFlows" summary:"Get approver options" dc:"Return the minimal user projection (id, name) for approver selection, filtered by keyword across username and nickname. At most 200 items are returned; excess results are not paged." permission:"oa:flow:query"`
	Keyword string `json:"keyword" dc:"Filter by username or nickname (fuzzy match); empty means all visible users" eg:"leader"`
}

// UserOptionsRes Approver user options response
type UserOptionsRes struct {
	List []*UserOptionItem `json:"list" dc:"Approver user option list" eg:"[]"`
}

// UserOptionItem Approver user minimal option projection
type UserOptionItem struct {
	Id   int64  `json:"id" dc:"User ID" eg:"2"`
	Name string `json:"name" dc:"User display name" eg:"leader"`
}
