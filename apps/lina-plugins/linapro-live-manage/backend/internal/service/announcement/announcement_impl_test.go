// announcement_impl_test.go covers the pure helpers of the announcement
// service: paging normalization and the play tenant-error passthrough
// predicate. Database-dependent paths are exercised by the E2E suite.

package announcement

import (
	"errors"
	"testing"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"

	"lina-core/pkg/bizerr"
	"lina-plugin-linapro-live-manage/backend/internal/service/play"
)

func TestNormalizePage(t *testing.T) {
	cases := []struct {
		name             string
		pageNum          int
		pageSize         int
		wantPageNum      int
		wantPageSize     int
	}{
		{"defaults on zero", 0, 0, 1, PageDefaultSize},
		{"defaults on negative", -3, -1, 1, PageDefaultSize},
		{"keeps valid input", 2, 25, 2, 25},
		{"clamps oversized page size", 1, 500, 1, PageMaxSize},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			gotPageNum, gotPageSize := normalizePage(tc.pageNum, tc.pageSize)
			if gotPageNum != tc.wantPageNum || gotPageSize != tc.wantPageSize {
				t.Fatalf("normalizePage(%d, %d) = (%d, %d), want (%d, %d)",
					tc.pageNum, tc.pageSize, gotPageNum, gotPageSize, tc.wantPageNum, tc.wantPageSize)
			}
		})
	}
}

func TestIsTenantErrorPassesPlayTenantContractThrough(t *testing.T) {
	if !isTenantError(bizerr.NewCode(play.CodePlayTenantRequired)) {
		t.Fatal("expected CodePlayTenantRequired to pass through")
	}
	if !isTenantError(bizerr.NewCode(play.CodePlayTenantInvalid)) {
		t.Fatal("expected CodePlayTenantInvalid to pass through")
	}
	if isTenantError(errors.New("database down")) {
		t.Fatal("expected plain infrastructure errors to stay swallowed")
	}
	if isTenantError(bizerr.NewCode(CodeAnnouncementNotFound)) {
		t.Fatal("expected announcement errors to stay swallowed")
	}
}

func TestIsTenantErrorIgnoresWrappedNonTenantErrors(t *testing.T) {
	wrapped := gerror.WrapCode(gcode.CodeNotFound, errors.New("missing"), "load")
	if isTenantError(wrapped) {
		t.Fatal("expected wrapped non-tenant errors to stay swallowed")
	}
}
