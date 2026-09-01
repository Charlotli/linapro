// This file declares the comment-approval-request request/response DTOs used
// by the linapro-oa-approval source plugin.

package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// Approval Request Comment API

// CommentReq defines the request for adding one reply to an approval request.
type CommentReq struct {
	g.Meta  `path:"/approval/request/{id}/comment" method:"put" tags:"ApprovalRequests" summary:"Reply to request" dc:"Append one reply to the approval timeline without changing the request status. Only the applicant and the snapshot approvers can reply while the request is pending." permission:"oa:request:comment"`
	Id      int64  `json:"id" v:"required" dc:"Approval request ID" eg:"1"`
	Comment string `json:"comment" v:"required|max-length:2000#gf.gvalid.rule.required|gf.gvalid.rule.max-length" dc:"Reply content" eg:"Invoice has been re-uploaded"`
}

// CommentRes Approval request comment response
type CommentRes struct{}
