// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// BookDao is the data access object for the table plugin_linapro_book_manage_book.
type BookDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  BookColumns        // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// BookColumns defines and stores column names for the table plugin_linapro_book_manage_book.
type BookColumns struct {
	Id                string // Book ID
	TenantId          string // Owning tenant ID, 0 means PLATFORM
	Title             string // Book title
	Author            string // Author
	Isbn              string // ISBN, unique per tenant when filled
	Category          string // Category: 1=technology, 2=literature, 3=management, 4=children, 5=other
	Publisher         string // Publisher
	PublishDate       string // Publish date with date-only semantics
	TotalQuantity     string // Total copy count
	AvailableQuantity string // Currently available copy count
	Location          string // Shelf location
	Status            string // Book status: 1=on shelf, 2=off shelf
	CoverUrl          string // Cover image URL address
	Remark            string // Remark
	CreatedBy         string // Creator
	UpdatedBy         string // Updater
	CreatedAt         string // Creation time
	UpdatedAt         string // Update time
	DeletedAt         string // Deletion time
}

// bookColumns holds the columns for the table plugin_linapro_book_manage_book.
var bookColumns = BookColumns{
	Id:                "id",
	TenantId:          "tenant_id",
	Title:             "title",
	Author:            "author",
	Isbn:              "isbn",
	Category:          "category",
	Publisher:         "publisher",
	PublishDate:       "publish_date",
	TotalQuantity:     "total_quantity",
	AvailableQuantity: "available_quantity",
	Location:          "location",
	Status:            "status",
	CoverUrl:          "cover_url",
	Remark:            "remark",
	CreatedBy:         "created_by",
	UpdatedBy:         "updated_by",
	CreatedAt:         "created_at",
	UpdatedAt:         "updated_at",
	DeletedAt:         "deleted_at",
}

// NewBookDao creates and returns a new DAO object for table data access.
func NewBookDao(handlers ...gdb.ModelHandler) *BookDao {
	return &BookDao{
		group:    "default",
		table:    "plugin_linapro_book_manage_book",
		columns:  bookColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *BookDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *BookDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *BookDao) Columns() BookColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *BookDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *BookDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *BookDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
