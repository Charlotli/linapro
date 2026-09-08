// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// AnnouncementDao is the data access object for the table plugin_linapro_live_manage_announcement.
type AnnouncementDao struct {
	table    string              // table is the underlying table name of the DAO.
	group    string              // group is the database configuration group name of the current DAO.
	columns  AnnouncementColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler  // handlers for customized model modification.
}

// AnnouncementColumns defines and stores column names for the table plugin_linapro_live_manage_announcement.
type AnnouncementColumns struct {
	Id        string // Announcement ID
	TenantId  string // Owning tenant ID, 0 means PLATFORM
	RoomId    string // Owning live room ID within the same tenant
	Title     string // Announcement title
	Content   string // Announcement body, plain text rendered with line breaks
	Enabled   string // Whether viewers can see the announcement
	Sort      string // Display order, smaller values first
	CreatedBy string // Creator
	UpdatedBy string // Updater
	CreatedAt string // Creation time
	UpdatedAt string // Update time
	DeletedAt string // Deletion time
}

// announcementColumns holds the columns for the table plugin_linapro_live_manage_announcement.
var announcementColumns = AnnouncementColumns{
	Id:        "id",
	TenantId:  "tenant_id",
	RoomId:    "room_id",
	Title:     "title",
	Content:   "content",
	Enabled:   "enabled",
	Sort:      "sort",
	CreatedBy: "created_by",
	UpdatedBy: "updated_by",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
	DeletedAt: "deleted_at",
}

// NewAnnouncementDao creates and returns a new DAO object for table data access.
func NewAnnouncementDao(handlers ...gdb.ModelHandler) *AnnouncementDao {
	return &AnnouncementDao{
		group:    "default",
		table:    "plugin_linapro_live_manage_announcement",
		columns:  announcementColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *AnnouncementDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *AnnouncementDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *AnnouncementDao) Columns() AnnouncementColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *AnnouncementDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *AnnouncementDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *AnnouncementDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
