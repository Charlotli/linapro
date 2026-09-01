// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"time"
)

// Record is the golang structure for table record.
type Record struct {
	Id        int64      `json:"id"        orm:"id"         description:"Approval record ID"`
	TenantId  int        `json:"tenantId"  orm:"tenant_id"  description:"Owning tenant ID, 0 means PLATFORM"`
	RequestId int64      `json:"requestId" orm:"request_id" description:"Owning approval request ID"`
	NodeOrder int        `json:"nodeOrder" orm:"node_order" description:"Related node order; 0 means request-level actions"`
	Action    int        `json:"action"    orm:"action"     description:"Action: 1=submit, 2=approve, 3=reject, 4=comment, 5=append approver, 6=withdraw"`
	ActorId   int64      `json:"actorId"   orm:"actor_id"   description:"Actor user ID"`
	Comment   string     `json:"comment"   orm:"comment"    description:"Reply or comment content attached to the action"`
	CreatedBy int64      `json:"createdBy" orm:"created_by" description:"Creator"`
	UpdatedBy int64      `json:"updatedBy" orm:"updated_by" description:"Updater"`
	CreatedAt *time.Time `json:"createdAt" orm:"created_at" description:"Creation time"`
	UpdatedAt *time.Time `json:"updatedAt" orm:"updated_at" description:"Update time"`
	DeletedAt *time.Time `json:"deletedAt" orm:"deleted_at" description:"Deletion time"`
}
