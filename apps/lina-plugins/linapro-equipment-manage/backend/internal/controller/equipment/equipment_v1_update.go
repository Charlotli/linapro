// equipment_v1_update.go implements the controller method that updates one
// equipment record.

package equipment

import (
	"context"

	v1 "lina-plugin-linapro-equipment-manage/backend/api/equipment/v1"
	equipmentsvc "lina-plugin-linapro-equipment-manage/backend/internal/service/equipment"
)

// Update updates equipment
func (c *ControllerV1) Update(ctx context.Context, req *v1.UpdateReq) (res *v1.UpdateRes, err error) {
	err = c.equipmentSvc.Update(ctx, equipmentsvc.UpdateInput{
		Id:            req.Id,
		EquipmentCode: req.EquipmentCode,
		EquipmentName: req.EquipmentName,
		EquipmentType: req.EquipmentType,
		BrandModel:    req.BrandModel,
		PurchaseDate:  req.PurchaseDate,
		PurchasePrice: req.PurchasePrice,
		Location:      req.Location,
		Owner:         req.Owner,
		Status:        req.Status,
		Remark:        req.Remark,
	})
	if err != nil {
		return nil, err
	}
	return &v1.UpdateRes{}, nil
}
