// liveroom_v1_list.go implements the controller method that serves the paged
// live-room list endpoint.

package liveroom

import (
	"context"

	v1 "lina-plugin-linapro-live-manage/backend/api/liveroom/v1"
	liveroomsvc "lina-plugin-linapro-live-manage/backend/internal/service/liveroom"
)

// List queries live rooms
func (c *ControllerV1) List(ctx context.Context, req *v1.ListReq) (res *v1.ListRes, err error) {
	out, err := c.roomSvc.List(ctx, liveroomsvc.ListInput{
		PageNum:  req.PageNum,
		PageSize: req.PageSize,
		RoomName: req.RoomName,
		RoomType: req.RoomType,
		Status:   req.Status,
	})
	if err != nil {
		return nil, err
	}
	items := make([]*v1.ListItem, 0, len(out.List))
	for _, item := range out.List {
		items = append(items, &v1.ListItem{
			RoomItem:      toAPIRoomItem(item.RoomEntity),
			CreatedByName: item.CreatedByName,
		})
	}
	return &v1.ListRes{List: items, Total: out.Total}, nil
}
