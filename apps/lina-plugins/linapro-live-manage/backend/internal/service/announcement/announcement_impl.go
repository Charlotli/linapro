// announcement_impl.go implements the announcement service: tenant-scoped
// admin CRUD over the announcement table with room ownership checks, plus the
// anonymous viewer list that resolves the room through the shared play
// service and never leaks room existence. Time fields are projected as Unix
// milliseconds at the service boundary per the API time contract.

package announcement

import (
	"context"
	"strings"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"

	"lina-core/pkg/bizerr"
	"lina-core/pkg/plugin/capability/tenantcap/tenantspi"
	"lina-plugin-linapro-live-manage/backend/internal/dao"
	"lina-plugin-linapro-live-manage/backend/internal/model/do"
	"lina-plugin-linapro-live-manage/backend/internal/model/entity"
	"lina-plugin-linapro-live-manage/backend/internal/service/play"
)

// List returns the paged admin announcement list of one room.
func (s *serviceImpl) List(ctx context.Context, in ListInput) ([]*Item, int, error) {
	if in.RoomId <= 0 {
		return nil, 0, bizerrRoomNotFound()
	}
	if err := s.ensureRoomInTenant(ctx, in.RoomId); err != nil {
		return nil, 0, err
	}

	pageNum, pageSize := normalizePage(in.PageNum, in.PageSize)
	cols := dao.Announcement.Columns()
	base := tenantspi.ApplyPluginTableFilter(ctx, s.pluginTableFilter(), dao.Announcement.Ctx(ctx), "").
		Where(cols.RoomId, in.RoomId)

	total, err := base.Count()
	if err != nil {
		return nil, 0, gerror.Wrap(err, "count announcements")
	}
	if total == 0 {
		return []*Item{}, 0, nil
	}

	var rows []*entity.Announcement
	err = base.
		OrderAsc(cols.Sort).
		OrderAsc(cols.Id).
		Page(pageNum, pageSize).
		Scan(&rows)
	if err != nil {
		return nil, 0, gerror.Wrap(err, "list announcements")
	}
	items := make([]*Item, 0, len(rows))
	for _, row := range rows {
		items = append(items, projectItem(row))
	}
	return items, total, nil
}

// Create inserts one announcement for a room owned by the current tenant.
func (s *serviceImpl) Create(ctx context.Context, in CreateInput) (int64, error) {
	title := strings.TrimSpace(in.Title)
	if title == "" {
		return 0, bizerrTitleRequired()
	}
	if err := s.ensureRoomInTenant(ctx, in.RoomId); err != nil {
		return 0, err
	}
	enabled := announcementEnabled
	if in.Enabled != nil {
		enabled = *in.Enabled
	}
	bizCtx := s.bizCtxSvc.Current(ctx)
	result, err := tenantspi.ApplyPluginTableFilter(ctx, s.pluginTableFilter(), dao.Announcement.Ctx(ctx), "").
		Data(do.Announcement{
			TenantId:  bizCtx.TenantID,
			RoomId:    in.RoomId,
			Title:     title,
			Content:   in.Content,
			Enabled:   enabled,
			Sort:      in.Sort,
			CreatedBy: int64(bizCtx.UserID),
			UpdatedBy: int64(bizCtx.UserID),
		}).
		Insert()
	if err != nil {
		return 0, gerror.Wrap(err, "create announcement")
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, gerror.Wrap(err, "read created announcement id")
	}
	return id, nil
}

// Update mutates the mutable fields of one announcement owned by the current tenant.
func (s *serviceImpl) Update(ctx context.Context, in UpdateInput) error {
	if in.Id <= 0 {
		return bizerr.NewCode(CodeAnnouncementNotFound)
	}
	existing, err := s.loadOwned(ctx, in.Id)
	if err != nil {
		return err
	}
	update := do.Announcement{
		UpdatedBy: int64(s.bizCtxSvc.Current(ctx).UserID),
	}
	if in.Title != nil {
		title := strings.TrimSpace(*in.Title)
		if title == "" {
			return bizerrTitleRequired()
		}
		update.Title = title
	}
	if in.Content != nil {
		update.Content = *in.Content
	}
	if in.Enabled != nil {
		update.Enabled = *in.Enabled
	}
	if in.Sort != nil {
		update.Sort = *in.Sort
	}
	if _, err := tenantspi.ApplyPluginTableFilter(ctx, s.pluginTableFilter(), dao.Announcement.Ctx(ctx), "").
		Where(dao.Announcement.Columns().Id, existing.Id).
		Data(update).
		Update(); err != nil {
		return gerror.Wrap(err, "update announcement")
	}
	return nil
}

// Delete soft-deletes one announcement owned by the current tenant.
func (s *serviceImpl) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return bizerr.NewCode(CodeAnnouncementNotFound)
	}
	existing, err := s.loadOwned(ctx, id)
	if err != nil {
		return err
	}
	if _, err := tenantspi.ApplyPluginTableFilter(ctx, s.pluginTableFilter(), dao.Announcement.Ctx(ctx), "").
		Where(dao.Announcement.Columns().Id, existing.Id).
		Data(do.Announcement{
			DeletedAt: timePtr(gtime.Now().Time),
			UpdatedBy: int64(s.bizCtxSvc.Current(ctx).UserID),
		}).
		Update(); err != nil {
		return gerror.Wrap(err, "delete announcement")
	}
	return nil
}

// ForRoom returns the enabled announcements of one room for anonymous viewers.
func (s *serviceImpl) ForRoom(ctx context.Context, in ForRoomInput) ([]*ViewerItem, error) {
	room, err := s.playInfoSvc.ResolveRoom(ctx, play.ResolveRoomInput{
		TenantId: in.TenantId,
		RoomCode: in.RoomCode,
	})
	if err != nil {
		// Tenant errors keep the public play contract verbatim; any other
		// failure degrades to an empty list because the viewer panel must
		// stay silent on infrastructure problems.
		if isTenantError(err) {
			return nil, err
		}
		return []*ViewerItem{}, nil
	}
	if room == nil {
		return []*ViewerItem{}, nil
	}
	cols := dao.Announcement.Columns()
	var rows []*entity.Announcement
	err = tenantspi.ApplyPluginTableFilter(ctx, s.pluginTableFilter(), dao.Announcement.Ctx(ctx), "").
		Where(cols.RoomId, room.Id).
		Where(cols.Enabled, announcementEnabled).
		OrderAsc(cols.Sort).
		OrderAsc(cols.Id).
		Limit(ViewerLimit).
		Scan(&rows)
	if err != nil {
		return []*ViewerItem{}, nil
	}
	items := make([]*ViewerItem, 0, len(rows))
	for _, row := range rows {
		items = append(items, &ViewerItem{
			Id:        row.Id,
			Title:     row.Title,
			Content:   row.Content,
			UpdatedAt: millis(row.UpdatedAt),
		})
	}
	return items, nil
}

// loadOwned resolves one tenant-scoped announcement by primary key.
func (s *serviceImpl) loadOwned(ctx context.Context, id int64) (*entity.Announcement, error) {
	cols := dao.Announcement.Columns()
	var row *entity.Announcement
	err := tenantspi.ApplyPluginTableFilter(ctx, s.pluginTableFilter(), dao.Announcement.Ctx(ctx), "").
		Where(cols.Id, id).
		Scan(&row)
	if err != nil {
		return nil, gerror.Wrap(err, "load announcement")
	}
	if row == nil {
		return nil, bizerr.NewCode(CodeAnnouncementNotFound)
	}
	return row, nil
}

// ensureRoomInTenant verifies the room exists within the current tenant scope.
func (s *serviceImpl) ensureRoomInTenant(ctx context.Context, roomID int64) error {
	cols := dao.Room.Columns()
	var room *entity.Room
	err := tenantspi.ApplyPluginTableFilter(ctx, s.pluginTableFilter(), dao.Room.Ctx(ctx), "").
		Where(cols.Id, roomID).
		Scan(&room)
	if err != nil {
		return gerror.Wrap(err, "load live room for announcement")
	}
	if room == nil {
		return bizerr.NewCode(CodeAnnouncementRoomNotFound)
	}
	return nil
}

// pluginTableFilter returns the tenant table-filter slice from the injected tenant service.
func (s *serviceImpl) pluginTableFilter() tenantcapFilterService {
	if s == nil || s.tenantSvc == nil {
		return nil
	}
	return s.tenantSvc.Filter()
}

// normalizePage clamps the admin paging input to the documented defaults and cap.
func normalizePage(pageNum int, pageSize int) (int, int) {
	if pageNum <= 0 {
		pageNum = 1
	}
	if pageSize <= 0 {
		pageSize = PageDefaultSize
	}
	if pageSize > PageMaxSize {
		pageSize = PageMaxSize
	}
	return pageNum, pageSize
}

// projectItem maps one entity row onto the admin projection.
func projectItem(row *entity.Announcement) *Item {
	return &Item{
		Id:        row.Id,
		RoomId:    row.RoomId,
		Title:     row.Title,
		Content:   row.Content,
		Enabled:   row.Enabled,
		Sort:      row.Sort,
		CreatedAt: millis(row.CreatedAt),
		UpdatedAt: millis(row.UpdatedAt),
	}
}

// millis converts a nullable timestamp to Unix milliseconds; nil becomes 0.
func millis(value *time.Time) int64 {
	if value == nil {
		return 0
	}
	return value.UnixMilli()
}

// timePtr converts a time value into a heap pointer for DO writes.
func timePtr(value time.Time) *time.Time {
	return &value
}

// tenantcapFilterService aliases the tenant filter service type returned by
// the tenant capability, keeping the helper signatures short.
type tenantcapFilterService = tenantcapFilter
