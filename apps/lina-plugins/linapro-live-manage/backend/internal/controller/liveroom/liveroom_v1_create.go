// liveroom_v1_create.go implements the controller method that creates a new
// live room.

package liveroom

import (
	"context"

	v1 "lina-plugin-linapro-live-manage/backend/api/liveroom/v1"
	liveroomsvc "lina-plugin-linapro-live-manage/backend/internal/service/liveroom"
)

// Create creates a live room
func (c *ControllerV1) Create(ctx context.Context, req *v1.CreateReq) (res *v1.CreateRes, err error) {
	status := 0
	if req.Status != nil {
		status = *req.Status
	}
	id, err := c.roomSvc.Create(ctx, liveroomsvc.CreateInput{
		RoomCode:    req.RoomCode,
		RoomName:    req.RoomName,
		RoomType:    req.RoomType,
		Status:      status,
		Description: req.Description,
	})
	if err != nil {
		return nil, err
	}
	return &v1.CreateRes{Id: id}, nil
}
