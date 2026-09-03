// calendar_ics.go assembles the RFC 5545 calendar stream from the resolved
// feed. Rendering is a pure function over the feed projection so escaping,
// UID stability, all-day events, and UTC conversion stay unit-testable
// without a database, as documented in the OpenSpec change
// add-live-calendar-subscription.

package calendar

import (
	"fmt"
	"strings"
	"time"
)

// ICS output constants. All event times are emitted as UTC with the RFC 5545
// "Z" suffix; X-WR-TIMEZONE advertises the community-facing display zone so
// calendar clients present pending all-day and timed entries consistently.
const (
	// icsProdID identifies the producing application.
	icsProdID = "-//LinaPro//linapro-live-manage//CN"
	// icsTimezone is the community display timezone advertised to clients.
	icsTimezone = "Asia/Shanghai"
	// icsUIDDomain namespaces event UIDs to this plugin.
	icsUIDDomain = "linapro-live-manage"
	// icsUIDPrefix keeps event UIDs traceable to live records.
	icsUIDPrefix = "linapro-live-"
	// icsLineBreak is the RFC 5545 required CRLF line break.
	icsLineBreak = "\r\n"
)

// icsDefaultEventDuration fills DTEND for lives that carry a start time but no
// stored end time: the plugin has no end-time column, so timed entries adopt
// the two-hour community gathering convention instead of inventing per-live
// end times.
const icsDefaultEventDuration = 2 * time.Hour

// RenderICS renders one feed into a complete RFC 5545 calendar document.
// Feeds without entries render a valid empty calendar: the room name is
// deliberately omitted so empty documents rendered for different rejection
// reasons stay indistinguishable and leak no room existence.
func RenderICS(feed *CalendarFeed, watchURL string) string {
	var builder strings.Builder
	builder.WriteString("BEGIN:VCALENDAR")
	builder.WriteString(icsLineBreak)
	builder.WriteString("VERSION:2.0")
	builder.WriteString(icsLineBreak)
	builder.WriteString("PRODID:" + icsProdID)
	builder.WriteString(icsLineBreak)
	builder.WriteString("CALSCALE:GREGORIAN")
	builder.WriteString(icsLineBreak)
	if len(feed.Events) > 0 && feed.RoomName != "" {
		builder.WriteString("X-WR-CALNAME:" + escapeICSText(feed.RoomName))
		builder.WriteString(icsLineBreak)
	}
	builder.WriteString("X-WR-TIMEZONE:" + icsTimezone)
	builder.WriteString(icsLineBreak)

	generatedAt := feed.GeneratedAt
	if generatedAt.IsZero() {
		generatedAt = time.Now()
	}
	stamp := formatICSUTC(generatedAt)
	for _, event := range feed.Events {
		builder.WriteString(renderEvent(event, feed.RoomName, stamp, watchURL))
	}

	builder.WriteString("END:VCALENDAR")
	builder.WriteString(icsLineBreak)
	return builder.String()
}

// renderEvent renders one VEVENT. Lives with a start time become timed entries
// in UTC; lives without one become all-day entries on the live date so no
// concrete time is invented. The UID is keyed by the live record ID, so
// re-subscribing or refreshing never duplicates entries.
func renderEvent(event CalendarEvent, roomName string, stamp string, watchURL string) string {
	var builder strings.Builder
	builder.WriteString("BEGIN:VEVENT")
	builder.WriteString(icsLineBreak)
	builder.WriteString(fmt.Sprintf("UID:%s%d@%s", icsUIDPrefix, event.LiveId, icsUIDDomain))
	builder.WriteString(icsLineBreak)
	builder.WriteString("DTSTAMP:" + stamp)
	builder.WriteString(icsLineBreak)
	builder.WriteString("SUMMARY:" + escapeICSText(event.Title))
	builder.WriteString(icsLineBreak)
	builder.WriteString(renderEventTime(event))
	if description := eventDescription(watchURL, roomName); description != "" {
		builder.WriteString("DESCRIPTION:" + escapeICSText(description))
		builder.WriteString(icsLineBreak)
	}
	builder.WriteString("END:VEVENT")
	builder.WriteString(icsLineBreak)
	return builder.String()
}

// renderEventTime renders the DTSTART/DTEND pair for one event.
func renderEventTime(event CalendarEvent) string {
	var builder strings.Builder
	if event.StartTime == nil {
		builder.WriteString("DTSTART;VALUE=DATE:" + event.LiveDate.Format("20060102"))
		builder.WriteString(icsLineBreak)
		return builder.String()
	}
	start := event.StartTime.UTC()
	end := start.Add(icsDefaultEventDuration)
	builder.WriteString("DTSTART:" + formatICSUTC(start))
	builder.WriteString(icsLineBreak)
	builder.WriteString("DTEND:" + formatICSUTC(end))
	builder.WriteString(icsLineBreak)
	return builder.String()
}

// eventDescription composes the event description: the viewer page link when
// available plus the room name for context.
func eventDescription(watchURL string, roomName string) string {
	var parts []string
	if watchURL != "" {
		parts = append(parts, "Watch: "+watchURL)
	}
	if roomName != "" {
		parts = append(parts, roomName)
	}
	return strings.Join(parts, "\n")
}

// formatICSUTC renders one moment in the RFC 5545 UTC form YYYYMMDDTHHMMSSZ.
func formatICSUTC(value time.Time) string {
	return value.UTC().Format("20060102T150405Z")
}

// escapeICSText escapes text per RFC 5545 section 3.3.11: backslashes,
// semicolons, commas, and line breaks carry a backslash or protocol-level
// escape so titles and room names cannot break the line structure.
func escapeICSText(value string) string {
	replaced := strings.ReplaceAll(value, "\\", "\\\\")
	replaced = strings.ReplaceAll(replaced, ";", "\\;")
	replaced = strings.ReplaceAll(replaced, ",", "\\,")
	replaced = strings.ReplaceAll(replaced, "\r\n", "\\n")
	replaced = strings.ReplaceAll(replaced, "\n", "\\n")
	return replaced
}
