// This file declares the update-live request/response DTOs used by the
// linapro-live-manage source plugin.

package v1

import (
	"github.com/gogf/gf/v2/frame/g"
)

// Live Content Update API

// UpdateReq defines the request for updating live content.
type UpdateReq struct {
	g.Meta           `path:"/live/{id}" method:"put" tags:"LiveContents" summary:"Update live content" dc:"Update the specified live content. All fields are optional; omitted fields keep their current values. The live-room binding is immutable after creation: updating RoomId to a different room is rejected. SongList must be a valid JSON array text when provided." permission:"live:live:edit"`
	Id               int64   `json:"id" v:"required" dc:"Live content ID" eg:"1"`
	RoomId           *int64  `json:"roomId" dc:"Immutable live-room binding; passing a different room is rejected" eg:"1"`
	Title            *string `json:"title" v:"length:1,512#gf.gvalid.rule.length" dc:"Live title shown on list and cards" eg:"Sunday service live (update)"`
	Content          *string `json:"content" dc:"Live description or rich content body" eg:"<p>Updated content</p>"`
	LiveDate         *string `json:"liveDate" v:"date#gf.gvalid.rule.date" dc:"Calendar date of the live event in YYYY-MM-DD format with date-only semantics" eg:"2026-04-26"`
	PageUrl          *string `json:"pageUrl" dc:"External page link of the live event" eg:"https://example.com/live/sunday-service"`
	PushUrl          *string `json:"pushUrl" dc:"Stream push URL for the publisher" eg:"rtmp://example.com/live/sunday"`
	LiveUrl          *string `json:"liveUrl" dc:"Play URL for viewers, such as HLS or RTMP" eg:"https://example.com/hls/sunday.m3u8"`
	CoverUrl         *string `json:"coverUrl" dc:"Cover image URL for cards and previews" eg:"https://example.com/assets/cover.jpg"`
	SongName         *string `json:"songName" dc:"Featured song or hymn name" eg:"Amazing grace"`
	SongList         *string `json:"songList" dc:"Song list as a JSON array text of {name, singer, order} objects; empty means none" eg:"[{\"name\":\"Amazing grace\",\"singer\":\"Choir\",\"order\":1}]"`
	LeadSinger       *string `json:"leadSinger" dc:"Lead singer or worship team" eg:"Choir"`
	Accompaniment    *string `json:"accompaniment" dc:"Accompaniment info such as band or format" eg:"Piano"`
	Host             *string `json:"host" dc:"Host name or role" eg:"John"`
	SermonTitle      *string `json:"sermonTitle" dc:"Sermon title when the live has a standalone sermon item" eg:"Walking by faith"`
	Preacher         *string `json:"preacher" dc:"Preacher name" eg:"Pastor Wang"`
	PreacherIdentity *string `json:"preacherIdentity" dc:"Preacher identity or title, such as pastor" eg:"Pastor"`
	ScriptureRef     *string `json:"scriptureRef" dc:"Scripture reference, e.g. John 3:16" eg:"John 3:16"`
	ScriptureContent *string `json:"scriptureContent" dc:"Scripture content for display" eg:"For God so loved the world..."`
	Outline          *string `json:"outline" dc:"Sermon outline or agenda for preview" eg:"1. Source of faith; 2. Practice of faith"`
	DeviceInfo       *string `json:"deviceInfo" dc:"Device info free text, e.g. camera and microphone" eg:"Canon EOS camera, Shure SM7B microphone"`
	Reception        *string `json:"reception" dc:"Reception arrangement or contact info" eg:"Mary 13800000000"`
	State            *int    `json:"state" dc:"Live state: 0=not started, 1=ongoing, 2=finished" eg:"1"`
	IsPublic         *int    `json:"isPublic" dc:"Visibility: 1=public, 0=private" eg:"1"`
	StartTime        *int64  `json:"startTime" dc:"Actual start time of the live as Unix timestamp in milliseconds" eg:"1776756000000"`
}

// UpdateRes Live content update response
type UpdateRes struct{}
