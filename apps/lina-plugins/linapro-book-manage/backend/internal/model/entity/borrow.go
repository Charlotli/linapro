// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"time"
)

// Borrow is the golang structure for table borrow.
type Borrow struct {
	Id         int64      `json:"id"         orm:"id"          description:"Borrow record ID"`
	TenantId   int        `json:"tenantId"   orm:"tenant_id"   description:"Owning tenant ID, 0 means PLATFORM"`
	BookId     int64      `json:"bookId"     orm:"book_id"     description:"Borrowed book ID within the same tenant"`
	Borrower   string     `json:"borrower"   orm:"borrower"    description:"Borrower name"`
	BorrowDate time.Time  `json:"borrowDate" orm:"borrow_date" description:"Borrow date with date-only semantics"`
	DueDate    time.Time  `json:"dueDate"    orm:"due_date"    description:"Expected return date with date-only semantics"`
	ReturnDate time.Time  `json:"returnDate" orm:"return_date" description:"Actual return date with date-only semantics"`
	Status     int        `json:"status"     orm:"status"      description:"Borrow status: 1=borrowed, 2=returned"`
	Remark     string     `json:"remark"     orm:"remark"      description:"Remark"`
	CreatedBy  int64      `json:"createdBy"  orm:"created_by"  description:"Creator"`
	UpdatedBy  int64      `json:"updatedBy"  orm:"updated_by"  description:"Updater"`
	CreatedAt  *time.Time `json:"createdAt"  orm:"created_at"  description:"Creation time"`
	UpdatedAt  *time.Time `json:"updatedAt"  orm:"updated_at"  description:"Update time"`
	DeletedAt  *time.Time `json:"deletedAt"  orm:"deleted_at"  description:"Deletion time"`
}
