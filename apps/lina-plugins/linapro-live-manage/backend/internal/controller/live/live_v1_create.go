// live_v1_create.go implements the controller method that creates a new live
// content record.

package live

import (
	"context"

	v1 "lina-plugin-linapro-live-manage/backend/api/live/v1"
	livesvc "lina-plugin-linapro-live-manage/backend/internal/service/live"
)

// Create creates live content
func (c *ControllerV1) Create(ctx context.Context, req *v1.CreateReq) (res *v1.CreateRes, err error) {
	state := 0
	if req.State != nil {
		state = *req.State
	}
	isPublic := 1
	if req.IsPublic != nil {
		isPublic = *req.IsPublic
	}
	id, err := c.liveSvc.Create(ctx, livesvc.CreateInput{
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
		State:            state,
		IsPublic:         isPublic,
		StartTime:        req.StartTime,
	})
	if err != nil {
		return nil, err
	}
	return &v1.CreateRes{Id: id}, nil
}
