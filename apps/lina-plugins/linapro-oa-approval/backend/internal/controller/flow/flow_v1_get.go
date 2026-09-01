// flow_v1_get.go implements the controller method that serves the approval
// flow detail endpoint.

package flow

import (
	"context"

	v1 "lina-plugin-linapro-oa-approval/backend/api/flow/v1"
	flowsvc "lina-plugin-linapro-oa-approval/backend/internal/service/flow"
)

// Get returns approval flow details
func (c *ControllerV1) Get(ctx context.Context, req *v1.GetReq) (res *v1.GetRes, err error) {
	detail, err := c.flowSvc.GetById(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &v1.GetRes{
		FlowItem: toAPIFlowItem(&flowsvc.FlowItem{FlowEntity: detail.FlowEntity}),
		Nodes:    toAPINodeItems(detail.Nodes),
		Fields:   toAPIFieldItems(detail.Fields),
	}, nil
}
