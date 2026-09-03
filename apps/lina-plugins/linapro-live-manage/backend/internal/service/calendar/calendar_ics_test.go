// This file verifies the ICS rendering pure functions: calendar document
// structure, UID stability, timed and all-day event rendering, RFC 5545 text
// escaping, room-name suppression on empty feeds, and the subscription window
// computation. Tests are self-contained pure-function checks, so no database
// or request context is required.

package calendar

import (
	"strings"
	"testing"
	"time"
)

// TestCalendarWindowBounds verifies the subscription window reaches back seven
// days and ahead ninety days, both truncated to whole dates so DATE column
// comparisons stay inclusive.
func TestCalendarWindowBounds(t *testing.T) {
	now := time.Date(2026, 9, 3, 15, 42, 17, 0, time.UTC)
	start, end := calendarWindow(now)

	expectedStart := time.Date(2026, 8, 27, 0, 0, 0, 0, time.UTC)
	expectedEnd := time.Date(2026, 12, 2, 0, 0, 0, 0, time.UTC)
	if !start.Equal(expectedStart) {
		t.Fatalf("expected window start %v, got %v", expectedStart, start)
	}
	if !end.Equal(expectedEnd) {
		t.Fatalf("expected window end %v, got %v", expectedEnd, end)
	}
}

// TestRenderICSDocumentStructure verifies the calendar envelope, the room
// display name, and per-event UID stability for a feed with entries.
func TestRenderICSDocumentStructure(t *testing.T) {
	startTime := time.Date(2026, 9, 6, 10, 0, 0, 0, time.UTC)
	feed := &CalendarFeed{
		RoomName:    "Main hall",
		GeneratedAt: time.Date(2026, 9, 3, 8, 0, 0, 0, time.UTC),
		Events: []CalendarEvent{
			{LiveId: 42, Title: "Sunday service", LiveDate: startTime, StartTime: &startTime},
		},
	}

	document := RenderICS(feed, "https://example.com/x/linapro-live-manage/h5?room=MAIN-HALL")
	lines := strings.Split(document, "\r\n")
	if lines[0] != "BEGIN:VCALENDAR" {
		t.Fatalf("expected calendar begin, got %q", lines[0])
	}
	if !containsPrefix(lines, "VERSION:2.0") ||
		!containsPrefix(lines, "PRODID:"+icsProdID) ||
		!containsPrefix(lines, "CALSCALE:GREGORIAN") ||
		!containsPrefix(lines, "X-WR-TIMEZONE:"+icsTimezone) {
		t.Fatalf("missing calendar headers in document:\n%s", document)
	}
	if !containsPrefix(lines, "X-WR-CALNAME:Main hall") {
		t.Fatalf("expected room display name header, got:\n%s", document)
	}
	if !containsPrefix(lines, "UID:linapro-live-42@linapro-live-manage") {
		t.Fatalf("expected stable live-keyed UID, got:\n%s", document)
	}
	if !containsPrefix(lines, "DTSTAMP:20260903T080000Z") {
		t.Fatalf("expected UTC generation stamp, got:\n%s", document)
	}
	if lines[len(lines)-1] != "" || !strings.HasSuffix(document, "END:VCALENDAR\r\n") {
		t.Fatalf("expected CRLF-terminated calendar end, got %q", lines[len(lines)-1])
	}
}

// TestRenderICSTimedEvent verifies the UTC start/end pair and the watch link
// description for a live with a recorded start time. The end time adopts the
// default two-hour gathering duration because the schema stores no end time.
func TestRenderICSTimedEvent(t *testing.T) {
	startTime := time.Date(2026, 9, 6, 10, 0, 0, 0, time.UTC)
	feed := &CalendarFeed{
		RoomName:    "Main hall",
		GeneratedAt: time.Date(2026, 9, 3, 8, 0, 0, 0, time.UTC),
		Events: []CalendarEvent{
			{LiveId: 42, Title: "Sunday service", LiveDate: startTime, StartTime: &startTime},
		},
	}

	document := RenderICS(feed, "https://example.com/x/linapro-live-manage/h5?room=MAIN-HALL")
	lines := strings.Split(document, "\r\n")
	if !containsPrefix(lines, "DTSTART:20260906T100000Z") {
		t.Fatalf("expected UTC start, got:\n%s", document)
	}
	if !containsPrefix(lines, "DTEND:20260906T120000Z") {
		t.Fatalf("expected default two-hour UTC end, got:\n%s", document)
	}
	if !containsPrefix(lines, `DESCRIPTION:Watch: https://example.com/x/linapro-live-manage/h5?room=MAIN-HALL\nMain hall`) {
		t.Fatalf("expected watch link and room name description, got:\n%s", document)
	}
}

// TestRenderICSAllDayEvent verifies that lives without a start time render as
// all-day entries on the live date, without inventing concrete times.
func TestRenderICSAllDayEvent(t *testing.T) {
	liveDate := time.Date(2026, 9, 6, 0, 0, 0, 0, time.UTC)
	feed := &CalendarFeed{
		RoomName:    "Main hall",
		GeneratedAt: time.Date(2026, 9, 3, 8, 0, 0, 0, time.UTC),
		Events: []CalendarEvent{
			{LiveId: 7, Title: "Prayer meeting", LiveDate: liveDate},
		},
	}

	document := RenderICS(feed, "")
	lines := strings.Split(document, "\r\n")
	if !containsPrefix(lines, "DTSTART;VALUE=DATE:20260906") {
		t.Fatalf("expected all-day start on the live date, got:\n%s", document)
	}
	if containsPrefix(lines, "DTEND:") {
		t.Fatalf("all-day entries must not carry a DTEND, got:\n%s", document)
	}
	if !containsPrefix(lines, "DESCRIPTION:Main hall") {
		t.Fatalf("expected room-name-only description without a watch link, got:\n%s", document)
	}
}

// TestRenderICSEmptyFeedHidesRoomName verifies every contentless case renders
// the same valid empty calendar without the room display name, so empty feeds
// rendered for different reasons stay indistinguishable and leak no room
// existence.
func TestRenderICSEmptyFeedHidesRoomName(t *testing.T) {
	missingRoomFeed := &CalendarFeed{
		GeneratedAt: time.Date(2026, 9, 3, 8, 0, 0, 0, time.UTC),
		Events:      []CalendarEvent{},
	}
	knownRoomFeed := &CalendarFeed{
		RoomName:    "Main hall",
		GeneratedAt: time.Date(2026, 9, 3, 8, 0, 0, 0, time.UTC),
		Events:      []CalendarEvent{},
	}

	missingRoomDocument := RenderICS(missingRoomFeed, "https://example.com/h5")
	knownRoomDocument := RenderICS(knownRoomFeed, "https://example.com/h5")
	if missingRoomDocument != knownRoomDocument {
		t.Fatal("empty feeds must render identically regardless of room resolution")
	}
	lines := strings.Split(missingRoomDocument, "\r\n")
	if containsPrefix(lines, "X-WR-CALNAME:") {
		t.Fatalf("empty feeds must not carry the room display name, got:\n%s", missingRoomDocument)
	}
	if containsPrefix(lines, "BEGIN:VEVENT") {
		t.Fatalf("empty feeds must not carry events, got:\n%s", missingRoomDocument)
	}
	if !strings.HasSuffix(missingRoomDocument, "END:VCALENDAR\r\n") {
		t.Fatalf("empty feeds must stay valid calendars, got:\n%s", missingRoomDocument)
	}
}

// TestEscapeICSText verifies the RFC 5545 escaping of backslashes, list
// separators, and line breaks so text values cannot break the document
// structure.
func TestEscapeICSText(t *testing.T) {
	escaped := escapeICSText(`标题; 主题, 第一部分\反斜杠
第二行`)
	expected := `标题\; 主题\, 第一部分\\反斜杠\n第二行`
	if escaped != expected {
		t.Fatalf("expected %q, got %q", expected, escaped)
	}
}

// TestRenderEventUIDStability verifies re-rendering the same live produces the
// same UID so calendar clients never duplicate entries on refresh.
func TestRenderEventUIDStability(t *testing.T) {
	event := CalendarEvent{LiveId: 42, Title: "Sunday service"}
	first := renderEvent(event, "", "20260903T080000Z", "")
	second := renderEvent(event, "", "20260910T090000Z", "")

	firstUID := findLine(strings.Split(first, "\r\n"), "UID:")
	secondUID := findLine(strings.Split(second, "\r\n"), "UID:")
	if firstUID == "" || firstUID != secondUID {
		t.Fatalf("expected stable UID %q across renders, got %q and %q", "linapro-live-42@linapro-live-manage", firstUID, secondUID)
	}
}

// findLine returns the first line starting with the given prefix, or empty.
func findLine(lines []string, prefix string) string {
	for _, line := range lines {
		if strings.HasPrefix(line, prefix) {
			return line
		}
	}
	return ""
}

// containsPrefix reports whether any line starts with the given prefix.
func containsPrefix(lines []string, prefix string) bool {
	return findLine(lines, prefix) != ""
}
