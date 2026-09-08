// Package view implements the anonymous viewer watch-statistics service for
// the linapro-live-manage source plugin. It records anonymous watch sessions
// idempotently keyed by a browser-generated session key and aggregates the
// online count (sessions with a heartbeat inside the freshness window) and
// the total deduplicated view count per live. Like the play service, this
// package deliberately does not read the request business context: the tenant
// is an explicit input validated against the tenant lifecycle, as documented
// in the OpenSpec change add-live-share-and-view-stats.
package view

import (
	"context"
	"time"

	"lina-core/pkg/plugin/capability/tenantcap"
	"lina-plugin-linapro-live-manage/backend/internal/service/play"
)

// Statistics tuning constants.
const (
	// onlineWindow is the heartbeat freshness window: a session counts as
	// online when its latest heartbeat is newer than now minus this duration.
	onlineWindow = 60 * time.Second
)

// Service defines the anonymous public watch-statistics service for the
// plugin.
type Service interface {
	// Heartbeat records one viewer watch session for the resolved playable
	// live and returns the aggregate counters. Tenant resolution is delegated
	// to the play service, so the tenant semantics can never drift from the
	// play contract: the platform tenant (0) is forced when the tenant
	// capability is disabled, and when enabled a missing or invalid TenantId
	// returns the play tenant errors verbatim. Sessions are upserted on the
	// unique session key, so repeated heartbeats never inflate the total
	// count. When the room has no publicly playable live, the service returns
	// zero counters without writing any row so invalid room codes and absent
	// content are indistinguishable. Database access is at most two
	// statements per call regardless of viewer scale.
	Heartbeat(ctx context.Context, in HeartbeatInput) (*Stats, error)
	// BatchGet aggregates the online and total view counters for the supplied
	// live ID set in one grouped query and maps the results by live ID. Live
	// IDs without any session are absent from the returned map; missing
	// counters mean zero. The query count is exactly one regardless of the
	// input size.
	BatchGet(ctx context.Context, liveIDs []int64) (map[int64]*Stats, error)
}

// Ensure serviceImpl implements Service.
var _ Service = (*serviceImpl)(nil)

// serviceImpl implements Service.
type serviceImpl struct {
	tenantSvc tenantcap.Service // Tenant capability bridge for tenant validation
	playSvc   play.Service      // Shared playable-live resolution service
}

// New creates and returns a new Service instance.
func New(tenantSvc tenantcap.Service, playSvc play.Service) Service {
	return &serviceImpl{
		tenantSvc: tenantSvc,
		playSvc:   playSvc,
	}
}

// HeartbeatInput defines input for Heartbeat function.
type HeartbeatInput struct {
	TenantId   *int   // Tenant ID scoping the anonymous query; nil means the parameter was absent
	RoomCode   string // Live room code, unique within the tenant
	SessionKey string // Browser-generated anonymous viewer session key
}

// Stats defines the per-live watch counters.
type Stats struct {
	OnlineCount int64 // Sessions with a heartbeat inside the online window
	TotalViews  int64 // Deduplicated watch sessions recorded for the live
}
