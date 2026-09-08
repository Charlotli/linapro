// This file defines linapro-live-manage viewer watch-statistics business
// error codes.

package view

import (
	"github.com/gogf/gf/v2/errors/gcode"

	"lina-core/pkg/bizerr"
)

var (
	// CodeViewSessionKeyInvalid reports that the heartbeat session key is
	// empty or longer than the column bound.
	CodeViewSessionKeyInvalid = bizerr.MustDefine(
		"LIVE_MANAGE_VIEW_SESSION_KEY_INVALID",
		"Session key is invalid",
		gcode.CodeInvalidParameter,
	)
)
