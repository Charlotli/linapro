// approval_v1_comment.go implements the controller method that appends one
// reply to an approval request timeline.

package approval

import (
	"context"

	v1 "lina-plugin-linapro-oa-approval/backend/api/approval/v1"
	approvalsvc "lina-plugin-linapro-oa-approval/backend/internal/service/approval"
)

// Comment replies to an approval request
func (c *ControllerV1) Comment(ctx context.Context, req *v1.CommentReq) (res *v1.CommentRes, err error) {
	err = c.approvalSvc.Comment(ctx, approvalsvc.ActionInput{
		Id:      req.Id,
		Comment: req.Comment,
	})
	if err != nil {
		return nil, err
	}
	return &v1.CommentRes{}, nil
}
