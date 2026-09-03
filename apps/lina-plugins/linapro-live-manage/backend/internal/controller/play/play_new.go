// This file wires the linapro-live-manage public play controller. The
// controller converts the service projection into the API DTO and owns no
// business logic beyond response mapping.

package play

import (
	playapi "lina-plugin-linapro-live-manage/backend/api/play"
	playsvc "lina-plugin-linapro-live-manage/backend/internal/service/play"
)

// ControllerV1 is the public viewer play controller.
type ControllerV1 struct {
	playSvc playsvc.Service // Anonymous play-info service
}

// NewV1 creates and returns a new public play controller instance.
func NewV1(playSvc playsvc.Service) playapi.IPlayV1 {
	return &ControllerV1{playSvc: playSvc}
}
