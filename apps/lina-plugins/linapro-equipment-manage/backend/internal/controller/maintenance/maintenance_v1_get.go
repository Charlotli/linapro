// maintenance_v1_get.go implements the controller method that serves the
// maintenance record detail endpoint.

package maintenance

import (
	"context"

	v1 "lina-plugin-linapro-equipment-manage/backend/api/maintenance/v1"
	maintenancesvc "lina-plugin-linapro-equipment-manage/backend/internal/service/maintenance"
)

// Get returns maintenance record details
func (c *ControllerV1) Get(ctx context.Context, req *v1.GetReq) (res *v1.GetRes, err error) {
	item, err := c.maintenanceSvc.GetById(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	projection := toAPIMaintenanceItem(&maintenancesvc.MaintenanceItem{MaintenanceEntity: item})
	return &v1.GetRes{MaintenanceItem: projection}, nil
}
