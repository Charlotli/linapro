// flow_v1_delete.go implements the controller method that deletes approval
// flows by ID list.

package flow

import (
	"context"

	v1 "lina-plugin-linapro-oa-approval/backend/api/flow/v1"
)

// Delete deletes approval flows by ID list.
func (c *ControllerV1) Delete(ctx context.Context, req *v1.DeleteReq) (res *v1.DeleteRes, err error) {
	err = c.flowSvc.Delete(ctx, req.Ids)
	if err != nil {
		return nil, err
	}
	return &v1.DeleteRes{}, nil
}
