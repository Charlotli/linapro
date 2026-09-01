// equipment_v1_get.go implements the controller method that serves the
// equipment detail endpoint.

package equipment

import (
	"context"

	v1 "lina-plugin-linapro-equipment-manage/backend/api/equipment/v1"
)

// Get returns equipment details
func (c *ControllerV1) Get(ctx context.Context, req *v1.GetReq) (res *v1.GetRes, err error) {
	item, err := c.equipmentSvc.GetById(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	projection := toAPIEquipmentItem(item)
	return &v1.GetRes{EquipmentItem: projection}, nil
}
