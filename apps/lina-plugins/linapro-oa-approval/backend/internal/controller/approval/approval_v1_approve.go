// approval_v1_approve.go implements the controller method that approves one
// pending approval request.

package approval

import (
	"context"

	v1 "lina-plugin-linapro-oa-approval/backend/api/approval/v1"
	approvalsvc "lina-plugin-linapro-oa-approval/backend/internal/service/approval"
)

// Approve approves an approval request
func (c *ControllerV1) Approve(ctx context.Context, req *v1.ApproveReq) (res *v1.ApproveRes, err error) {
	err = c.approvalSvc.Approve(ctx, approvalsvc.ActionInput{
		Id:      req.Id,
		Comment: req.Comment,
	})
	if err != nil {
		return nil, err
	}
	return &v1.ApproveRes{}, nil
}
