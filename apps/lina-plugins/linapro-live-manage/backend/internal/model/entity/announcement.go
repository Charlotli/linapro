// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"time"
)

// Announcement is the golang structure for table announcement.
type Announcement struct {
	Id        int64      `json:"id"        orm:"id"         description:"Announcement ID"`
	TenantId  int        `json:"tenantId"  orm:"tenant_id"  description:"Owning tenant ID, 0 means PLATFORM"`
	RoomId    int64      `json:"roomId"    orm:"room_id"    description:"Owning live room ID within the same tenant"`
	Title     string     `json:"title"     orm:"title"      description:"Announcement title"`
	Content   string     `json:"content"   orm:"content"    description:"Announcement body, plain text rendered with line breaks"`
	Enabled   bool       `json:"enabled"   orm:"enabled"    description:"Whether viewers can see the announcement"`
	Sort      int        `json:"sort"      orm:"sort"       description:"Display order, smaller values first"`
	CreatedBy int64      `json:"createdBy" orm:"created_by" description:"Creator"`
	UpdatedBy int64      `json:"updatedBy" orm:"updated_by" description:"Updater"`
	CreatedAt *time.Time `json:"createdAt" orm:"created_at" description:"Creation time"`
	UpdatedAt *time.Time `json:"updatedAt" orm:"updated_at" description:"Update time"`
	DeletedAt *time.Time `json:"deletedAt" orm:"deleted_at" description:"Deletion time"`
}
