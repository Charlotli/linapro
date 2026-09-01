// liveroom_v1_delete.go implements the controller method that deletes live
// rooms by ID list.

package liveroom

import (
	"context"

	v1 "lina-plugin-linapro-live-manage/backend/api/liveroom/v1"
)

// Delete deletes live rooms by ID list.
func (c *ControllerV1) Delete(ctx context.Context, req *v1.DeleteReq) (res *v1.DeleteRes, err error) {
	err = c.roomSvc.Delete(ctx, req.Ids)
	if err != nil {
		return nil, err
	}
	return &v1.DeleteRes{}, nil
}
