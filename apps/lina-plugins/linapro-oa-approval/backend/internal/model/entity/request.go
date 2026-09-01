// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"time"
)

// Request is the golang structure for table request.
type Request struct {
	Id                int64      `json:"id"                orm:"id"                  description:"Approval request ID"`
	TenantId          int        `json:"tenantId"          orm:"tenant_id"           description:"Owning tenant ID, 0 means PLATFORM"`
	FlowId            int64      `json:"flowId"            orm:"flow_id"             description:"Originating approval flow ID"`
	FlowType          int        `json:"flowType"          orm:"flow_type"           description:"Flow type: 1=fund request, 2=write-off, 3=expense reimbursement"`
	Title             string     `json:"title"             orm:"title"               description:"Request title"`
	Amount            float64    `json:"amount"            orm:"amount"              description:"Requested amount with two decimal places"`
	Content           string     `json:"content"           orm:"content"             description:"Applicant statement content"`
	Attachments       string     `json:"attachments"       orm:"attachments"         description:"Attachment URL list stored as a JSON string array of URL addresses, max 10 items"`
	Status            int        `json:"status"            orm:"status"              description:"Request status: 1=pending, 2=approved, 3=rejected, 4=withdrawn"`
	CurrentNodeOrder  int        `json:"currentNodeOrder"  orm:"current_node_order"  description:"Pending node order while the request is pending"`
	CurrentApproverId int64      `json:"currentApproverId" orm:"current_approver_id" description:"Pending approver user ID while the request is pending"`
	NodeSnapshot      string     `json:"nodeSnapshot"      orm:"node_snapshot"       description:"Frozen node list as a JSON array with order, approverId and approverName fields at submit time"`
	ApplicantId       int64      `json:"applicantId"       orm:"applicant_id"        description:"Applicant user ID"`
	CreatedBy         int64      `json:"createdBy"         orm:"created_by"          description:"Creator"`
	UpdatedBy         int64      `json:"updatedBy"         orm:"updated_by"          description:"Updater"`
	CreatedAt         *time.Time `json:"createdAt"         orm:"created_at"          description:"Creation time"`
	UpdatedAt         *time.Time `json:"updatedAt"         orm:"updated_at"          description:"Update time"`
	DeletedAt         *time.Time `json:"deletedAt"         orm:"deleted_at"          description:"Deletion time"`
	FormData          string     `json:"formData"          orm:"form_data"           description:"Submitted dynamic form values as a JSON object keyed by field key"`
	FormSnapshot      string     `json:"formSnapshot"      orm:"form_snapshot"       description:"Frozen form field definitions as a JSON array at submit time"`
}
