// flow_v1_update.go implements the controller method that updates one approval
// flow with optional full node replacement.

package flow

import (
	"context"

	v1 "lina-plugin-linapro-oa-approval/backend/api/flow/v1"
	"lina-plugin-linapro-oa-approval/backend/internal/service/approval"
	flowsvc "lina-plugin-linapro-oa-approval/backend/internal/service/flow"
)

// Update updates an approval flow
func (c *ControllerV1) Update(ctx context.Context, req *v1.UpdateReq) (res *v1.UpdateRes, err error) {
	var nodes []flowsvc.NodeInput
	if req.Nodes != nil {
		nodes = toServiceUpdateNodeInputs(req.Nodes)
	}
	var fields *[]approval.FieldConfig
	if req.Fields != nil {
		converted := toServiceUpdateFieldInputs(req.Fields)
		fields = &converted
	}
	err = c.flowSvc.Update(ctx, flowsvc.UpdateInput{
		Id:          req.Id,
		FlowType:    req.FlowType,
		FlowName:    req.FlowName,
		Description: req.Description,
		Status:      req.Status,
		Nodes:       nodes,
		Fields:      fields,
	})
	if err != nil {
		return nil, err
	}
	return &v1.UpdateRes{}, nil
}
