// liveroom_v1_options.go implements the controller method that serves the
// bounded live-room option candidates endpoint.

package liveroom

import (
	"context"

	v1 "lina-plugin-linapro-live-manage/backend/api/liveroom/v1"
	liveroomsvc "lina-plugin-linapro-live-manage/backend/internal/service/liveroom"
)

// Options returns bounded live-room candidates
func (c *ControllerV1) Options(ctx context.Context, req *v1.OptionsReq) (res *v1.OptionsRes, err error) {
	includeDisabled := false
	if req.IncludeDisabled != nil {
		includeDisabled = *req.IncludeDisabled
	}
	items, err := c.roomSvc.Options(ctx, liveroomsvc.OptionsInput{
		Keyword:         req.Keyword,
		IncludeDisabled: includeDisabled,
		Limit:           req.Limit,
	})
	if err != nil {
		return nil, err
	}
	list := make([]*v1.OptionItem, 0, len(items))
	for _, item := range items {
		list = append(list, &v1.OptionItem{
			Id:       item.Id,
			RoomCode: item.RoomCode,
			RoomName: item.RoomName,
			Status:   item.Status,
		})
	}
	return &v1.OptionsRes{List: list}, nil
}
