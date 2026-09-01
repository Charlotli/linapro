// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
)

// Record is the golang structure of table plugin_linapro_oa_approval_record for DAO operations like Where/Data.
type Record struct {
	g.Meta    `orm:"table:plugin_linapro_oa_approval_record, do:true"`
	Id        any        // Approval record ID
	TenantId  any        // Owning tenant ID, 0 means PLATFORM
	RequestId any        // Owning approval request ID
	NodeOrder any        // Related node order; 0 means request-level actions
	Action    any        // Action: 1=submit, 2=approve, 3=reject, 4=comment, 5=append approver, 6=withdraw
	ActorId   any        // Actor user ID
	Comment   any        // Reply or comment content attached to the action
	CreatedBy any        // Creator
	UpdatedBy any        // Updater
	CreatedAt *time.Time // Creation time
	UpdatedAt *time.Time // Update time
	DeletedAt *time.Time // Deletion time
}
