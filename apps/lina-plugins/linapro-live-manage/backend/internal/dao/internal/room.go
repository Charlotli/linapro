// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// RoomDao is the data access object for the table plugin_linapro_live_manage_room.
type RoomDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  RoomColumns        // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// RoomColumns defines and stores column names for the table plugin_linapro_live_manage_room.
type RoomColumns struct {
	Id          string // Live room ID
	TenantId    string // Owning tenant ID, 0 means PLATFORM
	RoomCode    string // Room code exposed to external systems, unique per tenant
	RoomName    string // Room name
	RoomType    string // Room type: 1=gathering, 2=event, 3=other
	Status      string // Room status: 0=idle, 1=live, 2=disabled
	Description string // Room description
	CreatedBy   string // Creator
	UpdatedBy   string // Updater
	CreatedAt   string // Creation time
	UpdatedAt   string // Update time
	DeletedAt   string // Deletion time
}

// roomColumns holds the columns for the table plugin_linapro_live_manage_room.
var roomColumns = RoomColumns{
	Id:          "id",
	TenantId:    "tenant_id",
	RoomCode:    "room_code",
	RoomName:    "room_name",
	RoomType:    "room_type",
	Status:      "status",
	Description: "description",
	CreatedBy:   "created_by",
	UpdatedBy:   "updated_by",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
	DeletedAt:   "deleted_at",
}

// NewRoomDao creates and returns a new DAO object for table data access.
func NewRoomDao(handlers ...gdb.ModelHandler) *RoomDao {
	return &RoomDao{
		group:    "default",
		table:    "plugin_linapro_live_manage_room",
		columns:  roomColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *RoomDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *RoomDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *RoomDao) Columns() RoomColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *RoomDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *RoomDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *RoomDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
