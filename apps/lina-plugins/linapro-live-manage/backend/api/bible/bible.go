// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package bible

import (
	"context"

	"lina-plugin-linapro-live-manage/backend/api/bible/v1"
)

// IBibleV1 is the anonymous viewer bible API interface. The gf CLI is
// unavailable in the current environment, so this interface file is
// hand-maintained in the exact generated format; task records note this
// exception, consistent with the sibling play/view interfaces.
type IBibleV1 interface {
	Books(ctx context.Context, req *v1.BooksReq) (res *v1.BooksRes, err error)
	Chapter(ctx context.Context, req *v1.ChapterReq) (res *v1.ChapterRes, err error)
}
