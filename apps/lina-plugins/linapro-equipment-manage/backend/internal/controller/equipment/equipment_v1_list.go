// equipment_v1_list.go implements the controller method that serves the paged
// equipment list endpoint.

package equipment

import (
	"context"

	v1 "lina-plugin-linapro-equipment-manage/backend/api/equipment/v1"
	equipmentsvc "lina-plugin-linapro-equipment-manage/backend/internal/service/equipment"
)

// List queries equipment
func (c *ControllerV1) List(ctx context.Context, req *v1.ListReq) (res *v1.ListRes, err error) {
	out, err := c.equipmentSvc.List(ctx, equipmentsvc.ListInput{
		PageNum:  req.PageNum,
		PageSize: req.PageSize,
		Name:     req.Name,
		Type:     req.Type,
		Status:   req.Status,
	})
	if err != nil {
		return nil, err
	}
	items := make([]*v1.EquipmentItem, 0, len(out.List))
	for _, item := range out.List {
		projection := toAPIEquipmentItem(item)
		items = append(items, &projection)
	}
	return &v1.ListRes{List: items, Total: out.Total}, nil
}
