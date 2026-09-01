// This file defines linapro-equipment-manage maintenance business error codes.

package maintenance

import (
	"github.com/gogf/gf/v2/errors/gcode"

	"lina-core/pkg/bizerr"
)

var (
	// CodeMaintenanceNotFound reports that the requested maintenance record
	// does not exist in the current tenant.
	CodeMaintenanceNotFound = bizerr.MustDefine(
		"EQUIPMENT_MAINTENANCE_NOT_FOUND",
		"Maintenance record does not exist",
		gcode.CodeNotFound,
	)
	// CodeMaintTypeInvalid reports an unknown maintenance type value.
	CodeMaintTypeInvalid = bizerr.MustDefine(
		"EQUIPMENT_MAINT_TYPE_INVALID",
		"Maintenance type is invalid",
		gcode.CodeInvalidParameter,
	)
	// CodeMaintenanceDateInvalid reports a malformed maintenance date value.
	CodeMaintenanceDateInvalid = bizerr.MustDefine(
		"EQUIPMENT_MAINT_DATE_INVALID",
		"Maintenance date must be in YYYY-MM-DD format",
		gcode.CodeInvalidParameter,
	)
	// CodeMaintenanceDeleteRequired reports that a delete operation received
	// no record IDs.
	CodeMaintenanceDeleteRequired = bizerr.MustDefine(
		"EQUIPMENT_MAINTENANCE_DELETE_REQUIRED",
		"Select at least one maintenance record to delete",
		gcode.CodeInvalidParameter,
	)
)
