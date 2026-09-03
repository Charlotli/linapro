// Package calendar implements the anonymous public ICS calendar subscription
// service for the linapro-live-manage source plugin. It renders the public
// live schedule of one live-room code as an RFC 5545 calendar stream so
// viewers can subscribe from their system calendar. Like the sibling play
// package, this package deliberately does not read the request business
// context: the tenant is an explicit input, and the authoritative visibility
// boundary is the is_public flag, as documented in the OpenSpec change
// add-live-calendar-subscription.
package calendar

import (
	"context"
	"time"

	"lina-core/pkg/plugin/capability/tenantcap"
)

// Service defines the anonymous public ICS calendar subscription for the
// plugin.
type Service interface {
	// Feed resolves the public live schedule of one room inside the
	// subscription window and projects it into a calendar feed. Tenant
	// resolution mirrors the play endpoint: when the tenant capability is
	// disabled the platform tenant (0) is forced and TenantId is ignored;
	// when enabled, TenantId is required and must reference an existing
	// enabled tenant. Unlike the play endpoint the boundary never fails the
	// caller: a missing tenant parameter, an invalid or disabled tenant, a
	// missing room, and a room without public lives in the window all return
	// an empty feed so calendar clients keep a stable valid document and the
	// room existence stays undisclosed. Database access is one point query
	// for the room plus one window query for the lives.
	Feed(ctx context.Context, in FeedInput) (*CalendarFeed, error)
}

// Ensure serviceImpl implements Service.
var _ Service = (*serviceImpl)(nil)

// serviceImpl implements Service.
type serviceImpl struct {
	tenantSvc tenantcap.Service // Tenant capability bridge for tenant validation
}

// New creates and returns a new Service instance.
func New(tenantSvc tenantcap.Service) Service {
	return &serviceImpl{
		tenantSvc: tenantSvc,
	}
}

// FeedInput defines input for the Feed function.
type FeedInput struct {
	TenantId *int   // Tenant ID scoping the anonymous query; nil means the parameter was absent
	RoomCode string // Live room code, unique within the tenant
	WatchURL string // Absolute viewer H5 page URL embedded into event descriptions; empty omits the link
}

// CalendarFeed defines the renderable projection of one room's public live
// schedule. RoomName stays empty whenever no public live resolved so empty
// feeds rendered for different reasons stay indistinguishable.
type CalendarFeed struct {
	RoomName    string          // Live room name; empty for feeds without entries
	GeneratedAt time.Time       // Feed generation stamp for DTSTAMP
	Events      []CalendarEvent // Public lives inside the subscription window, ordered by live date
}

// CalendarEvent defines one public live rendered as one VEVENT.
type CalendarEvent struct {
	LiveId    int64      // Live content ID backing the stable UID
	Title     string     // Live title rendered as SUMMARY
	LiveDate  time.Time  // Calendar date of the live event
	StartTime *time.Time // Planned or actual start time; nil renders an all-day event
}
