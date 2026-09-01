// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"time"
)

// FlowNode is the golang structure for table flow_node.
type FlowNode struct {
	Id         int64      `json:"id"         orm:"id"          description:"Flow node ID"`
	TenantId   int        `json:"tenantId"   orm:"tenant_id"   description:"Owning tenant ID, 0 means PLATFORM"`
	FlowId     int64      `json:"flowId"     orm:"flow_id"     description:"Owning approval flow ID"`
	NodeOrder  int        `json:"nodeOrder"  orm:"node_order"  description:"Approval order starting from 1, unique per flow excluding soft-deleted rows"`
	ApproverId int64      `json:"approverId" orm:"approver_id" description:"Approver user ID"`
	CreatedBy  int64      `json:"createdBy"  orm:"created_by"  description:"Creator"`
	UpdatedBy  int64      `json:"updatedBy"  orm:"updated_by"  description:"Updater"`
	CreatedAt  *time.Time `json:"createdAt"  orm:"created_at"  description:"Creation time"`
	UpdatedAt  *time.Time `json:"updatedAt"  orm:"updated_at"  description:"Update time"`
	DeletedAt  *time.Time `json:"deletedAt"  orm:"deleted_at"  description:"Deletion time"`
}
