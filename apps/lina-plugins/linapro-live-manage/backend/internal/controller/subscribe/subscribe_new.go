// This file wires the linapro-live-manage public calendar subscription
// controller. The controller converts the service feed into the raw calendar
// HTTP response and owns no business logic beyond request mapping.

package subscribe

import (
	subscribeapi "lina-plugin-linapro-live-manage/backend/api/subscribe"
	calendarsvc "lina-plugin-linapro-live-manage/backend/internal/service/calendar"
)

// ControllerV1 is the public viewer calendar subscription controller.
type ControllerV1 struct {
	calendarSvc calendarsvc.Service // Anonymous calendar subscription service
}

// NewV1 creates and returns a new public calendar subscription controller
// instance.
func NewV1(calendarSvc calendarsvc.Service) subscribeapi.ISubscribeV1 {
	return &ControllerV1{calendarSvc: calendarSvc}
}
