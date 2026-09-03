// This file verifies the anonymous viewer play selection logic: business
// error metadata, play-URL visibility by live state, presentability rules,
// state priority ordering, play-info projection withholding, and the
// tenant-capability-aware tenant resolution. Tests are self-contained and
// rely only on pure functions plus a stub tenant service, so no database or
// request context is required.

package play

import (
	"context"
	"testing"
	"time"

	"lina-core/pkg/bizerr"
	"lina-core/pkg/plugin/capability/capmodel"
	"lina-core/pkg/plugin/capability/tenantcap"
	"lina-plugin-linapro-live-manage/backend/internal/model/entity"
	"lina-plugin-linapro-live-manage/backend/internal/service/liveroom"
)

// TestPlayBusinessErrorMetadata verifies play errors expose stable runtime
// codes and i18n keys instead of fixed-language text.
func TestPlayBusinessErrorMetadata(t *testing.T) {
	testCases := []struct {
		name        string
		code        *bizerr.Code
		runtimeCode string
		messageKey  string
	}{
		{
			name:        "play not found",
			code:        CodePlayNotFound,
			runtimeCode: "LIVE_MANAGE_PLAY_NOT_FOUND",
			messageKey:  "error.live.manage.play.not.found",
		},
		{
			name:        "play tenant required",
			code:        CodePlayTenantRequired,
			runtimeCode: "LIVE_MANAGE_PLAY_TENANT_REQUIRED",
			messageKey:  "error.live.manage.play.tenant.required",
		},
		{
			name:        "play tenant invalid",
			code:        CodePlayTenantInvalid,
			runtimeCode: "LIVE_MANAGE_PLAY_TENANT_INVALID",
			messageKey:  "error.live.manage.play.tenant.invalid",
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

// TestPlayURLVisible verifies the play URL stays hidden for not-started
// lives and appears for ongoing and finished lives.
func TestPlayURLVisible(t *testing.T) {
	if playURLVisible(playLiveStateNotStarted) {
		t.Fatal("not-started lives must not expose the play URL")
	}
	if !playURLVisible(playLiveStateOngoing) {
		t.Fatal("ongoing lives must expose the play URL")
	}
	if !playURLVisible(playLiveStateFinished) {
		t.Fatal("finished lives must expose the play URL for replay")
	}
}

// TestPresentable verifies the room-status gate: ongoing lives require the
// room to be live, while replay and preview states remain presentable.
func TestPresentable(t *testing.T) {
	liveRoom := &entity.Room{Status: liveroom.RoomStatusLive}
	idleRoom := &entity.Room{Status: liveroom.RoomStatusIdle}
	disabledRoom := &entity.Room{Status: liveroom.RoomStatusDisabled}

	if !presentable(playLiveStateOngoing, liveRoom) {
		t.Fatal("ongoing live must be presentable when the room is live")
	}
	if presentable(playLiveStateOngoing, idleRoom) {
		t.Fatal("ongoing live must not be presentable when the room is idle")
	}
	if presentable(playLiveStateOngoing, disabledRoom) {
		t.Fatal("ongoing live must not be presentable when the room is disabled")
	}
	if presentable(playLiveStateOngoing, nil) {
		t.Fatal("ongoing live must not be presentable without a room")
	}
	if !presentable(playLiveStateFinished, disabledRoom) {
		t.Fatal("replay must remain presentable regardless of room status")
	}
	if !presentable(playLiveStateNotStarted, idleRoom) {
		t.Fatal("preview must remain presentable regardless of room status")
	}
}

// TestPlayStateOrder verifies the candidate priority: ongoing first, replay
// second, preview last.
func TestPlayStateOrder(t *testing.T) {
	expected := []int{playLiveStateOngoing, playLiveStateFinished, playLiveStateNotStarted}
	if len(playStateOrder) != len(expected) {
		t.Fatalf("expected %d states, got %d", len(expected), len(playStateOrder))
	}
	for index, state := range expected {
		if playStateOrder[index] != state {
			t.Fatalf("expected state %d at position %d, got %d", state, index, playStateOrder[index])
		}
	}
}

// TestBuildPlayInfoWithholdsFields verifies the projection drops the play URL
// for not-started lives and keeps administrative fields out of the payload.
func TestBuildPlayInfoWithholdsFields(t *testing.T) {
	startTime := time.Date(2026, 4, 26, 10, 0, 0, 0, time.UTC)
	room := &entity.Room{Id: 7, RoomCode: "MAIN-HALL", RoomName: "Main hall", Status: liveroom.RoomStatusLive}
	live := &entity.Live{
		Id:        42,
		RoomId:    room.Id,
		Title:     "Sunday service",
		LiveUrl:   "https://example.com/hls/sunday.m3u8",
		PushUrl:   "rtmp://example.com/push/sunday",
		LiveDate:  startTime,
		StartTime: &startTime,
		State:     playLiveStateOngoing,
		IsPublic:  playVisibilityPublic,
	}

	info := buildPlayInfo(room, live)
	if info.LiveUrl != live.LiveUrl {
		t.Fatalf("expected ongoing play URL %q, got %q", live.LiveUrl, info.LiveUrl)
	}
	if info.RoomCode != room.RoomCode || info.RoomName != room.RoomName {
		t.Fatal("expected room identity to be projected")
	}
	if info.State != playLiveStateOngoing {
		t.Fatalf("expected ongoing state, got %d", info.State)
	}

	live.State = playLiveStateNotStarted
	info = buildPlayInfo(room, live)
	if info.LiveUrl != "" {
		t.Fatalf("expected empty play URL for not-started live, got %q", info.LiveUrl)
	}
	// The projection struct has no push-URL field by construction; assert the
	// generated API DTO also omits it via the contract test in the api package.
	if info.State != playLiveStateNotStarted {
		t.Fatalf("expected not-started state, got %d", info.State)
	}
}

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
	tenantID, err := svc.resolveTenantID(context.Background(), nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if tenantID != playPlatformTenantID {
		t.Fatalf("expected platform tenant, got %d", tenantID)
	}
	tenantID, err = svc.resolveTenantID(context.Background(), &missing)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if tenantID != playPlatformTenantID {
		t.Fatalf("expected platform tenant, got %d", tenantID)
	}
}

// TestResolveTenantIDRequiresParameter verifies the required-tenant error when
// the capability is enabled and the parameter is absent.
func TestResolveTenantIDRequiresParameter(t *testing.T) {
	svc := New(stubTenantService{available: true}).(*serviceImpl)
	_, err := svc.resolveTenantID(context.Background(), nil)
	assertPlayCode(t, err, CodePlayTenantRequired)
}

// TestResolveTenantIDValidatesLifecycle verifies existence and lifecycle
// validation of the explicit tenant parameter.
func TestResolveTenantIDValidatesLifecycle(t *testing.T) {
	svc := New(stubTenantService{
		available: true,
		directory: stubTenantDirectory{
			tenants: map[tenantcap.TenantID]*tenantcap.TenantInfo{
				1: {ID: 1, Status: tenantStatusActive},
				2: {ID: 2, Status: "suspended"},
			},
		},
	}).(*serviceImpl)

	active := 1
	tenantID, err := svc.resolveTenantID(context.Background(), &active)
	if err != nil {
		t.Fatalf("expected active tenant to resolve, got %v", err)
	}
	if tenantID != active {
		t.Fatalf("expected tenant %d, got %d", active, tenantID)
	}

	suspended := 2
	_, err = svc.resolveTenantID(context.Background(), &suspended)
	assertPlayCode(t, err, CodePlayTenantInvalid)

	unknown := 99
	_, err = svc.resolveTenantID(context.Background(), &unknown)
	assertPlayCode(t, err, CodePlayTenantInvalid)

	negative := -1
	_, err = svc.resolveTenantID(context.Background(), &negative)
	assertPlayCode(t, err, CodePlayTenantInvalid)
}

// assertPlayCode reports whether the error carries the expected play code.
func assertPlayCode(t *testing.T, err error, code *bizerr.Code) {
	t.Helper()
	if err == nil {
		t.Fatal("expected an error, got none")
	}
	messageErr, ok := bizerr.As(err)
	if !ok {
		t.Fatalf("expected structured business error, got %T: %v", err, err)
	}
	if messageErr.RuntimeCode() != code.RuntimeCode() {
		t.Fatalf("expected runtime code %q, got %q", code.RuntimeCode(), messageErr.RuntimeCode())
	}
}
