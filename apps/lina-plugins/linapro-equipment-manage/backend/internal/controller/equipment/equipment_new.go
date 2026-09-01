// This file wires the linapro-equipment-manage equipment controller and
// shared response mappers.

package equipment

import (
	"time"

	"lina-core/pkg/apitime"
	equipmentapi "lina-plugin-linapro-equipment-manage/backend/api/equipment"
	v1 "lina-plugin-linapro-equipment-manage/backend/api/equipment/v1"
	equipmentsvc "lina-plugin-linapro-equipment-manage/backend/internal/service/equipment"
)

// ControllerV1 is the equipment controller.
type ControllerV1 struct {
	equipmentSvc equipmentsvc.Service // equipment service
}

// NewV1 creates and returns a new equipment controller instance.
func NewV1(equipmentSvc equipmentsvc.Service) equipmentapi.IEquipmentV1 {
	return &ControllerV1{equipmentSvc: equipmentSvc}
}

// toAPIEquipmentItem converts the service-layer list item into the API DTO
// projection returned through HTTP responses.
func toAPIEquipmentItem(item *equipmentsvc.ListItem) v1.EquipmentItem {
	if item == nil || item.EquipmentEntity == nil {
		return v1.EquipmentItem{}
	}
	return v1.EquipmentItem{
		Id:            item.Id,
		EquipmentCode: item.EquipmentCode,
		EquipmentName: item.EquipmentName,
		EquipmentType: item.EquipmentType,
		BrandModel:    item.BrandModel,
		PurchaseDate:  formatDate(item.PurchaseDate),
		PurchasePrice: item.PurchasePrice,
		Location:      item.Location,
		Owner:         item.Owner,
		Status:        item.Status,
		Remark:        item.Remark,
		CreatedBy:     item.CreatedBy,
		CreatedByName: item.CreatedByName,
		CreatedAt:     apitime.Milli(item.CreatedAt),
		UpdatedAt:     apitime.Milli(item.UpdatedAt),
	}
}

// dateLayout is the date-only wire format of date-only fields.
const dateLayout = "2006-01-02"

// formatDate renders the stored date in the date-only wire format.
func formatDate(value time.Time) string {
	if value.IsZero() {
		return ""
	}
	return value.Format(dateLayout)
}
