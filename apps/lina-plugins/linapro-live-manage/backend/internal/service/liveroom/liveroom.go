// Package liveroom implements tenant-scoped live-room maintenance services for
// the linapro-live-manage source plugin. It owns the
// plugin_linapro_live_manage_room table access and consumes host capability
// seams for business context, tenant filtering, and user projections. The room
// enumeration values defined here are the plugin-wide source of truth and are
// reused by the live-content service for room availability checks.
package liveroom

import (
	"context"

	"lina-core/pkg/plugin/capability/bizctxcap"
	"lina-core/pkg/plugin/capability/tenantcap"
	"lina-core/pkg/plugin/capability/usercap"
)

// Room type values (matching the plugin_live_room_type dictionary).
const (
	RoomTypeGathering = 1 // Normal gathering room
	RoomTypeEvent     = 2 // Event room
	RoomTypeOther     = 3 // Other room
)

// Room status values (matching the plugin_live_room_status dictionary).
const (
	RoomStatusIdle     = 0 // Idle
	RoomStatusLive     = 1 // Live now
	RoomStatusDisabled = 2 // Disabled
)

// RoomStatus validates one room status value and reports whether it is known
// to the plugin. Unknown statuses are rejected before persistence.
func RoomStatus(value int) bool {
	switch value {
	case RoomStatusIdle, RoomStatusLive, RoomStatusDisabled:
		return true
	default:
		return false
	}
}

// RoomType validates one room type value and reports whether it is known to
// the plugin. Unknown types are rejected before persistence.
func RoomType(value int) bool {
	switch value {
	case RoomTypeGathering, RoomTypeEvent, RoomTypeOther:
		return true
	default:
		return false
	}
}

// roomStatusValues lists every valid room status for validation messages.
var roomStatusValues = []int{RoomStatusIdle, RoomStatusLive, RoomStatusDisabled}

// roomTypeValues lists every valid room type for validation messages.
var roomTypeValues = []int{RoomTypeGathering, RoomTypeEvent, RoomTypeOther}

// OptionsLimit is the hard maximum cap for live-room candidate lists.
const OptionsLimit = 200

// OptionsDefaultLimit is the default cap applied when a candidate query
// requests no explicit limit.
const OptionsDefaultLimit = 100

// Service defines tenant-scoped live-room CRUD and bounded option candidates
// for the plugin.
type Service interface {
	// List returns live rooms visible to ctx's tenant using the supplied page
	// and fuzzy room-name/type/status filters. Status filtering is skipped when
	// Status is nil. It returns DAO or tenant-scope errors.
	List(ctx context.Context, in ListInput) (*ListOutput, error)
	// Options returns the bounded minimal room projection for select controls,
	// filtered by keyword. Disabled rooms are excluded unless IncludeDisabled is
	// true. The result is capped at 100 items by default and never exceeds
	// OptionsLimit items regardless of the requested Limit.
	Options(ctx context.Context, in OptionsInput) ([]*OptionItem, error)
	// GetById returns one tenant-visible live room by primary key with creator
	// display metadata, or CodeRoomNotFound when the row is outside scope or
	// absent.
	GetById(ctx context.Context, id int64) (*ListItem, error)
	// Create creates a new live room owned by the current tenant and user from
	// ctx. The room code must be unique within the tenant; a conflict returns
	// CodeRoomCodeExists.
	Create(ctx context.Context, in CreateInput) (int64, error)
	// Update updates the live-room fields. It only updates the current tenant's
	// row and returns CodeRoomNotFound for missing scope.
	Update(ctx context.Context, in UpdateInput) error
	// Delete soft-deletes live rooms by IDs. Rooms still referenced by live
	// content are rejected as a whole with CodeRoomReferenced; empty IDs return
	// a business error.
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
	PageNum  int    // Page number, starting from 1
	PageSize int    // Page size
	RoomName string // Room name, supports fuzzy search
	RoomType int    // Room type: 1=gathering 2=event 3=other; 0 means no filter
	Status   *int   // Room status: 0=idle 1=live 2=disabled; nil means no filter
}

// OptionsInput defines input for Options function.
type OptionsInput struct {
	Keyword         string // Keyword, matches room name or room code fuzzily
	IncludeDisabled bool   // Whether to include disabled rooms
	Limit           int    // Maximum options to return; <=0 means OptionsDefaultLimit, capped at OptionsLimit
}

// OptionItem defines one minimal live-room option projection.
type OptionItem struct {
	Id       int64  // Live room ID
	RoomCode string // Room code
	RoomName string // Room name
	Status   int    // Room status: 0=idle 1=live 2=disabled
}

// ListItem defines a single list item.
type ListItem struct {
	*RoomEntity          // Live room entity
	CreatedByName string // Creator username
}

// ListOutput defines output for List function.
type ListOutput struct {
	List  []*ListItem // List items
	Total int         // Total count
}

// CreateInput defines input for Create function.
type CreateInput struct {
	RoomCode    string // Room code, unique per tenant
	RoomName    string // Room name
	RoomType    int    // Room type: 1=gathering 2=event 3=other
	Status      int    // Room status: 0=idle 1=live 2=disabled
	Description string // Room description
}

// UpdateInput defines input for Update function.
type UpdateInput struct {
	Id          int64   // Live room ID
	RoomCode    *string // Room code, unique per tenant
	RoomName    *string // Room name
	RoomType    *int    // Room type: 1=gathering 2=event 3=other
	Status      *int    // Room status: 0=idle 1=live 2=disabled
	Description *string // Room description
}
