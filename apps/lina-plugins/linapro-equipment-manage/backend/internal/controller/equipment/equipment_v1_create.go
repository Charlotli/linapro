// equipment_v1_create.go implements the controller method that creates one
// equipment record.

package equipment

import (
	"context"

	v1 "lina-plugin-linapro-equipment-manage/backend/api/equipment/v1"
	equipmentsvc "lina-plugin-linapro-equipment-manage/backend/internal/service/equipment"
)

// Create creates equipment
func (c *ControllerV1) Create(ctx context.Context, req *v1.CreateReq) (res *v1.CreateRes, err error) {
	status := 1
	if req.Status != nil {
		status = *req.Status
	}
	id, err := c.equipmentSvc.Create(ctx, equipmentsvc.CreateInput{
		EquipmentCode: req.EquipmentCode,
		EquipmentName: req.EquipmentName,
		EquipmentType: req.EquipmentType,
		BrandModel:    req.BrandModel,
		PurchaseDate:  req.PurchaseDate,
		PurchasePrice: req.PurchasePrice,
		Location:      req.Location,
		Owner:         req.Owner,
		Status:        status,
		Remark:        req.Remark,
	})
	if err != nil {
		return nil, err
	}
	return &v1.CreateRes{Id: id}, nil
}
