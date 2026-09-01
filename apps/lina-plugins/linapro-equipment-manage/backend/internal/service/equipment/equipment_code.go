// This file defines linapro-equipment-manage equipment business error codes.

package equipment

import (
	"github.com/gogf/gf/v2/errors/gcode"

	"lina-core/pkg/bizerr"
)

var (
	// CodeEquipmentNotFound reports that the requested equipment does not
	// exist in the current tenant.
	CodeEquipmentNotFound = bizerr.MustDefine(
		"EQUIPMENT_NOT_FOUND",
		"Equipment does not exist",
		gcode.CodeNotFound,
	)
	// CodeEquipmentCodeExists reports that the equipment code is already used
	// in the current tenant.
	CodeEquipmentCodeExists = bizerr.MustDefine(
		"EQUIPMENT_CODE_EXISTS",
		"Equipment code already exists",
		gcode.CodeInvalidParameter,
	)
	// CodeEquipmentTypeInvalid reports an unknown equipment type value.
	CodeEquipmentTypeInvalid = bizerr.MustDefine(
		"EQUIPMENT_TYPE_INVALID",
		"Equipment type is invalid",
		gcode.CodeInvalidParameter,
	)
	// CodeEquipmentStatusInvalid reports an unknown equipment status value.
	CodeEquipmentStatusInvalid = bizerr.MustDefine(
		"EQUIPMENT_STATUS_INVALID",
		"Equipment status is invalid",
		gcode.CodeInvalidParameter,
	)
	// CodeEquipmentDateInvalid reports a malformed date value.
	CodeEquipmentDateInvalid = bizerr.MustDefine(
		"EQUIPMENT_DATE_INVALID",
		"Date must be in YYYY-MM-DD format",
		gcode.CodeInvalidParameter,
	)
	// CodeEquipmentReferenced reports that equipment is still referenced by
	// maintenance records and cannot be deleted.
	CodeEquipmentReferenced = bizerr.MustDefine(
		"EQUIPMENT_REFERENCED",
		"Equipment is still referenced by maintenance records",
		gcode.CodeInvalidOperation,
	)
	// CodeEquipmentDeleteRequired reports that a delete operation received no
	// equipment IDs.
	CodeEquipmentDeleteRequired = bizerr.MustDefine(
		"EQUIPMENT_DELETE_REQUIRED",
		"Select at least one equipment record to delete",
		gcode.CodeInvalidParameter,
	)
)
