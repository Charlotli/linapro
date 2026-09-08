// announcement_code.go centralizes the bizerr codes of the announcement
// service so every caller-facing failure reason stays in one auditable place.

package announcement

import (
	"github.com/gogf/gf/v2/errors/gcode"

	"lina-core/pkg/bizerr"
)

var (
	// CodeAnnouncementNotFound marks a missing or soft-deleted announcement.
	CodeAnnouncementNotFound = bizerr.MustDefine(
		"LIVE_MANAGE_ANNOUNCEMENT_NOT_FOUND",
		"Announcement not found",
		gcode.CodeNotFound,
	)
	// CodeAnnouncementRoomNotFound marks a room that does not exist in the
	// current tenant, guarding cross-tenant announcement writes.
	CodeAnnouncementRoomNotFound = bizerr.MustDefine(
		"LIVE_MANAGE_ANNOUNCEMENT_ROOM_NOT_FOUND",
		"Live room not found",
		gcode.CodeNotFound,
	)
	// CodeAnnouncementTitleRequired marks an announcement missing its title.
	CodeAnnouncementTitleRequired = bizerr.MustDefine(
		"LIVE_MANAGE_ANNOUNCEMENT_TITLE_REQUIRED",
		"Announcement title is required",
		gcode.CodeInvalidParameter,
	)
)
