// live_v1_stop.go implements the controller method that stops one live
// content record.

package live

import (
	"context"

	v1 "lina-plugin-linapro-live-manage/backend/api/live/v1"
)

// Stop stops live content
func (c *ControllerV1) Stop(ctx context.Context, req *v1.StopReq) (res *v1.StopRes, err error) {
	err = c.liveSvc.Stop(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &v1.StopRes{}, nil
}
