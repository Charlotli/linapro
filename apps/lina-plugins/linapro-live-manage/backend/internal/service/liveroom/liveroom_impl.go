// liveroom_impl.go implements tenant-scoped live-room CRUD, bounded option
// candidates, and referenced-room deletion protection for the
// linapro-live-manage plugin. Tenant filters are applied at the database query
// stage so out-of-scope rows never leave the database.

package liveroom

import (
	"context"
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"

	"lina-core/pkg/bizerr"
	"lina-core/pkg/plugin/capability/tenantcap"
	"lina-core/pkg/plugin/capability/tenantcap/tenantspi"
	"lina-core/pkg/plugin/capability/usercap"
	"lina-plugin-linapro-live-manage/backend/internal/dao"
	"lina-plugin-linapro-live-manage/backend/internal/model/do"
)

// List queries live rooms with pagination and filters.
func (s *serviceImpl) List(ctx context.Context, in ListInput) (*ListOutput, error) {
	roomColumns := dao.Room.Columns()

	m := tenantspi.ApplyPluginTableFilter(ctx, s.pluginTableFilter(), dao.Room.Ctx(ctx), "")

	// Apply filters before counting so pagination stays consistent.
	if name := strings.TrimSpace(in.RoomName); name != "" {
		m = m.WhereLike(roomColumns.RoomName, "%"+name+"%")
	}
	if in.RoomType > 0 {
		m = m.Where(roomColumns.RoomType, in.RoomType)
	}
	if in.Status != nil {
		m = m.Where(roomColumns.Status, *in.Status)
	}

	total, err := m.Count()
	if err != nil {
		return nil, err
	}

	list := make([]*RoomEntity, 0)
	err = m.Page(in.PageNum, in.PageSize).
		OrderDesc(roomColumns.Id).
		Scan(&list)
	if err != nil {
		return nil, err
	}

	userNameMap, err := s.resolveCreatorNameMap(ctx, list)
	if err != nil {
		return nil, err
	}

	items := make([]*ListItem, 0, len(list))
	for _, room := range list {
		items = append(items, &ListItem{
			RoomEntity:    room,
			CreatedByName: userNameMap[room.CreatedBy],
		})
	}

	return &ListOutput{
		List:  items,
		Total: total,
	}, nil
}

// Options returns bounded live-room candidates for select controls.
func (s *serviceImpl) Options(ctx context.Context, in OptionsInput) ([]*OptionItem, error) {
	roomColumns := dao.Room.Columns()

	m := tenantspi.ApplyPluginTableFilter(ctx, s.pluginTableFilter(), dao.Room.Ctx(ctx), "")
	if condition := buildOptionsCondition(m.Builder(), in); condition != nil {
		m = m.Where(condition)
	}
	limit := clampOptionsLimit(in.Limit)

	list := make([]*RoomEntity, 0, limit)
	err := m.Fields(roomColumns.Id, roomColumns.RoomCode, roomColumns.RoomName, roomColumns.Status).
		Limit(limit).
		OrderDesc(roomColumns.Id).
		Scan(&list)
	if err != nil {
		return nil, err
	}

	items := make([]*OptionItem, 0, len(list))
	for _, room := range list {
		items = append(items, &OptionItem{
			Id:       room.Id,
			RoomCode: room.RoomCode,
			RoomName: room.RoomName,
			Status:   room.Status,
		})
	}
	return items, nil
}

// buildOptionsCondition composes the candidate filter condition on the given
// builder seed, or returns nil when no filter applies. Keyword alternatives are
// grouped inside one sub-builder so the OR clauses stay parenthesized within
// this condition instead of escaping the tenant predicate applied on the model.
func buildOptionsCondition(b *gdb.WhereBuilder, in OptionsInput) *gdb.WhereBuilder {
	roomColumns := dao.Room.Columns()
	keyword := strings.TrimSpace(in.Keyword)
	if keyword == "" && in.IncludeDisabled {
		return nil
	}
	if keyword != "" {
		keywordFilter := b.WhereLike(roomColumns.RoomName, "%"+keyword+"%").
			WhereOrLike(roomColumns.RoomCode, "%"+keyword+"%")
		b = b.Where(keywordFilter)
	}
	if !in.IncludeDisabled {
		b = b.WhereNotIn(roomColumns.Status, []int{RoomStatusDisabled})
	}
	return b
}

// clampOptionsLimit resolves the effective candidate cap from the requested
// limit: non-positive values fall back to the default cap and oversized values
// are clamped to the hard maximum.
func clampOptionsLimit(limit int) int {
	if limit <= 0 {
		return OptionsDefaultLimit
	}
	if limit > OptionsLimit {
		return OptionsLimit
	}
	return limit
}

// GetById retrieves one live room by ID.
func (s *serviceImpl) GetById(ctx context.Context, id int64) (*ListItem, error) {
	roomColumns := dao.Room.Columns()

	var room *RoomEntity
	err := tenantspi.ApplyPluginTableFilter(ctx, s.pluginTableFilter(), dao.Room.Ctx(ctx), "").
		Where(roomColumns.Id, id).
		Scan(&room)
	if err != nil {
		return nil, err
	}
	if room == nil {
		return nil, bizerr.NewCode(CodeRoomNotFound)
	}

	item := &ListItem{RoomEntity: room}
	if room.CreatedBy > 0 {
		names, err := s.resolveCreatorNameMap(ctx, []*RoomEntity{room})
		if err != nil {
			return nil, err
		}
		if name := names[room.CreatedBy]; name != "" {
			item.CreatedByName = name
		}
	}
	return item, nil
}

// Create creates a new live room.
func (s *serviceImpl) Create(ctx context.Context, in CreateInput) (int64, error) {
	if !RoomType(in.RoomType) {
		return 0, bizerr.NewCode(CodeRoomTypeInvalid)
	}
	if !RoomStatus(in.Status) {
		return 0, bizerr.NewCode(CodeRoomStatusInvalid)
	}
	roomColumns := dao.Room.Columns()

	// Pre-check the tenant-unique room code so users see a business error
	// instead of a raw database constraint failure.
	existing, err := tenantspi.ApplyPluginTableFilter(ctx, s.pluginTableFilter(), dao.Room.Ctx(ctx), "").
		Where(roomColumns.RoomCode, strings.TrimSpace(in.RoomCode)).
		Count()
	if err != nil {
		return 0, err
	}
	if existing > 0 {
		return 0, bizerr.NewCode(CodeRoomCodeExists)
	}

	var (
		bizCtx    = s.bizCtxSvc.Current(ctx)
		createdBy = int64(bizCtx.UserID)
		tenantID  = s.tenantFilterContext(ctx).TenantID
	)

	// GoFrame auto-fills created_at and updated_at.
	id, err := dao.Room.Ctx(ctx).Data(do.Room{
		TenantId:    tenantID,
		RoomCode:    strings.TrimSpace(in.RoomCode),
		RoomName:    in.RoomName,
		RoomType:    in.RoomType,
		Status:      in.Status,
		Description: in.Description,
		CreatedBy:   createdBy,
		UpdatedBy:   createdBy,
	}).InsertAndGetId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

// Update updates live-room information.
func (s *serviceImpl) Update(ctx context.Context, in UpdateInput) error {
	if in.RoomType != nil && !RoomType(*in.RoomType) {
		return bizerr.NewCode(CodeRoomTypeInvalid)
	}
	if in.Status != nil && !RoomStatus(*in.Status) {
		return bizerr.NewCode(CodeRoomStatusInvalid)
	}
	roomColumns := dao.Room.Columns()

	// Ensure the room exists in the current tenant before updating.
	existing, err := tenantspi.ApplyPluginTableFilter(ctx, s.pluginTableFilter(), dao.Room.Ctx(ctx), "").
		Where(roomColumns.Id, in.Id).
		Count()
	if err != nil {
		return err
	}
	if existing == 0 {
		return bizerr.NewCode(CodeRoomNotFound)
	}

	if in.RoomCode != nil {
		code := strings.TrimSpace(*in.RoomCode)
		conflict, err := tenantspi.ApplyPluginTableFilter(ctx, s.pluginTableFilter(), dao.Room.Ctx(ctx), "").
			Where(roomColumns.RoomCode, code).
			WhereNot(roomColumns.Id, in.Id).
			Count()
		if err != nil {
			return err
		}
		if conflict > 0 {
			return bizerr.NewCode(CodeRoomCodeExists)
		}
	}

	var (
		bizCtx    = s.bizCtxSvc.Current(ctx)
		updatedBy = int64(bizCtx.UserID)
		tenantID  = s.tenantFilterContext(ctx).TenantID
	)

	data := do.Room{UpdatedBy: updatedBy}
	if in.RoomCode != nil {
		data.RoomCode = strings.TrimSpace(*in.RoomCode)
	}
	if in.RoomName != nil {
		data.RoomName = *in.RoomName
	}
	if in.RoomType != nil {
		data.RoomType = *in.RoomType
	}
	if in.Status != nil {
		data.Status = *in.Status
	}
	if in.Description != nil {
		data.Description = *in.Description
	}

	_, err = dao.Room.Ctx(ctx).
		OmitNilData().
		Where(tenantspi.TenantFilterColumn, tenantID).
		Where(roomColumns.Id, in.Id).
		Data(data).
		Update()
	return err
}

// Delete soft-deletes live rooms by IDs with referenced-room protection.
func (s *serviceImpl) Delete(ctx context.Context, ids []int64) error {
	idList := normalizePositiveInt64IDs(ids)
	if len(idList) == 0 {
		return bizerr.NewCode(CodeRoomDeleteRequired)
	}

	// Reject the whole batch when any room is still referenced by live content
	// so callers never observe a partial delete.
	referenced, err := s.firstReferencedRoom(ctx, idList)
	if err != nil {
		return err
	}
	if referenced > 0 {
		return bizerr.NewCode(CodeRoomReferenced)
	}

	roomColumns := dao.Room.Columns()
	_, err = tenantspi.ApplyPluginTableFilter(ctx, s.pluginTableFilter(), dao.Room.Ctx(ctx), "").
		WhereIn(roomColumns.Id, idList).
		Delete()
	return err
}

// firstReferencedRoom returns the first room ID in the list that is still
// referenced by live content rows in the current tenant, or 0 when none. One
// bounded query covers the whole batch.
func (s *serviceImpl) firstReferencedRoom(ctx context.Context, roomIDs []int64) (int64, error) {
	liveColumns := dao.Live.Columns()
	var referenced struct {
		RoomId int64 `json:"roomId"`
	}
	err := tenantspi.ApplyPluginTableFilter(ctx, s.pluginTableFilter(), dao.Live.Ctx(ctx), "").
		Fields(liveColumns.RoomId).
		WhereIn(liveColumns.RoomId, roomIDs).
		Limit(1).
		Scan(&referenced)
	if err != nil {
		return 0, err
	}
	return referenced.RoomId, nil
}

// normalizePositiveInt64IDs drops non-positive identifiers from a batch ID list.
func normalizePositiveInt64IDs(ids []int64) []int64 {
	if len(ids) == 0 {
		return nil
	}
	result := make([]int64, 0, len(ids))
	for _, id := range ids {
		if id > 0 {
			result = append(result, id)
		}
	}
	return result
}

// pluginTableFilter returns the tenant table-filter slice from the injected tenant service.
func (s *serviceImpl) pluginTableFilter() tenantcap.FilterService {
	if s == nil || s.tenantSvc == nil {
		return nil
	}
	return s.tenantSvc.Filter()
}

// tenantFilterContext returns current tenant metadata for write ownership fields.
func (s *serviceImpl) tenantFilterContext(ctx context.Context) tenantcap.TenantFilterContext {
	if filter := s.pluginTableFilter(); filter != nil {
		return filter.Context(ctx)
	}
	return tenantcap.TenantFilterContext{}
}

// resolveCreatorNameMap resolves current-page creator display names through one
// user-domain batch call and leaves invisible or missing users blank.
func (s *serviceImpl) resolveCreatorNameMap(ctx context.Context, rooms []*RoomEntity) (map[int64]string, error) {
	names := make(map[int64]string)
	userIDs := creatorDomainIDs(rooms)
	if len(userIDs) == 0 {
		return names, nil
	}
	if s.userSvc == nil {
		return nil, gerror.New("linapro-live-manage requires host user capability")
	}
	result, err := s.userSvc.BatchGet(ctx, userIDs)
	if err != nil || result == nil {
		return names, err
	}
	for id, projection := range result.Items {
		storageID, ok := userDomainIDStorageID(id)
		if !ok {
			continue
		}
		names[storageID] = userInfoDisplayName(projection)
	}
	return names, nil
}

// creatorDomainIDs converts plugin-owned room creator storage values to
// user-domain IDs for host capability calls while de-duplicating the current page.
func creatorDomainIDs(rooms []*RoomEntity) []usercap.UserID {
	ids := make([]usercap.UserID, 0, len(rooms))
	seen := make(map[int64]struct{}, len(rooms))
	for _, room := range rooms {
		if room == nil || room.CreatedBy <= 0 {
			continue
		}
		if _, ok := seen[room.CreatedBy]; ok {
			continue
		}
		seen[room.CreatedBy] = struct{}{}
		ids = append(ids, usercap.UserID(strconv.FormatInt(room.CreatedBy, 10)))
	}
	return ids
}

// userDomainIDStorageID parses the current host user-domain ID encoding used by
// existing plugin-owned creator columns.
func userDomainIDStorageID(id usercap.UserID) (int64, bool) {
	storageID, err := strconv.ParseInt(strings.TrimSpace(string(id)), 10, 64)
	return storageID, err == nil && storageID > 0
}

// userInfoDisplayName chooses the stable creator display field from the
// current user-domain projection.
func userInfoDisplayName(user *usercap.UserInfo) string {
	if user == nil {
		return ""
	}
	if user.Username != "" {
		return user.Username
	}
	if user.Nickname != "" {
		return user.Nickname
	}
	if user.Label != "" {
		return user.Label
	}
	return string(user.ID)
}
