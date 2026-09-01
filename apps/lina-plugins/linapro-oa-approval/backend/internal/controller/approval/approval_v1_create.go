// approval_v1_create.go implements the controller method that submits one
// approval request.

package approval

import (
	"context"

	v1 "lina-plugin-linapro-oa-approval/backend/api/approval/v1"
	approvalsvc "lina-plugin-linapro-oa-approval/backend/internal/service/approval"
)

// Create submits an approval request
func (c *ControllerV1) Create(ctx context.Context, req *v1.CreateReq) (res *v1.CreateRes, err error) {
	id, err := c.approvalSvc.Create(ctx, approvalsvc.CreateInput{
		FlowId: req.FlowId,
		Title:  req.Title,
		Form:   req.Form,
	})
	if err != nil {
		return nil, err
	}
	return &v1.CreateRes{Id: id}, nil
}
