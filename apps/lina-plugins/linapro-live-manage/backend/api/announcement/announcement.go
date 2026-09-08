// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package announcement

import (
	"context"

	"lina-plugin-linapro-live-manage/backend/api/announcement/v1"
)

// IAnnouncementV1 is the tenant-scoped admin announcement API interface. The
// gf CLI is unavailable in the current environment, so this interface file is
// hand-maintained in the exact generated format; task records note this
// exception, consistent with the sibling play/view interfaces.
type IAnnouncementV1 interface {
	List(ctx context.Context, req *v1.ListReq) (res *v1.ListRes, err error)
	Create(ctx context.Context, req *v1.CreateReq) (res *v1.CreateRes, err error)
	Update(ctx context.Context, req *v1.UpdateReq) (res *v1.UpdateRes, err error)
	Delete(ctx context.Context, req *v1.DeleteReq) (res *v1.DeleteRes, err error)
}

// IAnnouncementViewerV1 is the anonymous viewer announcement API interface.
// It binds only inside the anonymous viewer route group and therefore carries
// no permission-gated admin surface.
type IAnnouncementViewerV1 interface {
	ViewerList(ctx context.Context, req *v1.ViewerListReq) (res *v1.ViewerListRes, err error)
}
