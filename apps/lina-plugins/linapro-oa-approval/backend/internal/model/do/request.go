// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
)

// Request is the golang structure of table plugin_linapro_oa_approval_request for DAO operations like Where/Data.
type Request struct {
	g.Meta            `orm:"table:plugin_linapro_oa_approval_request, do:true"`
	Id                any        // Approval request ID
	TenantId          any        // Owning tenant ID, 0 means PLATFORM
	FlowId            any        // Originating approval flow ID
	FlowType          any        // Flow type: 1=fund request, 2=write-off, 3=expense reimbursement
	Title             any        // Request title
	Amount            any        // Requested amount with two decimal places
	Content           any        // Applicant statement content
	Attachments       any        // Attachment URL list stored as a JSON string array of URL addresses, max 10 items
	Status            any        // Request status: 1=pending, 2=approved, 3=rejected, 4=withdrawn
	CurrentNodeOrder  any        // Pending node order while the request is pending
	CurrentApproverId any        // Pending approver user ID while the request is pending
	NodeSnapshot      any        // Frozen node list as a JSON array with order, approverId and approverName fields at submit time
	ApplicantId       any        // Applicant user ID
	CreatedBy         any        // Creator
	UpdatedBy         any        // Updater
	CreatedAt         *time.Time // Creation time
	UpdatedAt         *time.Time // Update time
	DeletedAt         *time.Time // Deletion time
	FormData          any        // Submitted dynamic form values as a JSON object keyed by field key
	FormSnapshot      any        // Frozen form field definitions as a JSON array at submit time
}
