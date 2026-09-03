// This file verifies the calendar feed tenant-resolution boundary: every
// rejected tenant case degrades to the empty-feed path instead of an error,
// keeping rejection reasons undisclosed per the OpenSpec change
// add-live-calendar-subscription. Tests are self-contained and rely only on a
// stub tenant service, so no database is required.

package calendar

import (
	"context"
	"testing"
	"time"

	"lina-core/pkg/bizerr"
	"lina-core/pkg/plugin/capability/capmodel"
	"lina-core/pkg/plugin/capability/tenantcap"
	"lina-plugin-linapro-live-manage/backend/internal/model/entity"
)

// stubTenantDirectory implements the plugin-visible tenant directory contract
// for tenant-resolution tests.
type stubTenantDirectory struct {
	tenants map[tenantcap.TenantID]*tenantcap.TenantInfo
}

// Get returns the stubbed tenant projection or a not-found error.
func (s stubTenantDirectory) Get(_ context.Context, tenantID tenantcap.TenantID) (*tenantcap.TenantInfo, error) {
	info, ok := s.tenants[tenantID]
	if !ok {
		return nil, bizerr.NewCode(tenantcap.CodeTenantForbidden, bizerr.P("tenantId", int(tenantID)))
	}
	return info, nil
}

// BatchGet is unused by these tests.
func (s stubTenantDirectory) BatchGet(context.Context, []tenantcap.TenantID) (*capmodel.BatchResult[*tenantcap.TenantInfo, tenantcap.TenantID], error) {
	return nil, nil
}

// List is unused by these tests.
func (s stubTenantDirectory) List(context.Context, tenantcap.ListInput) (*capmodel.PageResult[*tenantcap.TenantInfo], error) {
	return nil, nil
}

// EnsureVisible is unused by these tests.
func (s stubTenantDirectory) EnsureVisible(context.Context, []tenantcap.TenantID) error {
	return nil
}

// stubTenantService implements the plugin-visible tenant capability with a
// controllable availability flag and directory.
type stubTenantService struct {
	available bool
	directory stubTenantDirectory
}

// Available reports the stubbed availability.
func (s stubTenantService) Available(context.Context) bool {
	return s.available
}

// Status returns an empty capability status.
func (s stubTenantService) Status(context.Context) capmodel.CapabilityStatus {
	return capmodel.CapabilityStatus{}
}

// Context returns no context service.
func (s stubTenantService) Context() tenantcap.ContextService {
	return nil
}

// Directory returns the stubbed directory.
func (s stubTenantService) Directory() tenantcap.DirectoryService {
	return s.directory
}

// Membership returns no membership service.
func (s stubTenantService) Membership() tenantcap.MembershipService {
	return nil
}

// Plugins returns no plugin governance service.
func (s stubTenantService) Plugins() tenantcap.PluginService {
	return nil
}

// Filter returns no filter service.
func (s stubTenantService) Filter() tenantcap.FilterService {
	return nil
}

// TestResolveTenantIDWithoutCapability verifies the platform tenant is forced
// when the tenant capability is unavailable, regardless of the parameter.
func TestResolveTenantIDWithoutCapability(t *testing.T) {
	svc := New(stubTenantService{available: false}).(*serviceImpl)
	missing := 0
	tenantID, ok := svc.resolveTenantID(context.Background(), nil)
	if !ok || tenantID != calendarPlatformTenantID {
		t.Fatalf("expected platform tenant without capability, got %d ok=%v", tenantID, ok)
	}
	tenantID, ok = svc.resolveTenantID(context.Background(), &missing)
	if !ok || tenantID != calendarPlatformTenantID {
		t.Fatalf("expected explicit platform tenant without capability, got %d ok=%v", tenantID, ok)
	}
}

// TestResolveTenantIDRejectsInvalidCases verifies every rejected tenant case
// reports the empty-feed path instead of an error.
func TestResolveTenantIDRejectsInvalidCases(t *testing.T) {
	svc := New(stubTenantService{
		available: true,
		directory: stubTenantDirectory{
			tenants: map[tenantcap.TenantID]*tenantcap.TenantInfo{
				1: {ID: 1, Status: tenantStatusActive},
				2: {ID: 2, Status: "suspended"},
			},
		},
	}).(*serviceImpl)

	if _, ok := svc.resolveTenantID(context.Background(), nil); ok {
		t.Fatal("missing tenant parameter must take the empty-feed path when the capability is enabled")
	}
	negative := -1
	if _, ok := svc.resolveTenantID(context.Background(), &negative); ok {
		t.Fatal("negative tenant IDs must take the empty-feed path")
	}
	suspended := 2
	if _, ok := svc.resolveTenantID(context.Background(), &suspended); ok {
		t.Fatal("suspended tenants must take the empty-feed path")
	}
	unknown := 99
	if _, ok := svc.resolveTenantID(context.Background(), &unknown); ok {
		t.Fatal("unknown tenants must take the empty-feed path")
	}
	active := 1
	tenantID, ok := svc.resolveTenantID(context.Background(), &active)
	if !ok || tenantID != active {
		t.Fatalf("expected active tenant %d to resolve, got %d ok=%v", active, tenantID, ok)
	}
}

// TestBuildFeedProjectsEvents verifies the feed projection carries the room
// name and one event per live with stable identities.
func TestBuildFeedProjectsEvents(t *testing.T) {
	startTime := time.Date(2026, 9, 6, 10, 0, 0, 0, time.UTC)
	room := &entity.Room{Id: 7, RoomCode: "MAIN-HALL", RoomName: "Main hall"}
	lives := []*entity.Live{
		{Id: 42, Title: "Sunday service", LiveDate: startTime, StartTime: &startTime},
		{Id: 43, Title: "Prayer meeting", LiveDate: startTime},
	}

	feed := buildFeed(room, lives)
	if feed.RoomName != "Main hall" {
		t.Fatalf("expected room name projection, got %q", feed.RoomName)
	}
	if len(feed.Events) != 2 {
		t.Fatalf("expected 2 events, got %d", len(feed.Events))
	}
	if feed.Events[0].LiveId != 42 || feed.Events[0].StartTime == nil {
		t.Fatalf("expected timed event projection, got %+v", feed.Events[0])
	}
	if feed.Events[1].StartTime != nil {
		t.Fatalf("expected all-day event projection, got %+v", feed.Events[1])
	}
}
