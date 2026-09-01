// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
)

// Room is the golang structure of table plugin_linapro_live_manage_room for DAO operations like Where/Data.
type Room struct {
	g.Meta      `orm:"table:plugin_linapro_live_manage_room, do:true"`
	Id          any        // Live room ID
	TenantId    any        // Owning tenant ID, 0 means PLATFORM
	RoomCode    any        // Room code exposed to external systems, unique per tenant
	RoomName    any        // Room name
	RoomType    any        // Room type: 1=gathering, 2=event, 3=other
	Status      any        // Room status: 0=idle, 1=live, 2=disabled
	Description any        // Room description
	CreatedBy   any        // Creator
	UpdatedBy   any        // Updater
	CreatedAt   *time.Time // Creation time
	UpdatedAt   *time.Time // Update time
	DeletedAt   *time.Time // Deletion time
}
