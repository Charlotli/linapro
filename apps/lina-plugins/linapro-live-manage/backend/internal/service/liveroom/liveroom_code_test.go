// This file verifies linapro-live-manage structured business error metadata.

package liveroom

import (
	"testing"

	"lina-core/pkg/bizerr"
)

// TestRoomBusinessErrorMetadata verifies live-room errors expose stable runtime
// codes and i18n keys instead of fixed-language text.
func TestRoomBusinessErrorMetadata(t *testing.T) {
	testCases := []struct {
		name        string
		code        *bizerr.Code
		runtimeCode string
		messageKey  string
	}{
		{
			name:        "room not found",
			code:        CodeRoomNotFound,
			runtimeCode: "LIVE_MANAGE_ROOM_NOT_FOUND",
			messageKey:  "error.live.manage.room.not.found",
		},
		{
			name:        "room code exists",
			code:        CodeRoomCodeExists,
			runtimeCode: "LIVE_MANAGE_ROOM_CODE_EXISTS",
			messageKey:  "error.live.manage.room.code.exists",
		},
		{
			name:        "room referenced",
			code:        CodeRoomReferenced,
			runtimeCode: "LIVE_MANAGE_ROOM_REFERENCED",
			messageKey:  "error.live.manage.room.referenced",
		},
		{
			name:        "delete required",
			code:        CodeRoomDeleteRequired,
			runtimeCode: "LIVE_MANAGE_ROOM_DELETE_REQUIRED",
			messageKey:  "error.live.manage.room.delete.required",
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			err := bizerr.NewCode(testCase.code)
			messageErr, ok := bizerr.As(err)
			if !ok {
				t.Fatalf("expected structured business error, got %T", err)
			}
			if messageErr.RuntimeCode() != testCase.runtimeCode {
				t.Fatalf("expected runtime code %q, got %q", testCase.runtimeCode, messageErr.RuntimeCode())
			}
			if messageErr.MessageKey() != testCase.messageKey {
				t.Fatalf("expected message key %q, got %q", testCase.messageKey, messageErr.MessageKey())
			}
		})
	}
}

// TestRoomEnumValidation verifies room type and status named enums only accept
// dictionary-backed values.
func TestRoomEnumValidation(t *testing.T) {
	validStatuses := []int{RoomStatusIdle, RoomStatusLive, RoomStatusDisabled}
	for _, status := range validStatuses {
		if !RoomStatus(status) {
			t.Fatalf("room status %d should be valid", status)
		}
	}
	if RoomStatus(9) {
		t.Fatal("room status 9 should be invalid")
	}

	validTypes := []int{RoomTypeGathering, RoomTypeEvent, RoomTypeOther}
	for _, roomType := range validTypes {
		if !RoomType(roomType) {
			t.Fatalf("room type %d should be valid", roomType)
		}
	}
	if RoomType(4) {
		t.Fatal("room type 4 should be invalid")
	}
}
