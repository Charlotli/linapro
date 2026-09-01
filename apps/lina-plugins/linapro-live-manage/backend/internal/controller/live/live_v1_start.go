// live_v1_start.go implements the controller method that starts one live
// content record.

package live

import (
	"context"

	v1 "lina-plugin-linapro-live-manage/backend/api/live/v1"
)

// Start starts live content
func (c *ControllerV1) Start(ctx context.Context, req *v1.StartReq) (res *v1.StartRes, err error) {
	err = c.liveSvc.Start(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &v1.StartRes{}, nil
}
