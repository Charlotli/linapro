// flow_v1_create.go implements the controller method that creates one approval
// flow with its ordered nodes.

package flow

import (
	"context"

	v1 "lina-plugin-linapro-oa-approval/backend/api/flow/v1"
	flowsvc "lina-plugin-linapro-oa-approval/backend/internal/service/flow"
)

// Create creates an approval flow
func (c *ControllerV1) Create(ctx context.Context, req *v1.CreateReq) (res *v1.CreateRes, err error) {
	status := 1
	if req.Status != nil {
		status = *req.Status
	}
	id, err := c.flowSvc.Create(ctx, flowsvc.CreateInput{
		FlowType:    req.FlowType,
		FlowName:    req.FlowName,
		Description: req.Description,
		Status:      status,
		Nodes:       toServiceNodeInputs(req.Nodes),
		Fields:      toServiceCreateFieldInputs(req.Fields),
	})
	if err != nil {
		return nil, err
	}
	return &v1.CreateRes{Id: id}, nil
}
