// approval_v1_pending_count.go implements the controller method that serves
// the pending approval count endpoint.

package approval

import (
	"context"

	v1 "lina-plugin-linapro-oa-approval/backend/api/approval/v1"
)

// PendingCount counts pending approvals for the current user
func (c *ControllerV1) PendingCount(ctx context.Context, req *v1.PendingCountReq) (res *v1.PendingCountRes, err error) {
	count, err := c.approvalSvc.PendingCount(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.PendingCountRes{Count: count}, nil
}
