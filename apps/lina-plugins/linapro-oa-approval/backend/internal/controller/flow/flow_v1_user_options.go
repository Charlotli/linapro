// flow_v1_user_options.go implements the controller method that serves the
// bounded approver candidate endpoint.

package flow

import (
	"context"

	v1 "lina-plugin-linapro-oa-approval/backend/api/flow/v1"
	flowsvc "lina-plugin-linapro-oa-approval/backend/internal/service/flow"
)

// UserOptions returns bounded approver candidates
func (c *ControllerV1) UserOptions(ctx context.Context, req *v1.UserOptionsReq) (res *v1.UserOptionsRes, err error) {
	items, err := c.flowSvc.UserOptions(ctx, flowsvc.UserOptionsInput{
		Keyword: req.Keyword,
	})
	if err != nil {
		return nil, err
	}
	list := make([]*v1.UserOptionItem, 0, len(items))
	for _, item := range items {
		list = append(list, &v1.UserOptionItem{
			Id:   item.Id,
			Name: item.Name,
		})
	}
	return &v1.UserOptionsRes{List: list}, nil
}
