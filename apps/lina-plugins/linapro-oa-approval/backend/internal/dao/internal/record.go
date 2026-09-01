// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// RecordDao is the data access object for the table plugin_linapro_oa_approval_record.
type RecordDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  RecordColumns      // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// RecordColumns defines and stores column names for the table plugin_linapro_oa_approval_record.
type RecordColumns struct {
	Id        string // Approval record ID
	TenantId  string // Owning tenant ID, 0 means PLATFORM
	RequestId string // Owning approval request ID
	NodeOrder string // Related node order; 0 means request-level actions
	Action    string // Action: 1=submit, 2=approve, 3=reject, 4=comment, 5=append approver, 6=withdraw
	ActorId   string // Actor user ID
	Comment   string // Reply or comment content attached to the action
	CreatedBy string // Creator
	UpdatedBy string // Updater
	CreatedAt string // Creation time
	UpdatedAt string // Update time
	DeletedAt string // Deletion time
}

// recordColumns holds the columns for the table plugin_linapro_oa_approval_record.
var recordColumns = RecordColumns{
	Id:        "id",
	TenantId:  "tenant_id",
	RequestId: "request_id",
	NodeOrder: "node_order",
	Action:    "action",
	ActorId:   "actor_id",
	Comment:   "comment",
	CreatedBy: "created_by",
	UpdatedBy: "updated_by",
	CreatedAt: "created_at",
	UpdatedAt: "updated_at",
	DeletedAt: "deleted_at",
}

// NewRecordDao creates and returns a new DAO object for table data access.
func NewRecordDao(handlers ...gdb.ModelHandler) *RecordDao {
	return &RecordDao{
		group:    "default",
		table:    "plugin_linapro_oa_approval_record",
		columns:  recordColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *RecordDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *RecordDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *RecordDao) Columns() RecordColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *RecordDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *RecordDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *RecordDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
