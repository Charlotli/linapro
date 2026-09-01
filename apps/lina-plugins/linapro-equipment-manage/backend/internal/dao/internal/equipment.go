// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// EquipmentDao is the data access object for the table plugin_linapro_equipment_manage_equipment.
type EquipmentDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  EquipmentColumns   // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// EquipmentColumns defines and stores column names for the table plugin_linapro_equipment_manage_equipment.
type EquipmentColumns struct {
	Id            string // Equipment ID
	TenantId      string // Owning tenant ID, 0 means PLATFORM
	EquipmentCode string // Equipment code, unique per tenant
	EquipmentName string // Equipment name
	EquipmentType string // Equipment type: 1=office, 2=it, 3=production, 4=other
	BrandModel    string // Brand and model
	PurchaseDate  string // Purchase date with date-only semantics
	PurchasePrice string // Purchase price with two decimal places
	Location      string // Storage location
	Owner         string // Responsible person
	Status        string // Equipment status: 1=in use, 2=idle, 3=under repair, 4=scrapped
	Remark        string // Remark
	CreatedBy     string // Creator
	UpdatedBy     string // Updater
	CreatedAt     string // Creation time
	UpdatedAt     string // Update time
	DeletedAt     string // Deletion time
}

// equipmentColumns holds the columns for the table plugin_linapro_equipment_manage_equipment.
var equipmentColumns = EquipmentColumns{
	Id:            "id",
	TenantId:      "tenant_id",
	EquipmentCode: "equipment_code",
	EquipmentName: "equipment_name",
	EquipmentType: "equipment_type",
	BrandModel:    "brand_model",
	PurchaseDate:  "purchase_date",
	PurchasePrice: "purchase_price",
	Location:      "location",
	Owner:         "owner",
	Status:        "status",
	Remark:        "remark",
	CreatedBy:     "created_by",
	UpdatedBy:     "updated_by",
	CreatedAt:     "created_at",
	UpdatedAt:     "updated_at",
	DeletedAt:     "deleted_at",
}

// NewEquipmentDao creates and returns a new DAO object for table data access.
func NewEquipmentDao(handlers ...gdb.ModelHandler) *EquipmentDao {
	return &EquipmentDao{
		group:    "default",
		table:    "plugin_linapro_equipment_manage_equipment",
		columns:  equipmentColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *EquipmentDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *EquipmentDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *EquipmentDao) Columns() EquipmentColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *EquipmentDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *EquipmentDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *EquipmentDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
