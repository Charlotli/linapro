// This file defines shared play response DTOs for the linapro-live-manage
// public viewer API.

package v1

// PlayItem exposes the minimal live-content projection visible to anonymous
// viewers. Administrative fields such as the push URL are intentionally
// excluded from this projection.
type PlayItem struct {
	RoomId           int64  `json:"roomId" dc:"Owning live room ID" eg:"1"`
	RoomCode         string `json:"roomCode" dc:"Live room code from the play request" eg:"MAIN-HALL"`
	RoomName         string `json:"roomName" dc:"Live room name" eg:"Main hall room"`
	LiveId           int64  `json:"liveId" dc:"Live content ID; 0 means no live record exists for the requested state" eg:"1"`
	Title            string `json:"title" dc:"Live title shown as the page headline" eg:"Sunday service live"`
	CoverUrl         string `json:"coverUrl" dc:"Cover image URL for the player placeholder and preview" eg:"https://example.com/assets/cover.jpg"`
	LiveUrl          string `json:"liveUrl" dc:"HLS play URL for viewers; empty for not-started lives so the stream address is not exposed early" eg:"https://example.com/hls/sunday.m3u8"`
	LiveDate         string `json:"liveDate" dc:"Calendar date of the live event in YYYY-MM-DD format with date-only semantics" eg:"2026-04-26"`
	StartTime        *int64 `json:"startTime" dc:"Actual start time of the live as Unix timestamp in milliseconds; null means the live has not started" eg:"1776756000000"`
	State            int    `json:"state" dc:"Live state: 0=not started, 1=ongoing, 2=finished" eg:"1"`
	Host             string `json:"host" dc:"Host name or role; empty means unset" eg:"John"`
	SongName         string `json:"songName" dc:"Featured song or hymn name; empty means unset" eg:"Amazing grace"`
	SongList         string `json:"songList" dc:"Song list as a JSON array text of {name, singer, order} objects; empty means none" eg:"[{\"name\":\"Amazing grace\",\"singer\":\"Choir\",\"order\":1}]"`
	LeadSinger       string `json:"leadSinger" dc:"Lead singer or worship team; empty means unset" eg:"Choir"`
	Accompaniment    string `json:"accompaniment" dc:"Accompaniment info such as band or format; empty means unset" eg:"Piano"`
	SermonTitle      string `json:"sermonTitle" dc:"Sermon title when the live has a standalone sermon item; empty means unset" eg:"Walking by faith"`
	Preacher         string `json:"preacher" dc:"Preacher name; empty means unset" eg:"Pastor Wang"`
	PreacherIdentity string `json:"preacherIdentity" dc:"Preacher identity or title, such as pastor; empty means unset" eg:"Pastor"`
	ScriptureRef     string `json:"scriptureRef" dc:"Scripture reference, e.g. John 3:16; empty means unset" eg:"John 3:16"`
	ScriptureContent string `json:"scriptureContent" dc:"Scripture content for display; empty means unset" eg:"For God so loved the world..."`
	Outline          string `json:"outline" dc:"Sermon outline or agenda for preview; empty means unset" eg:"1. Source of faith; 2. Practice of faith"`
}
