// This file wires the linapro-equipment-manage maintenance controller and
// shared response mappers.

package maintenance

import (
	"time"

	"lina-core/pkg/apitime"
	maintenanceapi "lina-plugin-linapro-equipment-manage/backend/api/maintenance"
	v1 "lina-plugin-linapro-equipment-manage/backend/api/maintenance/v1"
	maintenancesvc "lina-plugin-linapro-equipment-manage/backend/internal/service/maintenance"
)

// dateLayout is the date-only wire format of date-only fields.
const dateLayout = "2006-01-02"

// ControllerV1 is the maintenance controller.
type ControllerV1 struct {
	maintenanceSvc maintenancesvc.Service // maintenance service
}

// NewV1 creates and returns a new maintenance controller instance.
func NewV1(maintenanceSvc maintenancesvc.Service) maintenanceapi.IMaintenanceV1 {
	return &ControllerV1{maintenanceSvc: maintenanceSvc}
}

// toAPIMaintenanceItem converts the service-layer list item into the API DTO
// projection returned through HTTP responses.
func toAPIMaintenanceItem(item *maintenancesvc.MaintenanceItem) v1.MaintenanceItem {
	if item == nil || item.MaintenanceEntity == nil {
		return v1.MaintenanceItem{}
	}
	return v1.MaintenanceItem{
		Id:            item.Id,
		EquipmentId:   item.EquipmentId,
		EquipmentName: item.EquipmentName,
		MaintType:     item.MaintType,
		MaintDate:     formatDate(item.MaintDate),
		Maintainer:    item.Maintainer,
		Cost:          item.Cost,
		Content:       item.Content,
		Result:        item.Result,
		Remark:        item.Remark,
		CreatedAt:     apitime.Milli(item.CreatedAt),
		UpdatedAt:     apitime.Milli(item.UpdatedAt),
	}
}

// formatDate renders the stored date in the date-only wire format.
func formatDate(value time.Time) string {
	if value.IsZero() {
		return ""
	}
	return value.Format(dateLayout)
}
