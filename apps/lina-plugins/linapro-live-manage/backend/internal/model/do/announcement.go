// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
)

// Announcement is the golang structure of table plugin_linapro_live_manage_announcement for DAO operations like Where/Data.
type Announcement struct {
	g.Meta    `orm:"table:plugin_linapro_live_manage_announcement, do:true"`
	Id        any        // Announcement ID
	TenantId  any        // Owning tenant ID, 0 means PLATFORM
	RoomId    any        // Owning live room ID within the same tenant
	Title     any        // Announcement title
	Content   any        // Announcement body, plain text rendered with line breaks
	Enabled   any        // Whether viewers can see the announcement
	Sort      any        // Display order, smaller values first
	CreatedBy any        // Creator
	UpdatedBy any        // Updater
	CreatedAt *time.Time // Creation time
	UpdatedAt *time.Time // Update time
	DeletedAt *time.Time // Deletion time
}
