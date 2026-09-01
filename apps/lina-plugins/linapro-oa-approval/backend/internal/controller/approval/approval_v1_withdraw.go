// approval_v1_withdraw.go implements the controller method that withdraws one
// pending approval request as the applicant.

package approval

import (
	"context"

	v1 "lina-plugin-linapro-oa-approval/backend/api/approval/v1"
)

// Withdraw withdraws an approval request
func (c *ControllerV1) Withdraw(ctx context.Context, req *v1.WithdrawReq) (res *v1.WithdrawRes, err error) {
	err = c.approvalSvc.Withdraw(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &v1.WithdrawRes{}, nil
}
