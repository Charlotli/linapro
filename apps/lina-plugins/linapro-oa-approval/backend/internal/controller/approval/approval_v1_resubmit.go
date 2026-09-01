// approval_v1_resubmit.go implements the controller method that re-submits
// one rejected or withdrawn approval request.

package approval

import (
	"context"

	v1 "lina-plugin-linapro-oa-approval/backend/api/approval/v1"
)

// Resubmit re-submits an approval request
func (c *ControllerV1) Resubmit(ctx context.Context, req *v1.ResubmitReq) (res *v1.ResubmitRes, err error) {
	err = c.approvalSvc.Resubmit(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &v1.ResubmitRes{}, nil
}
