// This file wires the linapro-live-manage admin announcement controller. The
// controller converts the service projection into the API DTO and owns no
// business logic beyond request mapping.

package announcement

import (
	announcementapi "lina-plugin-linapro-live-manage/backend/api/announcement"
	announcementsvc "lina-plugin-linapro-live-manage/backend/internal/service/announcement"
)

// ControllerV1 is the tenant-scoped admin announcement controller.
type ControllerV1 struct {
	announcementSvc announcementsvc.Service // Tenant-scoped announcement service
}

// NewV1 creates and returns a new admin announcement controller instance.
func NewV1(announcementSvc announcementsvc.Service) announcementapi.IAnnouncementV1 {
	return &ControllerV1{announcementSvc: announcementSvc}
}

// ControllerViewerV1 is the anonymous viewer announcement controller.
type ControllerViewerV1 struct {
	announcementSvc announcementsvc.Service // Tenant-scoped announcement service
}

// NewViewerV1 creates and returns a new viewer announcement controller
// instance sharing the same service as the admin controller.
func NewViewerV1(announcementSvc announcementsvc.Service) announcementapi.IAnnouncementViewerV1 {
	return &ControllerViewerV1{announcementSvc: announcementSvc}
}
