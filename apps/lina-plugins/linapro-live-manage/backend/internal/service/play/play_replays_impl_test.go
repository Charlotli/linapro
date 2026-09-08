// play_replays_impl_test.go covers the pure pagination normalization of the
// public replay-library query: defaults for omitted values and the hard
// page-size cap.

package play

import "testing"

// TestNormalizeReplayPage verifies the documented pagination contract: page
// defaults to 1, page size defaults to 10, and oversized page sizes clamp to
// the 50 cap.
func TestNormalizeReplayPage(t *testing.T) {
	cases := []struct {
		name             string
		page             int
		pageSize         int
		expectedPage     int
		expectedPageSize int
	}{
		{name: "defaults", page: 0, pageSize: 0, expectedPage: 1, expectedPageSize: 10},
		{name: "negative values", page: -3, pageSize: -1, expectedPage: 1, expectedPageSize: 10},
		{name: "valid values pass through", page: 4, pageSize: 25, expectedPage: 4, expectedPageSize: 25},
		{name: "oversized page size clamps", page: 2, pageSize: 500, expectedPage: 2, expectedPageSize: 50},
		{name: "page size at cap", page: 1, pageSize: 50, expectedPage: 1, expectedPageSize: 50},
	}
	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			page, pageSize := normalizeReplayPage(testCase.page, testCase.pageSize)
			if page != testCase.expectedPage || pageSize != testCase.expectedPageSize {
				t.Fatalf("expected (%d, %d), got (%d, %d)",
					testCase.expectedPage, testCase.expectedPageSize, page, pageSize)
			}
		})
	}
}
