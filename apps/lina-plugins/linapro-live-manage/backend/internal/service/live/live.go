// Package live implements tenant-scoped live-content maintenance services for
// the linapro-live-manage source plugin. It owns the
// plugin_linapro_live_manage_live table access, validates live-room
// availability for the live-room enumeration values shared from the liveroom
// package, and consumes host capability seams for business context, tenant
// filtering, and user projections.
package live

import (
	"context"

	"lina-core/pkg/plugin/capability/bizctxcap"
	"lina-core/pkg/plugin/capability/tenantcap"
	"lina-core/pkg/plugin/capability/usercap"
)

// Live state values (matching the plugin_live_state dictionary).
const (
	liveStateNotStarted = 0 // Not started
	liveStateOngoing    = 1 // Ongoing
	liveStateFinished   = 2 // Finished
)

// Live visibility values (matching the plugin_live_public dictionary).
const (
	liveVisibilityPrivate = 0 // Private
	liveVisibilityPublic  = 1 // Public
)

// Live state validation helper reused by create and update paths.
func liveStateValid(value int) bool {
	switch value {
	case liveStateNotStarted, liveStateOngoing, liveStateFinished:
		return true
	default:
		return false
	}
}

// liveStateTransition validates one explicit state action against the live
// state machine: only not-started lives can start and only ongoing lives can
// stop. It returns the target state or false when the transition is illegal.
func liveStateTransition(current int, action liveStateAction) (int, bool) {
	switch action {
	case liveStateActionStart:
		if current == liveStateNotStarted {
			return liveStateOngoing, true
		}
	case liveStateActionStop:
		if current == liveStateOngoing {
			return liveStateFinished, true
		}
	}
	return current, false
}

// liveStateAction identifies one explicit live state action.
type liveStateAction int

const (
	// liveStateActionStart transitions one not-started live to ongoing.
	liveStateActionStart liveStateAction = iota
	// liveStateActionStop transitions one ongoing live to finished.
	liveStateActionStop
)

// Live visibility validation helper reused by create and update paths.
func liveVisibilityValid(value int) bool {
	switch value {
	case liveVisibilityPrivate, liveVisibilityPublic:
		return true
	default:
		return false
	}
}

// Service defines tenant-scoped live-content CRUD for the plugin.
type Service interface {
	// List returns live content visible to ctx's tenant using the supplied page
	// and title/room/state/visibility/date-range filters. Live-room names are
	// resolved with one batch query for the current page. State and visibility
	// filtering is skipped when the corresponding pointer is nil. It returns DAO
	// or tenant-scope errors.
	List(ctx context.Context, in ListInput) (*ListOutput, error)
	// GetById returns one tenant-visible live content record by primary key with
	// the live-room name and creator display metadata, or CodeLiveNotFound when
	// the row is outside scope or absent.
	GetById(ctx context.Context, id int64) (*ListItem, error)
	// Create creates one live content record owned by the current tenant and
	// user from ctx. The live room must exist in the same tenant and must not be
	// disabled; violations return CodeRoomUnavailable. A non-empty SongList must
	// be a valid JSON array; violations return CodeSongListInvalid.
	Create(ctx context.Context, in CreateInput) (int64, error)
	// Update updates the live-content fields. Room availability is re-checked
	// when RoomId is updated. It only updates the current tenant's row and
	// returns CodeLiveNotFound for missing scope.
	Update(ctx context.Context, in UpdateInput) error
	// Start transitions one tenant-visible not-started live to ongoing. It
	// records the start time, flips the live room to live status, and rejects
	// rooms that are disabled or already hosting another ongoing live. State
	// machine violations return CodeLiveStateTransition; busy rooms return
	// CodeRoomBusy. The live and room updates run in one transaction.
	Start(ctx context.Context, id int64) error
	// Stop transitions one tenant-visible ongoing live to finished. It restores
	// the live room to idle when the room is still live, keeping the recorded
	// start time. State machine violations return CodeLiveStateTransition. The
	// live and room updates run in one transaction.
	Stop(ctx context.Context, id int64) error
	// Delete soft-deletes live content records by IDs. Empty IDs return a
	// business error.
	Delete(ctx context.Context, ids []int64) error
}

// Ensure serviceImpl implements Service.
var _ Service = (*serviceImpl)(nil)

// serviceImpl implements Service.
type serviceImpl struct {
	bizCtxSvc bizctxcap.Service // Business context bridge
	tenantSvc tenantcap.Service // Tenant capability bridge
	userSvc   usercap.Service   // User domain projection capability
}

// New creates and returns a new Service instance.
func New(
	bizCtxSvc bizctxcap.Service,
	tenantSvc tenantcap.Service,
	userSvc usercap.Service,
) Service {
	return &serviceImpl{
		bizCtxSvc: bizCtxSvc,
		tenantSvc: tenantSvc,
		userSvc:   userSvc,
	}
}

// ListInput defines input for List function.
type ListInput struct {
	PageNum       int    // Page number, starting from 1
	PageSize      int    // Page size
	Title         string // Title, supports fuzzy search
	RoomId        int64  // Live room ID; 0 means no filter
	State         *int   // Live state: 0=not started 1=ongoing 2=finished; nil means no filter
	IsPublic      *int   // Visibility: 1=public 0=private; nil means no filter
	LiveDateStart string // Live date range start, YYYY-MM-DD, inclusive; empty means unbounded
	LiveDateEnd   string // Live date range end, YYYY-MM-DD, inclusive; empty means unbounded
}

// ListItem defines a single list item.
type ListItem struct {
	*LiveEntity          // Live content entity
	RoomName      string // Live room name resolved in batch
	CreatedByName string // Creator username
}

// ListOutput defines output for List function.
type ListOutput struct {
	List  []*ListItem // List items
	Total int         // Total count
}

// CreateInput defines input for Create function.
type CreateInput struct {
	RoomId           int64  // Owning live room ID
	Title            string // Live title
	Content          string // Live description or rich content
	LiveDate         string // Live date, YYYY-MM-DD; empty defaults to today
	PageUrl          string // External page link
	PushUrl          string // Stream push URL
	LiveUrl          string // Play URL for viewers
	CoverUrl         string // Cover image URL
	SongName         string // Featured song name
	SongList         string // Song list JSON array text
	LeadSinger       string // Lead singer or worship team
	Accompaniment    string // Accompaniment info
	Host             string // Host name or role
	SermonTitle      string // Sermon title
	Preacher         string // Preacher name
	PreacherIdentity string // Preacher identity
	ScriptureRef     string // Scripture reference
	ScriptureContent string // Scripture content
	Outline          string // Sermon outline
	DeviceInfo       string // Device info free text
	Reception        string // Reception arrangement
	State            int    // Live state: 0=not started 1=ongoing 2=finished
	IsPublic         int    // Visibility: 1=public 0=private
	StartTime        *int64 // Start time as Unix milliseconds; nil means unset
}

// UpdateInput defines input for Update function.
type UpdateInput struct {
	Id               int64   // Live content ID
	RoomId           *int64  // Owning live room ID
	Title            *string // Live title
	Content          *string // Live description or rich content
	LiveDate         *string // Live date, YYYY-MM-DD
	PageUrl          *string // External page link
	PushUrl          *string // Stream push URL
	LiveUrl          *string // Play URL for viewers
	CoverUrl         *string // Cover image URL
	SongName         *string // Featured song name
	SongList         *string // Song list JSON array text
	LeadSinger       *string // Lead singer or worship team
	Accompaniment    *string // Accompaniment info
	Host             *string // Host name or role
	SermonTitle      *string // Sermon title
	Preacher         *string // Preacher name
	PreacherIdentity *string // Preacher identity
	ScriptureRef     *string // Scripture reference
	ScriptureContent *string // Scripture content
	Outline          *string // Sermon outline
	DeviceInfo       *string // Device info free text
	Reception        *string // Reception arrangement
	State            *int    // Live state: 0=not started 1=ongoing 2=finished
	IsPublic         *int    // Visibility: 1=public 0=private
	StartTime        *int64  // Start time as Unix milliseconds
}
