// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
)

// FlowNode is the golang structure of table plugin_linapro_oa_approval_flow_node for DAO operations like Where/Data.
type FlowNode struct {
	g.Meta     `orm:"table:plugin_linapro_oa_approval_flow_node, do:true"`
	Id         any        // Flow node ID
	TenantId   any        // Owning tenant ID, 0 means PLATFORM
	FlowId     any        // Owning approval flow ID
	NodeOrder  any        // Approval order starting from 1, unique per flow excluding soft-deleted rows
	ApproverId any        // Approver user ID
	CreatedBy  any        // Creator
	UpdatedBy  any        // Updater
	CreatedAt  *time.Time // Creation time
	UpdatedAt  *time.Time // Update time
	DeletedAt  *time.Time // Deletion time
}
