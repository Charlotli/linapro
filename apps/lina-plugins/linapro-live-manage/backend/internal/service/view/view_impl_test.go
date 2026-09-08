// view_impl_test.go covers the pure-logic guards of the watch-statistics
// service: session-key validation runs before any dependency is touched, so
// these tests stay self-contained without a database or play service.

package view

import (
	"context"
	"strings"
	"testing"
)

// TestHeartbeatRejectsInvalidSessionKey verifies the service-level guard for
// the session key: empty keys and keys longer than the column bound are
// rejected with the dedicated business error before the play service is
// consulted.
func TestHeartbeatRejectsInvalidSessionKey(t *testing.T) {
	service := &serviceImpl{}
	cases := []struct {
		name       string
		sessionKey string
	}{
		{name: "empty", sessionKey: "   "},
		{name: "too long", sessionKey: strings.Repeat("a", maxSessionKeyLength+1)},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			stats, err := service.Heartbeat(context.Background(), HeartbeatInput{
				RoomCode:   "ANY-ROOM",
				SessionKey: testCase.sessionKey,
			})
			if err == nil {
				t.Fatalf("expected session-key error for %q", testCase.sessionKey)
			}
			if stats != nil {
				t.Fatal("expected nil stats on session-key error")
			}
		})
	}
}

// TestMaxSessionKeyLengthMatchesColumnBound documents the accepted bounds: a
// key of exactly the column width is valid input shape.
func TestMaxSessionKeyLengthMatchesColumnBound(t *testing.T) {
	if maxSessionKeyLength != 64 {
		t.Fatalf("expected session key bound 64 to mirror session_key varchar(64), got %d", maxSessionKeyLength)
	}
}
