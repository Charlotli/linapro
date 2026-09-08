// play_v1_replays.go implements the controller method that serves the
// anonymous paginated replay-library endpoint consumed by the plugin H5
// viewer page.

package play

import (
	"context"

	v1 "lina-plugin-linapro-live-manage/backend/api/play/v1"
	playsvc "lina-plugin-linapro-live-manage/backend/internal/service/play"
)

// Replays lists one page of the room's public replay library.
func (c *ControllerV1) Replays(ctx context.Context, req *v1.ReplaysReq) (res *v1.ReplaysRes, err error) {
	input := playsvc.ReplaysInput{
		RoomCode: req.RoomCode,
		Page:     req.Page,
		PageSize: req.PageSize,
	}
	if req.TenantId != nil {
		tenantID := *req.TenantId
		input.TenantId = &tenantID
	}
	items, total, err := c.playSvc.Replays(ctx, input)
	if err != nil {
		return nil, err
	}
	list := make([]v1.ReplayItem, 0, len(items))
	for _, item := range items {
		var startTime *int64
		if item.StartTime != nil {
			milli := item.StartTime.UnixMilli()
			startTime = &milli
		}
		list = append(list, v1.ReplayItem{
			LiveId:    item.LiveId,
			Title:     item.Title,
			CoverUrl:  item.CoverUrl,
			LiveUrl:   item.LiveUrl,
			LiveDate:  item.LiveDate.Format(playDateLayout),
			StartTime: startTime,
		})
	}
	return &v1.ReplaysRes{List: list, Total: total}, nil
}
