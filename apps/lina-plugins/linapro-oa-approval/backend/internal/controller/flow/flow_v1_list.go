// flow_v1_list.go implements the controller method that serves the paged
// approval-flow list endpoint.

package flow

import (
	"context"

	v1 "lina-plugin-linapro-oa-approval/backend/api/flow/v1"
	flowsvc "lina-plugin-linapro-oa-approval/backend/internal/service/flow"
)

// List queries approval flows
func (c *ControllerV1) List(ctx context.Context, req *v1.ListReq) (res *v1.ListRes, err error) {
	out, err := c.flowSvc.List(ctx, flowsvc.ListInput{
		PageNum:  req.PageNum,
		PageSize: req.PageSize,
		FlowName: req.FlowName,
		FlowType: req.FlowType,
		Status:   req.Status,
	})
	if err != nil {
		return nil, err
	}
	items := make([]*v1.FlowItem, 0, len(out.List))
	for _, item := range out.List {
		projection := toAPIFlowItem(item)
		items = append(items, &projection)
	}
	return &v1.ListRes{List: items, Total: out.Total}, nil
}
