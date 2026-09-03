// This file defines linapro-live-manage viewer play business error codes.

package play

import (
	"github.com/gogf/gf/v2/errors/gcode"

	"lina-core/pkg/bizerr"
)

var (
	// CodePlayNotFound reports that the live room does not exist or hosts no
	// publicly visible live in any presentable state. One shared code keeps
	// private and cross-tenant records indistinguishable from missing rooms.
	CodePlayNotFound = bizerr.MustDefine(
		"LIVE_MANAGE_PLAY_NOT_FOUND",
		"No publicly visible live for this room",
		gcode.CodeNotFound,
	)
	// CodePlayTenantRequired reports that the tenant capability is enabled and
	// the anonymous play query did not carry the required tenant ID.
	CodePlayTenantRequired = bizerr.MustDefine(
		"LIVE_MANAGE_PLAY_TENANT_REQUIRED",
		"Tenant ID is required for the public play query",
		gcode.CodeInvalidParameter,
	)
	// CodePlayTenantInvalid reports that the supplied tenant ID does not
	// reference an existing enabled tenant.
	CodePlayTenantInvalid = bizerr.MustDefine(
		"LIVE_MANAGE_PLAY_TENANT_INVALID",
		"Tenant does not exist or is disabled",
		gcode.CodeInvalidParameter,
	)
)
