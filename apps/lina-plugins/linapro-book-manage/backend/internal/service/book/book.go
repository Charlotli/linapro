// Package book implements tenant-scoped book registry services for the
// linapro-book-manage source plugin. It owns the
// plugin_linapro_book_manage_book table access and consumes host capability
// seams for business context, tenant filtering, and user projections.
package book

import (
	"context"

	"lina-core/pkg/plugin/capability/bizctxcap"
	"lina-core/pkg/plugin/capability/tenantcap"
	"lina-core/pkg/plugin/capability/usercap"
)

// Book category values (matching the plugin_book_category dictionary).
const (
	CategoryTechnology = 1 // Technology
	CategoryLiterature = 2 // Literature
	CategoryManagement = 3 // Management
	CategoryChildren   = 4 // Children
	CategoryOther      = 5 // Other
)

// Book status values (matching the book status column semantics).
const (
	BookStatusOnShelf  = 1 // On shelf
	BookStatusOffShelf = 2 // Off shelf
)

// Category validates one book category value.
func Category(value int) bool {
	switch value {
	case CategoryTechnology, CategoryLiterature, CategoryManagement, CategoryChildren, CategoryOther:
		return true
	default:
		return false
	}
}

// BookStatus validates one book status value.
func BookStatus(value int) bool {
	switch value {
	case BookStatusOnShelf, BookStatusOffShelf:
		return true
	default:
		return false
	}
}

// Service defines tenant-scoped book registry services for the plugin.
type Service interface {
	// List returns books visible to ctx's tenant using the supplied page and
	// title/author/category/status filters. Status filtering is skipped when
	// Status is nil. It returns DAO or tenant-scope errors.
	List(ctx context.Context, in ListInput) (*ListOutput, error)
	// GetById returns one tenant-visible book by primary key, or CodeBookNotFound
	// when the row is outside scope or absent.
	GetById(ctx context.Context, id int64) (*BookEntity, error)
	// Create creates one book owned by the current tenant and user from ctx.
	// Available quantity defaults to the total quantity. ISBN must be unique
	// within the tenant when filled.
	Create(ctx context.Context, in CreateInput) (int64, error)
	// Update updates the book fields. Total quantity adjustments cannot go
	// below the number of copies currently lent out.
	Update(ctx context.Context, in UpdateInput) error
	// Delete soft-deletes books by IDs. Books still referenced by borrow
	// records are rejected as a whole; empty IDs return a business error.
	Delete(ctx context.Context, ids []int64) error
	// Options returns bounded minimal book candidates for select controls,
	// filtered by keyword across title and author. Only on-shelf books are
	// returned. The result never exceeds OptionsLimit items.
	Options(ctx context.Context, in OptionsInput) ([]*OptionItem, error)
}

// OptionsLimit bounds the book candidate list.
const OptionsLimit = 200

// OptionsInput defines input for Options function.
type OptionsInput struct {
	Keyword string // Keyword, matches title or author fuzzily
}

// OptionItem defines one minimal book option projection.
type OptionItem struct {
	Id                int64  // Book ID
	Title             string // Book title
	Author            string // Author
	AvailableQuantity int    // Currently available copy count
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
	Title    string // Title or author, supports fuzzy search
	Category int    // Category; 0 means no filter
	Status   *int   // Book status; nil means no filter
}

// ListOutput defines output for List function.
type ListOutput struct {
	List  []*BookEntity // List items
	Total int           // Total count
}

// CreateInput defines input for Create function.
type CreateInput struct {
	Title         string // Book title
	Author        string // Author
	Isbn          string // ISBN, unique per tenant when filled
	Category      int    // Category
	Publisher     string // Publisher
	PublishDate   string // Publish date, YYYY-MM-DD; empty means unset
	TotalQuantity int    // Total copy count
	Location      string // Shelf location
	Status        int    // Book status
	CoverUrl      string // Cover image URL
	Remark        string // Remark
}

// UpdateInput defines input for Update function.
type UpdateInput struct {
	Id            int64   // Book ID
	Title         *string // Book title; nil keeps current value
	Author        *string // Author; nil keeps current value
	Isbn          *string // ISBN; nil keeps current value
	Category      *int    // Category; nil keeps current value
	Publisher     *string // Publisher; nil keeps current value
	PublishDate   *string // Publish date, YYYY-MM-DD; nil keeps current value
	TotalQuantity *int    // Total copy count; nil keeps current value
	Location      *string // Shelf location; nil keeps current value
	Status        *int    // Book status; nil keeps current value
	CoverUrl      *string // Cover image URL; nil keeps current value
	Remark        *string // Remark; nil keeps current value
}
