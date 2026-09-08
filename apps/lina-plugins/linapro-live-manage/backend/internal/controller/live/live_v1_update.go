// live_v1_update.go implements the controller method that updates an existing
// live content record.

package live

import (
	"context"

	v1 "lina-plugin-linapro-live-manage/backend/api/live/v1"
	livesvc "lina-plugin-linapro-live-manage/backend/internal/service/live"
)

// Update updates live content
func (c *ControllerV1) Update(ctx context.Context, req *v1.UpdateReq) (res *v1.UpdateRes, err error) {
	err = c.liveSvc.Update(ctx, livesvc.UpdateInput{
		Id:               req.Id,
		RoomId:           req.RoomId,
		Title:            req.Title,
		Content:          req.Content,
		LiveDate:         req.LiveDate,
		PageUrl:          req.PageUrl,
		PushUrl:          req.PushUrl,
		LiveUrl:          req.LiveUrl,
		CoverUrl:         req.CoverUrl,
		SongName:         req.SongName,
		SongList:         req.SongList,
		LeadSinger:       req.LeadSinger,
		Accompaniment:    req.Accompaniment,
		Host:             req.Host,
		SermonTitle:      req.SermonTitle,
		Preacher:         req.Preacher,
		PreacherIdentity: req.PreacherIdentity,
		ScriptureRef:     req.ScriptureRef,
		ScriptureContent: req.ScriptureContent,
		Outline:          req.Outline,
		DeviceInfo:       req.DeviceInfo,
		Reception:        req.Reception,
		State:            req.State,
		IsPublic:         req.IsPublic,
		StartTime:        req.StartTime,
		ReplayEnabled:    req.ReplayEnabled,
	})
	if err != nil {
		return nil, err
	}
	return &v1.UpdateRes{}, nil
}
