// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
)

// Borrow is the golang structure of table plugin_linapro_book_manage_borrow for DAO operations like Where/Data.
type Borrow struct {
	g.Meta     `orm:"table:plugin_linapro_book_manage_borrow, do:true"`
	Id         any        // Borrow record ID
	TenantId   any        // Owning tenant ID, 0 means PLATFORM
	BookId     any        // Borrowed book ID within the same tenant
	Borrower   any        // Borrower name
	BorrowDate any        // Borrow date with date-only semantics
	DueDate    any        // Expected return date with date-only semantics
	ReturnDate any        // Actual return date with date-only semantics
	Status     any        // Borrow status: 1=borrowed, 2=returned
	Remark     any        // Remark
	CreatedBy  any        // Creator
	UpdatedBy  any        // Updater
	CreatedAt  *time.Time // Creation time
	UpdatedAt  *time.Time // Update time
	DeletedAt  *time.Time // Deletion time
}
