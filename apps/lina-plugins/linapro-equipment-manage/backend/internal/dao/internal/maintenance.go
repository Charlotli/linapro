// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// MaintenanceDao is the data access object for the table plugin_linapro_equipment_manage_maintenance.
type MaintenanceDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  MaintenanceColumns // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// MaintenanceColumns defines and stores column names for the table plugin_linapro_equipment_manage_maintenance.
type MaintenanceColumns struct {
	Id          string // Maintenance record ID
	TenantId    string // Owning tenant ID, 0 means PLATFORM
	EquipmentId string // Maintained equipment ID within the same tenant
	MaintType   string // Maintenance type: 1=maintenance, 2=repair, 3=inspection, 4=other
	MaintDate   string // Maintenance date with date-only semantics
	Maintainer  string // Maintainer name
	Cost        string // Maintenance cost with two decimal places
	Content     string // Maintenance content description
	Result      string // Maintenance result
	Remark      string // Remark
	CreatedBy   string // Creator
	UpdatedBy   string // Updater
	CreatedAt   string // Creation time
	UpdatedAt   string // Update time
	DeletedAt   string // Deletion time
}

// maintenanceColumns holds the columns for the table plugin_linapro_equipment_manage_maintenance.
var maintenanceColumns = MaintenanceColumns{
	Id:          "id",
	TenantId:    "tenant_id",
	EquipmentId: "equipment_id",
	MaintType:   "maint_type",
	MaintDate:   "maint_date",
	Maintainer:  "maintainer",
	Cost:        "cost",
	Content:     "content",
	Result:      "result",
	Remark:      "remark",
	CreatedBy:   "created_by",
	UpdatedBy:   "updated_by",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
	DeletedAt:   "deleted_at",
}

// NewMaintenanceDao creates and returns a new DAO object for table data access.
func NewMaintenanceDao(handlers ...gdb.ModelHandler) *MaintenanceDao {
	return &MaintenanceDao{
		group:    "default",
		table:    "plugin_linapro_equipment_manage_maintenance",
		columns:  maintenanceColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *MaintenanceDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *MaintenanceDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *MaintenanceDao) Columns() MaintenanceColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *MaintenanceDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *MaintenanceDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *MaintenanceDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
