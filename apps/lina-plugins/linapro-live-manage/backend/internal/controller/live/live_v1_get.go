// live_v1_get.go implements the controller method that serves the live-content
// detail endpoint.

package live

import (
	"context"

	v1 "lina-plugin-linapro-live-manage/backend/api/live/v1"
)

// Get returns live content details
func (c *ControllerV1) Get(ctx context.Context, req *v1.GetReq) (res *v1.GetRes, err error) {
	item, err := c.liveSvc.GetById(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &v1.GetRes{
		LiveItem:      toAPILiveItem(item.LiveEntity, item.RoomName),
		CreatedByName: item.CreatedByName,
	}, nil
}
