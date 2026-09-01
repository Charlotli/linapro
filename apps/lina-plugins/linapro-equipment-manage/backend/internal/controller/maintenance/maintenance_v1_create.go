// maintenance_v1_create.go implements the controller method that creates one
// maintenance record.

package maintenance

import (
	"context"

	v1 "lina-plugin-linapro-equipment-manage/backend/api/maintenance/v1"
	maintenancesvc "lina-plugin-linapro-equipment-manage/backend/internal/service/maintenance"
)

// Create creates a maintenance record
func (c *ControllerV1) Create(ctx context.Context, req *v1.CreateReq) (res *v1.CreateRes, err error) {
	id, err := c.maintenanceSvc.Create(ctx, maintenancesvc.CreateInput{
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
	return &v1.CreateRes{Id: id}, nil
}
