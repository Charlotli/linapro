// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package subscribe

import (
	"context"

	"lina-plugin-linapro-live-manage/backend/api/subscribe/v1"
)

// ISubscribeV1 is the public viewer calendar subscription API interface. The
// gf CLI is unavailable in the current environment, so this interface file is
// hand-maintained in the exact generated format; task records note this
// exception per the OpenSpec change add-live-h5-player and its successors.
type ISubscribeV1 interface {
	Subscribe(ctx context.Context, req *v1.SubscribeReq) (res *v1.SubscribeRes, err error)
}
