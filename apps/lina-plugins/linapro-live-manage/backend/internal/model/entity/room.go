// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"time"
)

// Room is the golang structure for table room.
type Room struct {
	Id          int64      `json:"id"          orm:"id"          description:"Live room ID"`
	TenantId    int        `json:"tenantId"    orm:"tenant_id"   description:"Owning tenant ID, 0 means PLATFORM"`
	RoomCode    string     `json:"roomCode"    orm:"room_code"   description:"Room code exposed to external systems, unique per tenant"`
	RoomName    string     `json:"roomName"    orm:"room_name"   description:"Room name"`
	RoomType    int        `json:"roomType"    orm:"room_type"   description:"Room type: 1=gathering, 2=event, 3=other"`
	Status      int        `json:"status"      orm:"status"      description:"Room status: 0=idle, 1=live, 2=disabled"`
	Description string     `json:"description" orm:"description" description:"Room description"`
	CreatedBy   int64      `json:"createdBy"   orm:"created_by"  description:"Creator"`
	UpdatedBy   int64      `json:"updatedBy"   orm:"updated_by"  description:"Updater"`
	CreatedAt   *time.Time `json:"createdAt"   orm:"created_at"  description:"Creation time"`
	UpdatedAt   *time.Time `json:"updatedAt"   orm:"updated_at"  description:"Update time"`
	DeletedAt   *time.Time `json:"deletedAt"   orm:"deleted_at"  description:"Deletion time"`
}
