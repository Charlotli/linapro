// calendar_impl.go implements the anonymous ICS schedule resolution: explicit
// tenant validation that degrades to an empty feed instead of an error, a
// point query for the room, and one bounded window query for public lives.
// The visibility boundary matches the play endpoint (is_public plus explicit
// tenant) as documented in the OpenSpec change add-live-calendar-subscription.

package calendar

import (
	"context"
	"strings"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"

	"lina-core/pkg/plugin/capability/tenantcap"
	"lina-plugin-linapro-live-manage/backend/internal/dao"
	"lina-plugin-linapro-live-manage/backend/internal/model/entity"
)

// calendar enumeration values mirror the plugin-wide dictionaries. They repeat
// the values owned by the sibling service packages because this package stays
// independent of admin-facing services.
const (
	// calendarVisibilityPublic matches the plugin_live_public value 1.
	calendarVisibilityPublic = 1
	// calendarPlatformTenantID is the platform tenant owning single-tenant data.
	calendarPlatformTenantID = 0
)

// tenantStatusActive mirrors the provider-owned tenant lifecycle status that
// allows normal tenant access. The host fallback and the tenant-core provider
// both store this value for active tenants.
const tenantStatusActive = "active"

// Subscription window bounds relative to the generation moment: finished
// lives stay visible for seven days so subscribers can jump back into the
// replay, and ninety days ahead cover a quarterly meeting schedule.
const (
	// calendarWindowPastDays is how far back the window reaches.
	calendarWindowPastDays = 7
	// calendarWindowFutureDays is how far ahead the window reaches.
	calendarWindowFutureDays = 90
)

// Feed resolves the public live schedule of one room inside the window.
func (s *serviceImpl) Feed(ctx context.Context, in FeedInput) (*CalendarFeed, error) {
	tenantID, ok := s.resolveTenantID(ctx, in.TenantId)
	if !ok {
		return emptyFeed(), nil
	}

	room, err := loadRoom(ctx, tenantID, in.RoomCode)
	if err != nil {
		return nil, err
	}
	if room == nil {
		return emptyFeed(), nil
	}

	lives, err := loadWindowLives(ctx, tenantID, room.Id)
	if err != nil {
		return nil, err
	}
	if len(lives) == 0 {
		return emptyFeed(), nil
	}

	return buildFeed(room, lives), nil
}

// resolveTenantID maps the explicit tenant parameter onto the query tenant. It
// forces the platform tenant when the tenant capability is disabled, requires
// the parameter when enabled, and validates existence plus lifecycle status.
// Every rejected case reports ok=false so the caller renders the same empty
// feed instead of an error, keeping rejection reasons undisclosed.
func (s *serviceImpl) resolveTenantID(ctx context.Context, tenantID *int) (int, bool) {
	if s.tenantSvc == nil || !s.tenantSvc.Available(ctx) {
		// Without the tenant capability every plugin row lives under the
		// platform tenant, so the parameter is ignored by design.
		return calendarPlatformTenantID, true
	}
	if tenantID == nil || *tenantID < 0 {
		return 0, false
	}
	info, err := s.tenantSvc.Directory().Get(ctx, tenantcap.TenantID(*tenantID))
	if err != nil {
		return 0, false
	}
	if info == nil || info.Status != tenantStatusActive {
		return 0, false
	}
	return *tenantID, true
}

// loadRoom resolves one tenant-scoped live room by its external code. The
// (tenant_id, room_code) unique index turns this into a point query.
func loadRoom(ctx context.Context, tenantID int, roomCode string) (*entity.Room, error) {
	cols := dao.Room.Columns()
	var room *entity.Room
	err := dao.Room.Ctx(ctx).
		Fields(cols.Id, cols.RoomCode, cols.RoomName).
		Where(cols.TenantId, tenantID).
		Where(cols.RoomCode, strings.TrimSpace(roomCode)).
		Scan(&room)
	if err != nil {
		return nil, gerror.Wrap(err, "load live room for calendar subscription")
	}
	return room, nil
}

// loadWindowLives loads the room's public lives inside the subscription
// window. The (tenant_id, room_id) index scopes the scan to one room, so the
// row count stays linear in that room's own schedule and bounded by the
// window; no pagination is needed by design.
func loadWindowLives(ctx context.Context, tenantID int, roomID int64) ([]*entity.Live, error) {
	windowStart, windowEnd := calendarWindow(time.Now())

	cols := dao.Live.Columns()
	var lives []*entity.Live
	err := dao.Live.Ctx(ctx).
		Fields(cols.Id, cols.Title, cols.LiveDate, cols.StartTime).
		Where(cols.TenantId, tenantID).
		Where(cols.RoomId, roomID).
		Where(cols.IsPublic, calendarVisibilityPublic).
		WhereGTE(cols.LiveDate, windowStart.Format(time.DateOnly)).
		WhereLTE(cols.LiveDate, windowEnd.Format(time.DateOnly)+" 23:59:59").
		OrderAsc(cols.LiveDate).
		Scan(&lives)
	if err != nil {
		return nil, gerror.Wrap(err, "load public lives for calendar subscription")
	}
	return lives, nil
}

// calendarWindow computes the subscription window for the generation moment:
// the past bound reaches back calendarWindowPastDays days and the future bound
// reaches ahead calendarWindowFutureDays days. Both bounds are date-accurate
// so DATE column comparisons stay inclusive.
func calendarWindow(now time.Time) (time.Time, time.Time) {
	start := now.AddDate(0, 0, -calendarWindowPastDays)
	end := now.AddDate(0, 0, calendarWindowFutureDays)
	return truncateDate(start), truncateDate(end)
}

// truncateDate strips the time-of-day, keeping the location of the input.
func truncateDate(value time.Time) time.Time {
	return time.Date(value.Year(), value.Month(), value.Day(), 0, 0, 0, 0, value.Location())
}

// emptyFeed returns the canonical empty feed shared by every rejected or
// contentless case so responses stay indistinguishable.
func emptyFeed() *CalendarFeed {
	return &CalendarFeed{Events: []CalendarEvent{}}
}

// buildFeed projects the room and public lives into the renderable feed. The
// room name is only attached when entries exist so empty feeds rendered for
// different reasons stay byte-identical.
func buildFeed(room *entity.Room, lives []*entity.Live) *CalendarFeed {
	events := make([]CalendarEvent, 0, len(lives))
	for _, live := range lives {
		events = append(events, CalendarEvent{
			LiveId:    live.Id,
			Title:     live.Title,
			LiveDate:  live.LiveDate,
			StartTime: live.StartTime,
		})
	}
	return &CalendarFeed{
		RoomName:    room.RoomName,
		GeneratedAt: time.Now(),
		Events:      events,
	}
}
