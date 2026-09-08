// bible_impl_test.go covers the pure request validation of the bible
// reading service: out-of-range volume and chapter inputs must answer with
// empty payloads without touching the database. Catalogue-dependent paths
// are exercised by the E2E suite.

package bible

import (
	"context"
	"testing"
)

func TestChapterRejectsInvalidInputWithoutDatabase(t *testing.T) {
	svc := &serviceImpl{}
	cases := []struct {
		name     string
		volumeSn int
		chapter  int
	}{
		{"zero volume", 0, 1},
		{"negative volume", -1, 1},
		{"zero chapter", 1, 0},
		{"negative chapter", 1, -2},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := svc.Chapter(context.Background(), ChapterInput{
				VolumeSn: tc.volumeSn,
				Chapter:  tc.chapter,
			})
			if err != nil {
				t.Fatalf("Chapter(%d, %d) returned error: %v", tc.volumeSn, tc.chapter, err)
			}
			if result == nil {
				t.Fatal("expected empty result payload")
			}
			if len(result.Verses) != 0 {
				t.Fatalf("expected empty verses, got %d", len(result.Verses))
			}
			if result.Book != "" {
				t.Fatalf("expected empty book name, got %q", result.Book)
			}
		})
	}
}
