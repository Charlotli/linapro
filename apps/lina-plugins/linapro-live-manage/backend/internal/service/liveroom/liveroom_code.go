// This file defines linapro-live-manage live-room business error codes.

package liveroom

import (
	"github.com/gogf/gf/v2/errors/gcode"

	"lina-core/pkg/bizerr"
)

var (
	// CodeRoomNotFound reports that the requested live room does not exist in
	// the current tenant.
	CodeRoomNotFound = bizerr.MustDefine(
		"LIVE_MANAGE_ROOM_NOT_FOUND",
		"Live room does not exist",
		gcode.CodeNotFound,
	)
	// CodeRoomCodeExists reports that the room code is already used in the
	// current tenant.
	CodeRoomCodeExists = bizerr.MustDefine(
		"LIVE_MANAGE_ROOM_CODE_EXISTS",
		"Live room code already exists",
		gcode.CodeInvalidParameter,
	)
	// CodeRoomTypeInvalid reports an unknown live-room type value.
	CodeRoomTypeInvalid = bizerr.MustDefine(
		"LIVE_MANAGE_ROOM_TYPE_INVALID",
		"Live room type is invalid",
		gcode.CodeInvalidParameter,
	)
	// CodeRoomStatusInvalid reports an unknown live-room status value.
	CodeRoomStatusInvalid = bizerr.MustDefine(
		"LIVE_MANAGE_ROOM_STATUS_INVALID",
		"Live room status is invalid",
		gcode.CodeInvalidParameter,
	)
	// CodeRoomReferenced reports that a live room is still referenced by live
	// content and cannot be deleted.
	CodeRoomReferenced = bizerr.MustDefine(
		"LIVE_MANAGE_ROOM_REFERENCED",
		"Live room is still referenced by live content",
		gcode.CodeInvalidOperation,
	)
	// CodeRoomDeleteRequired reports that a delete operation received no room IDs.
	CodeRoomDeleteRequired = bizerr.MustDefine(
		"LIVE_MANAGE_ROOM_DELETE_REQUIRED",
		"Select at least one live room to delete",
		gcode.CodeInvalidParameter,
	)
)
