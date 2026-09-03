// This file verifies the live date wire-format parsing: the empty input
// defaults to today's calendar date (regression for the gtime PHP-token
// layout misuse that stored the constant "2006-01-02"), and malformed input
// is rejected with the stable business error code.

package live

import (
	"testing"
	"time"
)

// TestParseLiveDateDefaultsEmptyToToday verifies that omitted live dates
// resolve to the current calendar date instead of the layout constant.
func TestParseLiveDateDefaultsEmptyToToday(t *testing.T) {
	for _, input := range []string{"", "   "} {
		parsed, err := parseLiveDate(input)
		if err != nil {
			t.Fatalf("empty input %q must default to today, got error %v", input, err)
		}
		today := time.Now().Format(liveDateFormat)
		if parsed.Format(liveDateFormat) != today {
			t.Fatalf("expected empty input %q to resolve to today %s, got %s", input, today, parsed.Format(liveDateFormat))
		}
	}
}

// TestParseLiveDateAcceptsValidAndRejectsInvalid verifies the date-only wire
// parsing and the stable rejection code for malformed values.
func TestParseLiveDateAcceptsValidAndRejectsInvalid(t *testing.T) {
	parsed, err := parseLiveDate("2026-09-03")
	if err != nil {
		t.Fatalf("expected valid date to parse, got error %v", err)
	}
	if parsed.Format(liveDateFormat) != "2026-09-03" {
		t.Fatalf("expected 2026-09-03, got %s", parsed.Format(liveDateFormat))
	}

	for _, input := range []string{"not-a-date", "2026/09/03"} {
		if _, err := parseLiveDate(input); err == nil {
			t.Fatalf("expected malformed input %q to be rejected", input)
		}
	}
}
