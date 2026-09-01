// maintenance_v1_update.go implements the controller method that updates one
// maintenance record.

package maintenance

import (
	"context"

	v1 "lina-plugin-linapro-equipment-manage/backend/api/maintenance/v1"
	maintenancesvc "lina-plugin-linapro-equipment-manage/backend/internal/service/maintenance"
)

// Update updates a maintenance record
func (c *ControllerV1) Update(ctx context.Context, req *v1.UpdateReq) (res *v1.UpdateRes, err error) {
	err = c.maintenanceSvc.Update(ctx, maintenancesvc.UpdateInput{
		Id:          req.Id,
		EquipmentId: req.EquipmentId,
		MaintType:   req.MaintType,
		MaintDate:   req.MaintDate,
		Maintainer:  req.Maintainer,
		Cost:        req.Cost,
		Content:     req.Content,
		Result:      req.Result,
		Remark:      req.Remark,
	})
	if err != nil {
		return nil, err
	}
	return &v1.UpdateRes{}, nil
}
