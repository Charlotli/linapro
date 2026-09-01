// This file wires the linapro-live-manage live-room controller and shared
// response mappers.

package liveroom

import (
	"lina-core/pkg/apitime"
	liveroomapi "lina-plugin-linapro-live-manage/backend/api/liveroom"
	v1 "lina-plugin-linapro-live-manage/backend/api/liveroom/v1"
	liveroomsvc "lina-plugin-linapro-live-manage/backend/internal/service/liveroom"
)

// ControllerV1 is the live-room controller.
type ControllerV1 struct {
	roomSvc liveroomsvc.Service // live-room service
}

// NewV1 creates and returns a new live-room controller instance.
func NewV1(roomSvc liveroomsvc.Service) liveroomapi.ILiveroomV1 {
	return &ControllerV1{roomSvc: roomSvc}
}

// toAPIRoomItem converts the service-layer plugin-local room entity into the
// API DTO projection returned through HTTP responses.
func toAPIRoomItem(e *liveroomsvc.RoomEntity) v1.RoomItem {
	if e == nil {
		return v1.RoomItem{}
	}
	return v1.RoomItem{
		Id:          e.Id,
		RoomCode:    e.RoomCode,
		RoomName:    e.RoomName,
		RoomType:    e.RoomType,
		Status:      e.Status,
		Description: e.Description,
		CreatedBy:   e.CreatedBy,
		UpdatedBy:   e.UpdatedBy,
		CreatedAt:   apitime.Milli(e.CreatedAt),
		UpdatedAt:   apitime.Milli(e.UpdatedAt),
	}
}
