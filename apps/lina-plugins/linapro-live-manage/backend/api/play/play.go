// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package play

import (
	"context"

	"lina-plugin-linapro-live-manage/backend/api/play/v1"
)

// IPlayV1 is the public viewer play API interface. The gf CLI is unavailable in
// the current environment, so this interface file is hand-maintained in the
// exact generated format; task records note this exception per the OpenSpec
// change add-live-h5-player.
type IPlayV1 interface {
	Get(ctx context.Context, req *v1.GetReq) (res *v1.GetRes, err error)
}
