// announcement_helpers.go carries the small pure helpers of the announcement
// service: error constructors and the shared play tenant-error passthrough
// predicate kept next to their single consumer.

package announcement

import (
	"lina-core/pkg/bizerr"
	"lina-plugin-linapro-live-manage/backend/internal/service/play"
)

// bizerrRoomNotFound builds the room-not-found business error.
func bizerrRoomNotFound() error {
	return bizerr.NewCode(CodeAnnouncementRoomNotFound)
}

// bizerrTitleRequired builds the missing-title business error.
func bizerrTitleRequired() error {
	return bizerr.NewCode(CodeAnnouncementTitleRequired)
}

// isTenantError reports whether the play service failure is a tenant
// contract error that must pass through to the anonymous caller verbatim,
// matching the viewer statistics service behavior.
func isTenantError(err error) bool {
	return bizerr.Is(err, play.CodePlayTenantRequired) || bizerr.Is(err, play.CodePlayTenantInvalid)
}
