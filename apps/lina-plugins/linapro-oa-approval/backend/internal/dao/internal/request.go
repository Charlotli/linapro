// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// RequestDao is the data access object for the table plugin_linapro_oa_approval_request.
type RequestDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  RequestColumns     // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// RequestColumns defines and stores column names for the table plugin_linapro_oa_approval_request.
type RequestColumns struct {
	Id                string // Approval request ID
	TenantId          string // Owning tenant ID, 0 means PLATFORM
	FlowId            string // Originating approval flow ID
	FlowType          string // Flow type: 1=fund request, 2=write-off, 3=expense reimbursement
	Title             string // Request title
	Amount            string // Requested amount with two decimal places
	Content           string // Applicant statement content
	Attachments       string // Attachment URL list stored as a JSON string array of URL addresses, max 10 items
	Status            string // Request status: 1=pending, 2=approved, 3=rejected, 4=withdrawn
	CurrentNodeOrder  string // Pending node order while the request is pending
	CurrentApproverId string // Pending approver user ID while the request is pending
	NodeSnapshot      string // Frozen node list as a JSON array with order, approverId and approverName fields at submit time
	ApplicantId       string // Applicant user ID
	CreatedBy         string // Creator
	UpdatedBy         string // Updater
	CreatedAt         string // Creation time
	UpdatedAt         string // Update time
	DeletedAt         string // Deletion time
	FormData          string // Submitted dynamic form values as a JSON object keyed by field key
	FormSnapshot      string // Frozen form field definitions as a JSON array at submit time
}

// requestColumns holds the columns for the table plugin_linapro_oa_approval_request.
var requestColumns = RequestColumns{
	Id:                "id",
	TenantId:          "tenant_id",
	FlowId:            "flow_id",
	FlowType:          "flow_type",
	Title:             "title",
	Amount:            "amount",
	Content:           "content",
	Attachments:       "attachments",
	Status:            "status",
	CurrentNodeOrder:  "current_node_order",
	CurrentApproverId: "current_approver_id",
	NodeSnapshot:      "node_snapshot",
	ApplicantId:       "applicant_id",
	CreatedBy:         "created_by",
	UpdatedBy:         "updated_by",
	CreatedAt:         "created_at",
	UpdatedAt:         "updated_at",
	DeletedAt:         "deleted_at",
	FormData:          "form_data",
	FormSnapshot:      "form_snapshot",
}

// NewRequestDao creates and returns a new DAO object for table data access.
func NewRequestDao(handlers ...gdb.ModelHandler) *RequestDao {
	return &RequestDao{
		group:    "default",
		table:    "plugin_linapro_oa_approval_request",
		columns:  requestColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *RequestDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *RequestDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *RequestDao) Columns() RequestColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *RequestDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *RequestDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *RequestDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
