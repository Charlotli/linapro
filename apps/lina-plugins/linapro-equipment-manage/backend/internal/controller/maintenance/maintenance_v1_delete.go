// maintenance_v1_delete.go implements the controller method that deletes
// maintenance records by ID list.

package maintenance

import (
	"context"

	v1 "lina-plugin-linapro-equipment-manage/backend/api/maintenance/v1"
)

// Delete deletes maintenance records by ID list.
func (c *ControllerV1) Delete(ctx context.Context, req *v1.DeleteReq) (res *v1.DeleteRes, err error) {
	err = c.maintenanceSvc.Delete(ctx, req.Ids)
	if err != nil {
		return nil, err
	}
	return &v1.DeleteRes{}, nil
}
