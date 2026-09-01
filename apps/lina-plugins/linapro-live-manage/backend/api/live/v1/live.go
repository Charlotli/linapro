// This file defines shared live-content response DTOs for the linapro-live-manage API.
package v1

// LiveItem exposes live-content fields visible to callers.
type LiveItem struct {
	Id               int64  `json:"id" dc:"Live content ID" eg:"1"`
	RoomId           int64  `json:"roomId" dc:"Owning live room ID within the same tenant" eg:"1"`
	RoomName         string `json:"roomName" dc:"Live room name resolved in batch for display" eg:"Main hall room"`
	Title            string `json:"title" dc:"Live title shown on list and cards" eg:"Sunday service live"`
	Content          string `json:"content" dc:"Live description or rich content body" eg:"<p>Welcome to the online Sunday service.</p>"`
	LiveDate         string `json:"liveDate" dc:"Calendar date of the live event in YYYY-MM-DD format with date-only semantics" eg:"2026-04-26"`
	PageUrl          string `json:"pageUrl" dc:"External page link of the live event" eg:"https://example.com/live/sunday-service"`
	PushUrl          string `json:"pushUrl" dc:"Stream push URL for the publisher" eg:"rtmp://example.com/live/sunday"`
	LiveUrl          string `json:"liveUrl" dc:"Play URL for viewers, such as HLS or RTMP" eg:"https://example.com/hls/sunday.m3u8"`
	CoverUrl         string `json:"coverUrl" dc:"Cover image URL for cards and previews" eg:"https://example.com/assets/cover.jpg"`
	SongName         string `json:"songName" dc:"Featured song or hymn name" eg:"Amazing grace"`
	SongList         string `json:"songList" dc:"Song list as a JSON array text of {name, singer, order} objects; empty means none" eg:"[{\"name\":\"Amazing grace\",\"singer\":\"Choir\",\"order\":1}]"`
	LeadSinger       string `json:"leadSinger" dc:"Lead singer or worship team" eg:"Choir"`
	Accompaniment    string `json:"accompaniment" dc:"Accompaniment info such as band or format" eg:"Piano"`
	Host             string `json:"host" dc:"Host name or role" eg:"John"`
	SermonTitle      string `json:"sermonTitle" dc:"Sermon title when the live has a standalone sermon item" eg:"Walking by faith"`
	Preacher         string `json:"preacher" dc:"Preacher name" eg:"Pastor Wang"`
	PreacherIdentity string `json:"preacherIdentity" dc:"Preacher identity or title, such as pastor" eg:"Pastor"`
	ScriptureRef     string `json:"scriptureRef" dc:"Scripture reference, e.g. John 3:16" eg:"John 3:16"`
	ScriptureContent string `json:"scriptureContent" dc:"Scripture content for display" eg:"For God so loved the world..."`
	Outline          string `json:"outline" dc:"Sermon outline or agenda for preview" eg:"1. Source of faith; 2. Practice of faith"`
	DeviceInfo       string `json:"deviceInfo" dc:"Device info free text, e.g. camera and microphone" eg:"Canon EOS camera, Shure SM7B microphone"`
	Reception        string `json:"reception" dc:"Reception arrangement or contact info" eg:"Mary 13800000000"`
	State            int    `json:"state" dc:"Live state: 0=not started, 1=ongoing, 2=finished" eg:"0"`
	IsPublic         int    `json:"isPublic" dc:"Visibility: 1=public, 0=private" eg:"1"`
	StartTime        *int64 `json:"startTime" dc:"Actual start time of the live as Unix timestamp in milliseconds; null means unset" eg:"1776756000000"`
	CreatedBy        int64  `json:"createdBy" dc:"Creator user ID" eg:"1"`
	UpdatedBy        int64  `json:"updatedBy" dc:"Last updated user ID" eg:"1"`
	CreatedAt        *int64 `json:"createdAt" dc:"Creation time as Unix timestamp in milliseconds" eg:"1776756000000"`
	UpdatedAt        *int64 `json:"updatedAt" dc:"Last updated time as Unix timestamp in milliseconds" eg:"1776757800000"`
}
