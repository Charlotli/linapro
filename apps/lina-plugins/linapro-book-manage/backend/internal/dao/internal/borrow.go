// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// BorrowDao is the data access object for the table plugin_linapro_book_manage_borrow.
type BorrowDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  BorrowColumns      // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// BorrowColumns defines and stores column names for the table plugin_linapro_book_manage_borrow.
type BorrowColumns struct {
	Id         string // Borrow record ID
	TenantId   string // Owning tenant ID, 0 means PLATFORM
	BookId     string // Borrowed book ID within the same tenant
	Borrower   string // Borrower name
	BorrowDate string // Borrow date with date-only semantics
	DueDate    string // Expected return date with date-only semantics
	ReturnDate string // Actual return date with date-only semantics
	Status     string // Borrow status: 1=borrowed, 2=returned
	Remark     string // Remark
	CreatedBy  string // Creator
	UpdatedBy  string // Updater
	CreatedAt  string // Creation time
	UpdatedAt  string // Update time
	DeletedAt  string // Deletion time
}

// borrowColumns holds the columns for the table plugin_linapro_book_manage_borrow.
var borrowColumns = BorrowColumns{
	Id:         "id",
	TenantId:   "tenant_id",
	BookId:     "book_id",
	Borrower:   "borrower",
	BorrowDate: "borrow_date",
	DueDate:    "due_date",
	ReturnDate: "return_date",
	Status:     "status",
	Remark:     "remark",
	CreatedBy:  "created_by",
	UpdatedBy:  "updated_by",
	CreatedAt:  "created_at",
	UpdatedAt:  "updated_at",
	DeletedAt:  "deleted_at",
}

// NewBorrowDao creates and returns a new DAO object for table data access.
func NewBorrowDao(handlers ...gdb.ModelHandler) *BorrowDao {
	return &BorrowDao{
		group:    "default",
		table:    "plugin_linapro_book_manage_borrow",
		columns:  borrowColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *BorrowDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *BorrowDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *BorrowDao) Columns() BorrowColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *BorrowDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *BorrowDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rolls back the transaction and returns the error if function f returns a non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note: Do not commit or roll back the transaction in function f,
// as it is automatically handled by this function.
func (dao *BorrowDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
