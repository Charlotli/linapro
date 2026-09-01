// equipment_v1_options.go implements the controller method that serves the
// bounded equipment candidate endpoint.

package equipment

import (
	"context"

	v1 "lina-plugin-linapro-equipment-manage/backend/api/equipment/v1"
	equipmentsvc "lina-plugin-linapro-equipment-manage/backend/internal/service/equipment"
)

// Options returns bounded equipment candidates
func (c *ControllerV1) Options(ctx context.Context, req *v1.OptionsReq) (res *v1.OptionsRes, err error) {
	items, err := c.equipmentSvc.Options(ctx, equipmentsvc.OptionsInput{
		Keyword: req.Keyword,
	})
	if err != nil {
		return nil, err
	}
	list := make([]*v1.OptionItem, 0, len(items))
	for _, item := range items {
		list = append(list, &v1.OptionItem{
			Id:            item.Id,
			EquipmentCode: item.EquipmentCode,
			EquipmentName: item.EquipmentName,
			Status:        item.Status,
		})
	}
	return &v1.OptionsRes{List: list}, nil
}
