// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package view

import (
	"context"

	"lina-plugin-linapro-live-manage/backend/api/view/v1"
)

// IViewV1 is the public viewer watch-statistics API interface. The gf CLI is
// unavailable in the current environment, so this interface file is
// hand-maintained in the exact generated format; task records note this
// exception per the OpenSpec change add-live-h5-player and its successors.
type IViewV1 interface {
	Heartbeat(ctx context.Context, req *v1.HeartbeatReq) (res *v1.HeartbeatRes, err error)
}
