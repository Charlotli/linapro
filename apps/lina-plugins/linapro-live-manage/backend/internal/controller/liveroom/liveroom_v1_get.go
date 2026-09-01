// liveroom_v1_get.go implements the controller method that serves the
// live-room detail endpoint.

package liveroom

import (
	"context"

	v1 "lina-plugin-linapro-live-manage/backend/api/liveroom/v1"
)

// Get returns live room details
func (c *ControllerV1) Get(ctx context.Context, req *v1.GetReq) (res *v1.GetRes, err error) {
	item, err := c.roomSvc.GetById(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &v1.GetRes{
		RoomItem:      toAPIRoomItem(item.RoomEntity),
		CreatedByName: item.CreatedByName,
	}, nil
}
