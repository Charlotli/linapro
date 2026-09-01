// approval_v1_list.go implements the controller method that serves the paged
// approval-request list endpoint.

package approval

import (
	"context"

	v1 "lina-plugin-linapro-oa-approval/backend/api/approval/v1"
	approvalsvc "lina-plugin-linapro-oa-approval/backend/internal/service/approval"
)

// List queries approval requests
func (c *ControllerV1) List(ctx context.Context, req *v1.ListReq) (res *v1.ListRes, err error) {
	out, err := c.approvalSvc.List(ctx, approvalsvc.ListInput{
		PageNum:  req.PageNum,
		PageSize: req.PageSize,
		Scope:    req.Scope,
		Title:    req.Title,
		FlowType: req.FlowType,
		Status:   req.Status,
	})
	if err != nil {
		return nil, err
	}
	items := make([]*v1.ListItem, 0, len(out.List))
	for _, item := range out.List {
		projection := toAPIRequestItem(item)
		items = append(items, &v1.ListItem{RequestItem: projection})
	}
	return &v1.ListRes{List: items, Total: out.Total}, nil
}
