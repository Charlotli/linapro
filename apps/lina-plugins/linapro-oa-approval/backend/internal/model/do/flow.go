// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
)

// Flow is the golang structure of table plugin_linapro_oa_approval_flow for DAO operations like Where/Data.
type Flow struct {
	g.Meta      `orm:"table:plugin_linapro_oa_approval_flow, do:true"`
	Id          any        // Approval flow ID
	TenantId    any        // Owning tenant ID, 0 means PLATFORM
	FlowType    any        // Flow type: 1=fund request, 2=write-off, 3=expense reimbursement
	FlowName    any        // Flow name, unique per tenant
	Description any        // Flow description
	Status      any        // Flow status: 1=enabled, 0=disabled
	CreatedBy   any        // Creator
	UpdatedBy   any        // Updater
	CreatedAt   *time.Time // Creation time
	UpdatedAt   *time.Time // Update time
	DeletedAt   *time.Time // Deletion time
	FormFields  any        // Configurable form field definitions as a JSON array
}
