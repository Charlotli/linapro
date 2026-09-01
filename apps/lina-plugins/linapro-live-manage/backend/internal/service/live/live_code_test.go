// This file verifies linapro-live-manage live-content business error metadata
// and song-list validation rules.

package live

import (
	"testing"

	"lina-core/pkg/bizerr"
)

// TestLiveBusinessErrorMetadata verifies live-content errors expose stable
// runtime codes and i18n keys instead of fixed-language text.
func TestLiveBusinessErrorMetadata(t *testing.T) {
	testCases := []struct {
		name        string
		code        *bizerr.Code
		runtimeCode string
		messageKey  string
	}{
		{
			name:        "live not found",
			code:        CodeLiveNotFound,
			runtimeCode: "LIVE_MANAGE_LIVE_NOT_FOUND",
			messageKey:  "error.live.manage.live.not.found",
		},
		{
			name:        "room unavailable",
			code:        CodeRoomUnavailable,
			runtimeCode: "LIVE_MANAGE_ROOM_UNAVAILABLE",
			messageKey:  "error.live.manage.room.unavailable",
		},
		{
			name:        "song list invalid",
			code:        CodeSongListInvalid,
			runtimeCode: "LIVE_MANAGE_SONG_LIST_INVALID",
			messageKey:  "error.live.manage.song.list.invalid",
		},
		{
			name:        "delete required",
			code:        CodeLiveDeleteRequired,
			runtimeCode: "LIVE_MANAGE_LIVE_DELETE_REQUIRED",
			messageKey:  "error.live.manage.live.delete.required",
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

// TestNormalizeSongList verifies the song-list payload stays empty when unset
// and only accepts valid JSON arrays otherwise.
func TestNormalizeSongList(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{name: "empty stays empty", input: "", expected: "", wantErr: false},
		{name: "blank stays empty", input: "   ", expected: "", wantErr: false},
		{
			name:     "valid array",
			input:    `[{"name":"Song A","singer":"Alice","order":1}]`,
			expected: `[{"name":"Song A","singer":"Alice","order":1}]`,
			wantErr:  false,
		},
		{name: "object rejected", input: `{"name":"Song A"}`, wantErr: true},
		{name: "plain text rejected", input: `Song A`, wantErr: true},
		{name: "broken json rejected", input: `[{"name":}`, wantErr: true},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			result, err := normalizeSongList(testCase.input)
			if testCase.wantErr {
				if err == nil {
					t.Fatalf("expected error for input %q", testCase.input)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error for input %q: %v", testCase.input, err)
			}
			if result != testCase.expected {
				t.Fatalf("expected %q, got %q", testCase.expected, result)
			}
		})
	}
}

// TestParseLiveDate verifies date-only parsing defaults to today for empty
// input and rejects malformed values.
func TestParseLiveDate(t *testing.T) {
	parsed, err := parseLiveDate("2026-04-26")
	if err != nil {
		t.Fatalf("expected valid parse, got %v", err)
	}
	if parsed.Format("2006-01-02") != "2026-04-26" {
		t.Fatalf("expected 2026-04-26, got %s", parsed.Format("2006-01-02"))
	}

	if _, err = parseLiveDate(""); err != nil {
		t.Fatalf("expected empty input to default to today, got %v", err)
	}

	if _, err = parseLiveDate("2026/04/26"); err == nil {
		t.Fatal("expected malformed date to be rejected")
	}
}

// TestLiveStateTransition verifies the explicit start/stop state machine:
// only not-started lives can start and only ongoing lives can stop.
func TestLiveStateTransition(t *testing.T) {
	testCases := []struct {
		name        string
		current     int
		action      liveStateAction
		wantTarget  int
		wantAllowed bool
	}{
		{name: "start from not-started", current: liveStateNotStarted, action: liveStateActionStart, wantTarget: liveStateOngoing, wantAllowed: true},
		{name: "start from ongoing rejected", current: liveStateOngoing, action: liveStateActionStart, wantAllowed: false},
		{name: "start from finished rejected", current: liveStateFinished, action: liveStateActionStart, wantAllowed: false},
		{name: "stop from ongoing", current: liveStateOngoing, action: liveStateActionStop, wantTarget: liveStateFinished, wantAllowed: true},
		{name: "stop from not-started rejected", current: liveStateNotStarted, action: liveStateActionStop, wantAllowed: false},
		{name: "stop from finished rejected", current: liveStateFinished, action: liveStateActionStop, wantAllowed: false},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			target, allowed := liveStateTransition(testCase.current, testCase.action)
			if allowed != testCase.wantAllowed {
				t.Fatalf("expected allowed=%v, got %v", testCase.wantAllowed, allowed)
			}
			if allowed && target != testCase.wantTarget {
				t.Fatalf("expected target state %d, got %d", testCase.wantTarget, target)
			}
		})
	}
}
