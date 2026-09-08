// This file wires the linapro-live-manage public watch-statistics controller.
// The controller converts the service projection into the API DTO and owns no
// business logic beyond response mapping.

package view

import (
	viewapi "lina-plugin-linapro-live-manage/backend/api/view"
	viewsvc "lina-plugin-linapro-live-manage/backend/internal/service/view"
)

// ControllerV1 is the public viewer watch-statistics controller.
type ControllerV1 struct {
	viewSvc viewsvc.Service // Anonymous watch-statistics service
}

// NewV1 creates and returns a new public watch-statistics controller
// instance.
func NewV1(viewSvc viewsvc.Service) viewapi.IViewV1 {
	return &ControllerV1{viewSvc: viewSvc}
}
