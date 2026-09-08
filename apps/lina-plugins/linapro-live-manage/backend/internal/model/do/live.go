// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"time"

	"github.com/gogf/gf/v2/frame/g"
)

// Live is the golang structure of table plugin_linapro_live_manage_live for DAO operations like Where/Data.
type Live struct {
	g.Meta           `orm:"table:plugin_linapro_live_manage_live, do:true"`
	Id               any        // Live content ID
	TenantId         any        // Owning tenant ID, 0 means PLATFORM
	RoomId           any        // Owning live room ID within the same tenant
	Title            any        // Live title shown on list and cards
	Content          any        // Live description or rich content body
	LiveDate         any        // Calendar date of the live event for list sorting and filtering
	PageUrl          any        // External page link of the live event
	PushUrl          any        // Stream push URL for the publisher
	LiveUrl          any        // Play URL for viewers, such as HLS or RTMP
	CoverUrl         any        // Cover image URL for cards and previews
	SongName         any        // Featured song or hymn name
	SongList         any        // Song list as a JSON array text, e.g. [{"name":"Song A","singer":"Alice","order":1}]
	LeadSinger       any        // Lead singer or worship team
	Accompaniment    any        // Accompaniment info such as band or format
	Host             any        // Host name or role
	SermonTitle      any        // Sermon title when the live has a standalone sermon item
	Preacher         any        // Preacher name
	PreacherIdentity any        // Preacher identity or title, such as pastor
	ScriptureRef     any        // Scripture reference, e.g. John 3:16
	ScriptureContent any        // Scripture content for display
	Outline          any        // Sermon outline or agenda for preview
	DeviceInfo       any        // Device info free text, e.g. camera and microphone
	Reception        any        // Reception arrangement or contact info
	State            any        // Live state: 0=not started, 1=ongoing, 2=finished
	IsPublic         any        // Visibility: 1=public, 0=private
	StartTime        *time.Time // Actual start time of the live for schedule and reminders
	CreatedBy        any        // Creator
	UpdatedBy        any        // Updater
	CreatedAt        *time.Time // Creation time
	UpdatedAt        *time.Time // Update time
	DeletedAt        *time.Time // Deletion time
	ReplayEnabled    any        // Whether the finished live is offered as a public replay; ongoing lives are unaffected
}
