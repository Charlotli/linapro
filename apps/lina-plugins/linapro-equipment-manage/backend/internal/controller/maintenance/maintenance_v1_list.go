// maintenance_v1_list.go implements the controller method that serves the
// paged maintenance record list endpoint.

package maintenance

import (
	"context"

	v1 "lina-plugin-linapro-equipment-manage/backend/api/maintenance/v1"
	maintenancesvc "lina-plugin-linapro-equipment-manage/backend/internal/service/maintenance"
)

// List queries maintenance records
func (c *ControllerV1) List(ctx context.Context, req *v1.ListReq) (res *v1.ListRes, err error) {
	out, err := c.maintenanceSvc.List(ctx, maintenancesvc.ListInput{
		PageNum:     req.PageNum,
		PageSize:    req.PageSize,
		EquipmentId: req.EquipmentId,
		MaintType:   req.MaintType,
		DateStart:   req.DateStart,
		DateEnd:     req.DateEnd,
	})
	if err != nil {
		return nil, err
	}
	items := make([]*v1.MaintenanceItem, 0, len(out.List))
	for _, item := range out.List {
		projection := toAPIMaintenanceItem(item)
		items = append(items, &projection)
	}
	return &v1.ListRes{List: items, Total: out.Total}, nil
}
