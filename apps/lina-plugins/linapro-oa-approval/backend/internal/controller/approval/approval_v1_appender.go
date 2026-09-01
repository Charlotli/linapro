// approval_v1_appender.go implements the controller method that appends one
// new approver node to a pending approval request.

package approval

import (
	"context"

	v1 "lina-plugin-linapro-oa-approval/backend/api/approval/v1"
	approvalsvc "lina-plugin-linapro-oa-approval/backend/internal/service/approval"
)

// Appender appends one approver to an approval request
func (c *ControllerV1) Appender(ctx context.Context, req *v1.AppenderReq) (res *v1.AppenderRes, err error) {
	err = c.approvalSvc.Appender(ctx, approvalsvc.AppenderInput{
		Id:         req.Id,
		ApproverId: req.ApproverId,
	})
	if err != nil {
		return nil, err
	}
	return &v1.AppenderRes{}, nil
}
