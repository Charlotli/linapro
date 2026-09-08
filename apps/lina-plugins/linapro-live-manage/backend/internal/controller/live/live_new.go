// This file wires the linapro-live-manage live-content controller and shared
// response mappers.

package live

import (
	"lina-core/pkg/apitime"
	liveapi "lina-plugin-linapro-live-manage/backend/api/live"
	v1 "lina-plugin-linapro-live-manage/backend/api/live/v1"
	livesvc "lina-plugin-linapro-live-manage/backend/internal/service/live"
)

// liveDateLayout is the date-only wire format of the liveDate field.
const liveDateLayout = "2006-01-02"

// ControllerV1 is the live-content controller.
type ControllerV1 struct {
	liveSvc livesvc.Service // live-content service
}

// NewV1 creates and returns a new live-content controller instance.
func NewV1(liveSvc livesvc.Service) liveapi.ILiveV1 {
	return &ControllerV1{liveSvc: liveSvc}
}

// toAPILiveItem converts the service-layer list item into the API DTO
// projection returned through HTTP responses. Date fields are rendered in the
// date-only wire format and time points as Unix milliseconds.
func toAPILiveItem(e *livesvc.LiveEntity, roomName string) v1.LiveItem {
	if e == nil {
		return v1.LiveItem{}
	}
	return v1.LiveItem{
		Id:               e.Id,
		RoomId:           e.RoomId,
		RoomName:         roomName,
		Title:            e.Title,
		Content:          e.Content,
		LiveDate:         e.LiveDate.Format(liveDateLayout),
		PageUrl:          e.PageUrl,
		PushUrl:          e.PushUrl,
		LiveUrl:          e.LiveUrl,
		CoverUrl:         e.CoverUrl,
		SongName:         e.SongName,
		SongList:         e.SongList,
		LeadSinger:       e.LeadSinger,
		Accompaniment:    e.Accompaniment,
		Host:             e.Host,
		SermonTitle:      e.SermonTitle,
		Preacher:         e.Preacher,
		PreacherIdentity: e.PreacherIdentity,
		ScriptureRef:     e.ScriptureRef,
		ScriptureContent: e.ScriptureContent,
		Outline:          e.Outline,
		DeviceInfo:       e.DeviceInfo,
		Reception:        e.Reception,
		State:            e.State,
		IsPublic:         e.IsPublic,
		ReplayEnabled:    e.ReplayEnabled,
		StartTime:        apitime.Milli(e.StartTime),
		CreatedBy:        e.CreatedBy,
		UpdatedBy:        e.UpdatedBy,
		CreatedAt:        apitime.Milli(e.CreatedAt),
		UpdatedAt:        apitime.Milli(e.UpdatedAt),
	}
}
