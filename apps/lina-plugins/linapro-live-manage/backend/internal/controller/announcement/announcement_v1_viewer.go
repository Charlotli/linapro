// announcement_v1_viewer.go implements the controller method serving the
// anonymous viewer announcement endpoint consumed by the plugin H5 viewer
// page.

package announcement

import (
	"context"

	v1 "lina-plugin-linapro-live-manage/backend/api/announcement/v1"
	announcementsvc "lina-plugin-linapro-live-manage/backend/internal/service/announcement"
)

// ViewerList returns the enabled announcements of one room for the H5 viewer.
func (c *ControllerViewerV1) ViewerList(ctx context.Context, req *v1.ViewerListReq) (res *v1.ViewerListRes, err error) {
	input := announcementsvc.ForRoomInput{RoomCode: req.RoomCode}
	if req.TenantId != nil {
		tenantID := *req.TenantId
		input.TenantId = &tenantID
	}
	items, err := c.announcementSvc.ForRoom(ctx, input)
	if err != nil {
		return nil, err
	}
	list := make([]*v1.ViewerAnnouncementItem, 0, len(items))
	for _, item := range items {
		list = append(list, &v1.ViewerAnnouncementItem{
			Id:        item.Id,
			Title:     item.Title,
			Content:   item.Content,
			UpdatedAt: item.UpdatedAt,
		})
	}
	return &v1.ViewerListRes{List: list}, nil
}
