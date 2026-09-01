// This file defines linapro-live-manage live-content business error codes.

package live

import (
	"github.com/gogf/gf/v2/errors/gcode"

	"lina-core/pkg/bizerr"
)

var (
	// CodeLiveNotFound reports that the requested live content record does not
	// exist in the current tenant.
	CodeLiveNotFound = bizerr.MustDefine(
		"LIVE_MANAGE_LIVE_NOT_FOUND",
		"Live content does not exist",
		gcode.CodeNotFound,
	)
	// CodeRoomUnavailable reports that the referenced live room does not exist
	// in the current tenant or is disabled.
	CodeRoomUnavailable = bizerr.MustDefine(
		"LIVE_MANAGE_ROOM_UNAVAILABLE",
		"Live room does not exist or is disabled",
		gcode.CodeInvalidParameter,
	)
	// CodeSongListInvalid reports that the song list is not a valid JSON array.
	CodeSongListInvalid = bizerr.MustDefine(
		"LIVE_MANAGE_SONG_LIST_INVALID",
		"Song list must be a valid JSON array",
		gcode.CodeInvalidParameter,
	)
	// CodeLiveStateInvalid reports an unknown live state value.
	CodeLiveStateInvalid = bizerr.MustDefine(
		"LIVE_MANAGE_LIVE_STATE_INVALID",
		"Live state is invalid",
		gcode.CodeInvalidParameter,
	)
	// CodeLivePublicInvalid reports an unknown live visibility value.
	CodeLivePublicInvalid = bizerr.MustDefine(
		"LIVE_MANAGE_LIVE_PUBLIC_INVALID",
		"Live visibility is invalid",
		gcode.CodeInvalidParameter,
	)
	// CodeLiveDateInvalid reports a malformed live date value.
	CodeLiveDateInvalid = bizerr.MustDefine(
		"LIVE_MANAGE_LIVE_DATE_INVALID",
		"Live date must be in YYYY-MM-DD format",
		gcode.CodeInvalidParameter,
	)
	// CodeLiveDeleteRequired reports that a delete operation received no IDs.
	CodeLiveDeleteRequired = bizerr.MustDefine(
		"LIVE_MANAGE_LIVE_DELETE_REQUIRED",
		"Select at least one live content record to delete",
		gcode.CodeInvalidParameter,
	)
	// CodeLiveStateTransition reports that the current live state does not
	// allow the requested start or stop action.
	CodeLiveStateTransition = bizerr.MustDefine(
		"LIVE_MANAGE_LIVE_STATE_TRANSITION",
		"Live state does not allow this action",
		gcode.CodeInvalidOperation,
	)
	// CodeRoomBusy reports that the live room already hosts another ongoing
	// live content record.
	CodeRoomBusy = bizerr.MustDefine(
		"LIVE_MANAGE_ROOM_BUSY",
		"Live room already has an ongoing live",
		gcode.CodeInvalidOperation,
	)
	// CodeRoomImmutable reports that the live-room binding of a live content
	// record cannot be changed after creation.
	CodeRoomImmutable = bizerr.MustDefine(
		"LIVE_MANAGE_ROOM_IMMUTABLE",
		"Live room binding cannot be changed after creation",
		gcode.CodeInvalidParameter,
	)
)
