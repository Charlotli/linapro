// announcement_v1_admin.go implements the admin controller methods serving
// the tenant-scoped announcement CRUD endpoints consumed by the plugin
// management console.

package announcement

import (
	"context"

	v1 "lina-plugin-linapro-live-manage/backend/api/announcement/v1"
	announcementsvc "lina-plugin-linapro-live-manage/backend/internal/service/announcement"
)

// List returns the paged announcement list of one room.
func (c *ControllerV1) List(ctx context.Context, req *v1.ListReq) (res *v1.ListRes, err error) {
	items, total, err := c.announcementSvc.List(ctx, announcementsvc.ListInput{
		RoomId:   req.RoomId,
		PageNum:  req.PageNum,
		PageSize: req.PageSize,
	})
	if err != nil {
		return nil, err
	}
	list := make([]*v1.AnnouncementItem, 0, len(items))
	for _, item := range items {
		list = append(list, &v1.AnnouncementItem{
			Id:        item.Id,
			RoomId:    item.RoomId,
			Title:     item.Title,
			Content:   item.Content,
			Enabled:   item.Enabled,
			Sort:      item.Sort,
			CreatedAt: item.CreatedAt,
			UpdatedAt: item.UpdatedAt,
		})
	}
	return &v1.ListRes{List: list, Total: total}, nil
}

// Create inserts one announcement and returns the new ID.
func (c *ControllerV1) Create(ctx context.Context, req *v1.CreateReq) (res *v1.CreateRes, err error) {
	id, err := c.announcementSvc.Create(ctx, announcementsvc.CreateInput{
		RoomId:  req.RoomId,
		Title:   req.Title,
		Content: req.Content,
		Enabled: req.Enabled,
		Sort:    req.Sort,
	})
	if err != nil {
		return nil, err
	}
	return &v1.CreateRes{Id: id}, nil
}

// Update mutates one announcement; nil fields keep stored values.
func (c *ControllerV1) Update(ctx context.Context, req *v1.UpdateReq) (res *v1.UpdateRes, err error) {
	if err := c.announcementSvc.Update(ctx, announcementsvc.UpdateInput{
		Id:      req.Id,
		Title:   req.Title,
		Content: req.Content,
		Enabled: req.Enabled,
		Sort:    req.Sort,
	}); err != nil {
		return nil, err
	}
	return &v1.UpdateRes{}, nil
}

// Delete soft-deletes one announcement.
func (c *ControllerV1) Delete(ctx context.Context, req *v1.DeleteReq) (res *v1.DeleteRes, err error) {
	if err := c.announcementSvc.Delete(ctx, req.Id); err != nil {
		return nil, err
	}
	return &v1.DeleteRes{}, nil
}
