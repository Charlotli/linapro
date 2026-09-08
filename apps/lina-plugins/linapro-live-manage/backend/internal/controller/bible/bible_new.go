// This file wires the linapro-live-manage bible controller. The controller
// converts the service projection into the API DTO and owns no business logic
// beyond response mapping.

package bible

import (
	bibleapi "lina-plugin-linapro-live-manage/backend/api/bible"
	biblesvc "lina-plugin-linapro-live-manage/backend/internal/service/bible"
)

// ControllerV1 is the anonymous viewer bible controller.
type ControllerV1 struct {
	bibleSvc biblesvc.Service // Static catalog reading service
}

// NewV1 creates and returns a new bible controller instance.
func NewV1(bibleSvc biblesvc.Service) bibleapi.IBibleV1 {
	return &ControllerV1{bibleSvc: bibleSvc}
}
