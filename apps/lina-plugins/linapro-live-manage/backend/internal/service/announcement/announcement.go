// Package announcement implements the per-room viewer announcement service
// for the linapro-live-manage source plugin. The admin side runs behind the
// host auth, tenancy, and permission chain and scopes every row to the
// request tenant; the viewer side answers anonymously through the shared
// play service room resolution so the tenant and room visibility contracts
// stay aligned with the public play endpoint, as documented in the OpenSpec
// change add-live-announcement-and-bible.
package announcement

import (
	"context"

	"lina-core/pkg/plugin/capability/bizctxcap"
	"lina-core/pkg/plugin/capability/tenantcap"
	"lina-core/pkg/plugin/capability/usercap"
	"lina-plugin-linapro-live-manage/backend/internal/service/play"
)

// ViewerLimit is the hard maximum cap for the anonymous viewer announcement
// list; announcements are short-lived room notices, so the bound protects the
// H5 payload without pagination.
const ViewerLimit = 50

// PageDefaultSize is the admin list page size applied when the input leaves it unset.
const PageDefaultSize = 10

// PageMaxSize is the hard maximum cap for the admin list page size.
const PageMaxSize = 100

// Service defines the announcement management and public listing contract.
type Service interface {
	// List returns the paged announcement list of one room for the admin
	// console, disabled rows included, ordered by sort then ID. The room must
	// belong to the current tenant; CodeAnnouncementRoomNotFound is returned
	// otherwise. Page defaults to 1 and PageSize defaults to 10 with a hard
	// cap of 100.
	List(ctx context.Context, in ListInput) ([]*Item, int, error)
	// Create inserts one announcement for a room owned by the current tenant
	// and returns the new ID. Enabled defaults to true when the input leaves
	// it nil.
	Create(ctx context.Context, in CreateInput) (int64, error)
	// Update mutates the mutable fields of one announcement owned by the
	// current tenant; nil fields keep their stored values. Returns
	// CodeAnnouncementNotFound when the row is missing or soft-deleted.
	Update(ctx context.Context, in UpdateInput) error
	// Delete soft-deletes one announcement owned by the current tenant.
	// Deleting an already deleted row succeeds silently to keep the admin
	// flow idempotent.
	Delete(ctx context.Context, id int64) error
	// ForRoom returns the enabled announcements of one room for the anonymous
	// viewer endpoint, ordered by sort then creation, capped at ViewerLimit.
	// Tenant resolution matches the public play endpoint, including verbatim
	// tenant errors. Unknown rooms and tenant failures answer with an empty
	// list so room existence is never leaked.
	ForRoom(ctx context.Context, in ForRoomInput) ([]*ViewerItem, error)
}

// Ensure serviceImpl implements Service.
var _ Service = (*serviceImpl)(nil)

// serviceImpl implements Service.
type serviceImpl struct {
	bizCtxSvc   bizctxcap.Service // Business context bridge for admin audit fields
	tenantSvc   tenantcap.Service // Tenant capability bridge for admin scoping
	userSvc     usercap.Service   // User domain projection capability
	playInfoSvc play.Service      // Shared play service for anonymous room resolution
}

// New creates and returns a new Service instance wired with its startup
// dependencies. All collaborators are host-owned capability services or the
// shared play service, resolved once during plugin route registration.
func New(
	bizCtxSvc bizctxcap.Service,
	tenantSvc tenantcap.Service,
	userSvc usercap.Service,
	playInfoSvc play.Service,
) Service {
	return &serviceImpl{
		bizCtxSvc:   bizCtxSvc,
		tenantSvc:   tenantSvc,
		userSvc:     userSvc,
		playInfoSvc: playInfoSvc,
	}
}

// ListInput defines input for List function.
type ListInput struct {
	RoomId   int64 // Owning live room ID within the current tenant
	PageNum  int   // Page number starting from 1; zero or negative defaults to 1
	PageSize int   // Page size; zero defaults to 10, values above the cap clamp to 100
}

// CreateInput defines input for Create function.
type CreateInput struct {
	RoomId  int64  // Owning live room ID within the current tenant
	Title   string // Announcement title
	Content string // Announcement body, plain text
	Enabled *bool  // Viewer visibility; nil defaults to true
	Sort    int    // Display order, smaller values first
}

// UpdateInput defines input for Update function. Nil fields keep stored values.
type UpdateInput struct {
	Id      int64   // Announcement ID
	Title   *string // Title; nil keeps the stored title
	Content *string // Body; nil keeps the stored body
	Enabled *bool   // Viewer visibility; nil keeps the stored state
	Sort    *int    // Display order; nil keeps the stored order
}

// ForRoomInput defines input for ForRoom function.
type ForRoomInput struct {
	TenantId *int   // Tenant ID scoping the anonymous query; nil means the parameter was absent
	RoomCode string // Live room code, unique within the tenant
}

// Item defines the admin-side announcement projection.
type Item struct {
	Id        int64  // Announcement ID
	RoomId    int64  // Owning live room ID
	Title     string // Announcement title
	Content   string // Announcement body
	Enabled   bool   // Viewer visibility
	Sort      int    // Display order
	CreatedAt int64  // Creation time as Unix milliseconds; 0 when unset
	UpdatedAt int64  // Last update time as Unix milliseconds; 0 when unset
}

// ViewerItem defines the anonymous viewer announcement projection.
type ViewerItem struct {
	Id        int64  // Announcement ID
	Title     string // Announcement title
	Content   string // Announcement body
	UpdatedAt int64  // Last update time as Unix milliseconds; 0 when unset
}
