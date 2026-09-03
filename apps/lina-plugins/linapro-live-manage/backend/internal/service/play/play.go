// Package play implements the anonymous viewer play-info service for the
// linapro-live-manage source plugin. It resolves one publicly visible live
// record for a live-room code under an explicit tenant context and projects
// the minimal viewer-facing fields. Unlike the tenant-scoped admin services,
// this package deliberately does not read the request business context: the
// tenant is an explicit input, and the authoritative visibility boundary is
// the is_public flag combined with the live state machine, as documented in
// the OpenSpec change add-live-h5-player.
package play

import (
	"context"
	"time"

	"lina-core/pkg/plugin/capability/tenantcap"
)

// Service defines the anonymous public play-info query for the plugin.
type Service interface {
	// Get resolves one playable live for the supplied tenant and live-room
	// code. Tenant resolution: when the tenant capability is disabled the
	// platform tenant (0) is forced and TenantId is ignored; when enabled,
	// TenantId is required (CodePlayTenantRequired when missing) and must
	// reference an existing enabled tenant (CodePlayTenantInvalid otherwise).
	// The candidate priority is ongoing public live first, latest finished
	// public live as replay second, latest not-started public live as preview
	// third; anything else returns CodePlayNotFound. The play URL is only
	// populated for ongoing and finished states, and ongoing lives additionally
	// require the room status to be live. Database access is two point queries
	// on tenant-prefixed indexes.
	Get(ctx context.Context, in GetInput) (*PlayInfo, error)
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

// GetInput defines input for Get function.
type GetInput struct {
	TenantId *int   // Tenant ID scoping the anonymous query; nil means the parameter was absent
	RoomCode string // Live room code, unique within the tenant
}

// PlayInfo defines the viewer-facing play projection of one resolved live.
type PlayInfo struct {
	RoomId           int64      // Owning live room ID
	RoomCode         string     // Live room code from the request
	RoomName         string     // Live room name
	LiveId           int64      // Live content ID; 0 means no live record for the resolved state
	Title            string     // Live title
	CoverUrl         string     // Cover image URL
	LiveUrl          string     // HLS play URL; empty for not-started lives
	LiveDate         time.Time  // Calendar date of the live event
	StartTime        *time.Time // Actual start time; nil when not started
	State            int        // Live state: 0=not started, 1=ongoing, 2=finished
	Host             string     // Host name or role
	SongName         string     // Featured song name
	SongList         string     // Song list JSON array text
	LeadSinger       string     // Lead singer or worship team
	Accompaniment    string     // Accompaniment info
	SermonTitle      string     // Sermon title
	Preacher         string     // Preacher name
	PreacherIdentity string     // Preacher identity or title
	ScriptureRef     string     // Scripture reference
	ScriptureContent string     // Scripture content
	Outline          string     // Sermon outline or agenda
}
