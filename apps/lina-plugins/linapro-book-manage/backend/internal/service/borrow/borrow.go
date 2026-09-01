// Package borrow implements tenant-scoped book borrow services for the
// linapro-book-manage source plugin. Borrow and return actions link the
// borrow record with the book's available copy count inside one transaction;
// concurrent borrows are guarded by a conditional update on the remaining
// available count.
package borrow

import (
	"context"

	"lina-core/pkg/plugin/capability/bizctxcap"
	"lina-core/pkg/plugin/capability/tenantcap"
)

// Borrow status values (matching the plugin_book_borrow_status dictionary).
const (
	BorrowStatusBorrowed = 1 // Borrowed
	BorrowStatusReturned = 2 // Returned
)

// Service defines tenant-scoped book borrow services for the plugin.
type Service interface {
	// List returns borrow records visible to ctx's tenant using the supplied
	// page and book/borrower/status filters. Book titles are resolved with one
	// batch query for the current page.
	List(ctx context.Context, in ListInput) (*ListOutput, error)
	// GetById returns one tenant-visible borrow record by primary key, or
	// CodeBorrowNotFound when the row is outside scope or absent.
	GetById(ctx context.Context, id int64) (*BorrowEntity, error)
	// Create borrows one book: the book must be on shelf with available
	// copies; the record is created and the available count decremented in one
	// transaction.
	Create(ctx context.Context, in CreateInput) (int64, error)
	// Return marks one borrowed record as returned. Only the record creator
	// can return it; the return date is recorded and the available count is
	// incremented in one transaction.
	Return(ctx context.Context, in ReturnInput) error
	// Update updates editable borrow fields (borrower, due date, remark).
	Update(ctx context.Context, in UpdateInput) error
	// Delete soft-deletes returned borrow records by IDs. Records still in
	// borrowed status are rejected; empty IDs return a business error.
	Delete(ctx context.Context, ids []int64) error
}

// Ensure serviceImpl implements Service.
var _ Service = (*serviceImpl)(nil)

// serviceImpl implements Service.
type serviceImpl struct {
	bizCtxSvc bizctxcap.Service // Business context bridge
	tenantSvc tenantcap.Service // Tenant capability bridge
}

// New creates and returns a new Service instance.
func New(
	bizCtxSvc bizctxcap.Service,
	tenantSvc tenantcap.Service,
) Service {
	return &serviceImpl{
		bizCtxSvc: bizCtxSvc,
		tenantSvc: tenantSvc,
	}
}

// ListInput defines input for List function.
type ListInput struct {
	PageNum  int    // Page number, starting from 1
	PageSize int    // Page size
	BookId   int64  // Book ID; 0 means no filter
	Borrower string // Borrower name, supports fuzzy search
	Status   *int   // Borrow status; nil means no filter
}

// ListOutput defines output for List function.
type ListOutput struct {
	List  []*BorrowItem // List items
	Total int           // Total count
}

// BorrowItem defines one list projection item.
type BorrowItem struct {
	*BorrowEntity
	BookTitle string // Book title resolved in batch
}

// CreateInput defines input for Create function.
type CreateInput struct {
	BookId     int64  // Borrowed book ID
	Borrower   string // Borrower name
	BorrowDate string // Borrow date, YYYY-MM-DD; empty defaults to today
	DueDate    string // Expected return date, YYYY-MM-DD; empty means unset
	Remark     string // Remark
}

// ReturnInput defines input for the return action.
type ReturnInput struct {
	Id         int64  // Borrow record ID
	ReturnDate string // Actual return date, YYYY-MM-DD; empty defaults to today
}

// UpdateInput defines input for Update function.
type UpdateInput struct {
	Id       int64   // Borrow record ID
	Borrower *string // Borrower name; nil keeps current value
	DueDate  *string // Expected return date, YYYY-MM-DD; nil keeps current value
	Remark   *string // Remark; nil keeps current value
}
