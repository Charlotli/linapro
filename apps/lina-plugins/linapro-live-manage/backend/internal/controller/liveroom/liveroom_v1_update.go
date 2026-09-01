// liveroom_v1_update.go implements the controller method that updates an
// existing live room.

package liveroom

import (
	"context"

	v1 "lina-plugin-linapro-live-manage/backend/api/liveroom/v1"
	liveroomsvc "lina-plugin-linapro-live-manage/backend/internal/service/liveroom"
)

// Update updates a live room
func (c *ControllerV1) Update(ctx context.Context, req *v1.UpdateReq) (res *v1.UpdateRes, err error) {
	err = c.roomSvc.Update(ctx, liveroomsvc.UpdateInput{
		Id:          req.Id,
		RoomCode:    req.RoomCode,
		RoomName:    req.RoomName,
		RoomType:    req.RoomType,
		Status:      req.Status,
		Description: req.Description,
	})
	if err != nil {
		return nil, err
	}
	return &v1.UpdateRes{}, nil
}
