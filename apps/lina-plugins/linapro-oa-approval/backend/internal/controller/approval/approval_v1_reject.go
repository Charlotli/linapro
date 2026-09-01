// approval_v1_reject.go implements the controller method that rejects one
// pending approval request.

package approval

import (
	"context"

	v1 "lina-plugin-linapro-oa-approval/backend/api/approval/v1"
	approvalsvc "lina-plugin-linapro-oa-approval/backend/internal/service/approval"
)

// Reject rejects an approval request
func (c *ControllerV1) Reject(ctx context.Context, req *v1.RejectReq) (res *v1.RejectRes, err error) {
	err = c.approvalSvc.Reject(ctx, approvalsvc.ActionInput{
		Id:      req.Id,
		Comment: req.Comment,
	})
	if err != nil {
		return nil, err
	}
	return &v1.RejectRes{}, nil
}
