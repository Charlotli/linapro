// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"time"
)

// Flow is the golang structure for table flow.
type Flow struct {
	Id          int64      `json:"id"          orm:"id"          description:"Approval flow ID"`
	TenantId    int        `json:"tenantId"    orm:"tenant_id"   description:"Owning tenant ID, 0 means PLATFORM"`
	FlowType    int        `json:"flowType"    orm:"flow_type"   description:"Flow type: 1=fund request, 2=write-off, 3=expense reimbursement"`
	FlowName    string     `json:"flowName"    orm:"flow_name"   description:"Flow name, unique per tenant"`
	Description string     `json:"description" orm:"description" description:"Flow description"`
	Status      int        `json:"status"      orm:"status"      description:"Flow status: 1=enabled, 0=disabled"`
	CreatedBy   int64      `json:"createdBy"   orm:"created_by"  description:"Creator"`
	UpdatedBy   int64      `json:"updatedBy"   orm:"updated_by"  description:"Updater"`
	CreatedAt   *time.Time `json:"createdAt"   orm:"created_at"  description:"Creation time"`
	UpdatedAt   *time.Time `json:"updatedAt"   orm:"updated_at"  description:"Update time"`
	DeletedAt   *time.Time `json:"deletedAt"   orm:"deleted_at"  description:"Deletion time"`
	FormFields  string     `json:"formFields"  orm:"form_fields" description:"Configurable form field definitions as a JSON array"`
}
