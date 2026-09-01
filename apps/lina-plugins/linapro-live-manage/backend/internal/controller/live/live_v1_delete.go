// live_v1_delete.go implements the controller method that deletes live content
// records by ID list.

package live

import (
	"context"

	v1 "lina-plugin-linapro-live-manage/backend/api/live/v1"
)

// Delete deletes live content by ID list.
func (c *ControllerV1) Delete(ctx context.Context, req *v1.DeleteReq) (res *v1.DeleteRes, err error) {
	err = c.liveSvc.Delete(ctx, req.Ids)
	if err != nil {
		return nil, err
	}
	return &v1.DeleteRes{}, nil
}
