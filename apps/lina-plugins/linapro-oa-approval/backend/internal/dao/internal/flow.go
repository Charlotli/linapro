// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// FlowDao is the data access object for the table plugin_linapro_oa_approval_flow.
type FlowDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  FlowColumns        // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// FlowColumns defines and stores column names for the table plugin_linapro_oa_approval_flow.
type FlowColumns struct {
	Id          string // Approval flow ID
	TenantId    string // Owning tenant ID, 0 means PLATFORM
	FlowType    string // Flow type: 1=fund request, 2=write-off, 3=expense reimbursement
	FlowName    string // Flow name, unique per tenant
	Description string // Flow description
	Status      string // Flow status: 1=enabled, 0=disabled
	CreatedBy   string // Creator
	UpdatedBy   string // Updater
	CreatedAt   string // Creation time
	UpdatedAt   string // Update time
	DeletedAt   string // Deletion time
	FormFields  string // Configurable form field definitions as a JSON array
}

// flowColumns holds the columns for the table plugin_linapro_oa_approval_flow.
var flowColumns = FlowColumns{
	Id:          "id",
	TenantId:    "tenant_id",
	FlowType:    "flow_type",
	FlowName:    "flow_name",
	Description: "description",
	Status:      "status",
	CreatedBy:   "created_by",
	UpdatedBy:   "updated_by",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
	DeletedAt:   "deleted_at",
	FormFields:  "form_fields",
}

// NewFlowDao creates and returns a new DAO object for table data access.
func NewFlowDao(handlers ...gdb.ModelHandler) *FlowDao {
	return &FlowDao{
		group:    "default",
		table:    "plugin_linapro_oa_approval_flow",
		columns:  flowColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *FlowDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *FlowDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *FlowDao) Columns() FlowColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *FlowDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *FlowDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *FlowDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
