// live_v1_list.go implements the controller method that serves the paged
// live-content list endpoint.

package live

import (
	"context"

	v1 "lina-plugin-linapro-live-manage/backend/api/live/v1"
	livesvc "lina-plugin-linapro-live-manage/backend/internal/service/live"
)

// List queries live content
func (c *ControllerV1) List(ctx context.Context, req *v1.ListReq) (res *v1.ListRes, err error) {
	out, err := c.liveSvc.List(ctx, livesvc.ListInput{
		PageNum:       req.PageNum,
		PageSize:      req.PageSize,
		Title:         req.Title,
		RoomId:        req.RoomId,
		State:         req.State,
		IsPublic:      req.IsPublic,
		LiveDateStart: req.LiveDateStart,
		LiveDateEnd:   req.LiveDateEnd,
	})
	if err != nil {
		return nil, err
	}
	items := make([]*v1.ListItem, 0, len(out.List))
	for _, item := range out.List {
		items = append(items, &v1.ListItem{
			LiveItem:      toAPILiveItem(item.LiveEntity, item.RoomName),
			CreatedByName: item.CreatedByName,
			OnlineCount:   item.OnlineCount,
			TotalViews:    item.TotalViews,
		})
	}
	return &v1.ListRes{List: items, Total: out.Total}, nil
}
