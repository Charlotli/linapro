// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"time"
)

// Live is the golang structure for table live.
type Live struct {
	Id               int64      `json:"id"               orm:"id"                description:"Live content ID"`
	TenantId         int        `json:"tenantId"         orm:"tenant_id"         description:"Owning tenant ID, 0 means PLATFORM"`
	RoomId           int64      `json:"roomId"           orm:"room_id"           description:"Owning live room ID within the same tenant"`
	Title            string     `json:"title"            orm:"title"             description:"Live title shown on list and cards"`
	Content          string     `json:"content"          orm:"content"           description:"Live description or rich content body"`
	LiveDate         time.Time  `json:"liveDate"         orm:"live_date"         description:"Calendar date of the live event for list sorting and filtering"`
	PageUrl          string     `json:"pageUrl"          orm:"page_url"          description:"External page link of the live event"`
	PushUrl          string     `json:"pushUrl"          orm:"push_url"          description:"Stream push URL for the publisher"`
	LiveUrl          string     `json:"liveUrl"          orm:"live_url"          description:"Play URL for viewers, such as HLS or RTMP"`
	CoverUrl         string     `json:"coverUrl"         orm:"cover_url"         description:"Cover image URL for cards and previews"`
	SongName         string     `json:"songName"         orm:"song_name"         description:"Featured song or hymn name"`
	SongList         string     `json:"songList"         orm:"song_list"         description:"Song list as a JSON array text, e.g. [{\"name\":\"Song A\",\"singer\":\"Alice\",\"order\":1}]"`
	LeadSinger       string     `json:"leadSinger"       orm:"lead_singer"       description:"Lead singer or worship team"`
	Accompaniment    string     `json:"accompaniment"    orm:"accompaniment"     description:"Accompaniment info such as band or format"`
	Host             string     `json:"host"             orm:"host"              description:"Host name or role"`
	SermonTitle      string     `json:"sermonTitle"      orm:"sermon_title"      description:"Sermon title when the live has a standalone sermon item"`
	Preacher         string     `json:"preacher"         orm:"preacher"          description:"Preacher name"`
	PreacherIdentity string     `json:"preacherIdentity" orm:"preacher_identity" description:"Preacher identity or title, such as pastor"`
	ScriptureRef     string     `json:"scriptureRef"     orm:"scripture_ref"     description:"Scripture reference, e.g. John 3:16"`
	ScriptureContent string     `json:"scriptureContent" orm:"scripture_content" description:"Scripture content for display"`
	Outline          string     `json:"outline"          orm:"outline"           description:"Sermon outline or agenda for preview"`
	DeviceInfo       string     `json:"deviceInfo"       orm:"device_info"       description:"Device info free text, e.g. camera and microphone"`
	Reception        string     `json:"reception"        orm:"reception"         description:"Reception arrangement or contact info"`
	State            int        `json:"state"            orm:"state"             description:"Live state: 0=not started, 1=ongoing, 2=finished"`
	IsPublic         int        `json:"isPublic"         orm:"is_public"         description:"Visibility: 1=public, 0=private"`
	StartTime        *time.Time `json:"startTime"        orm:"start_time"        description:"Actual start time of the live for schedule and reminders"`
	CreatedBy        int64      `json:"createdBy"        orm:"created_by"        description:"Creator"`
	UpdatedBy        int64      `json:"updatedBy"        orm:"updated_by"        description:"Updater"`
	CreatedAt        *time.Time `json:"createdAt"        orm:"created_at"        description:"Creation time"`
	UpdatedAt        *time.Time `json:"updatedAt"        orm:"updated_at"        description:"Update time"`
	DeletedAt        *time.Time `json:"deletedAt"        orm:"deleted_at"        description:"Deletion time"`
}
