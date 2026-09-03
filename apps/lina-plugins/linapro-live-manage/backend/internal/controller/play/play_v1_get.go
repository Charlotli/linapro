// play_v1_get.go implements the controller method that serves the anonymous
// public play-info endpoint consumed by the plugin H5 viewer page.

package play

import (
	"context"

	v1 "lina-plugin-linapro-live-manage/backend/api/play/v1"
	playsvc "lina-plugin-linapro-live-manage/backend/internal/service/play"
)

// playDateLayout is the date-only wire format of the liveDate field.
const playDateLayout = "2006-01-02"

// Get resolves one publicly visible live for anonymous viewers.
func (c *ControllerV1) Get(ctx context.Context, req *v1.GetReq) (res *v1.GetRes, err error) {
	info, err := c.playSvc.Get(ctx, toServiceInput(req))
	if err != nil {
		return nil, err
	}
	return &v1.GetRes{PlayItem: toAPIPlayItem(info)}, nil
}

// toServiceInput converts the request DTO into the service input. The tenant
// parameter is nil when absent so the service can enforce its own presence
// rules per tenant-capability state.
func toServiceInput(req *v1.GetReq) playsvc.GetInput {
	input := playsvc.GetInput{
		RoomCode: req.RoomCode,
	}
	if req.TenantId != nil {
		tenantID := *req.TenantId
		input.TenantId = &tenantID
	}
	return input
}

// toAPIPlayItem converts the service projection into the API DTO. Date fields
// render in the date-only wire format and time points as Unix milliseconds.
func toAPIPlayItem(info *playsvc.PlayInfo) v1.PlayItem {
	if info == nil {
		return v1.PlayItem{}
	}
	var startTime *int64
	if info.StartTime != nil {
		milli := info.StartTime.UnixMilli()
		startTime = &milli
	}
	return v1.PlayItem{
		RoomId:           info.RoomId,
		RoomCode:         info.RoomCode,
		RoomName:         info.RoomName,
		LiveId:           info.LiveId,
		Title:            info.Title,
		CoverUrl:         info.CoverUrl,
		LiveUrl:          info.LiveUrl,
		LiveDate:         info.LiveDate.Format(playDateLayout),
		StartTime:        startTime,
		State:            info.State,
		Host:             info.Host,
		SongName:         info.SongName,
		SongList:         info.SongList,
		LeadSinger:       info.LeadSinger,
		Accompaniment:    info.Accompaniment,
		SermonTitle:      info.SermonTitle,
		Preacher:         info.Preacher,
		PreacherIdentity: info.PreacherIdentity,
		ScriptureRef:     info.ScriptureRef,
		ScriptureContent: info.ScriptureContent,
		Outline:          info.Outline,
	}
}
