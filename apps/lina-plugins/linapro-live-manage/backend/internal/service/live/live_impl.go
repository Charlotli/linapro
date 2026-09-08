// live_impl.go implements tenant-scoped live-content CRUD, live-room
// availability validation, song-list JSON validation, and batched room-name
// assembly for the linapro-live-manage plugin. Tenant filters are applied at
// the database query stage and list assembly stays bounded to the current page.

package live

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"

	"lina-core/pkg/bizerr"
	"lina-core/pkg/plugin/capability/tenantcap"
	"lina-core/pkg/plugin/capability/tenantcap/tenantspi"
	"lina-core/pkg/plugin/capability/usercap"
	"lina-plugin-linapro-live-manage/backend/internal/dao"
	"lina-plugin-linapro-live-manage/backend/internal/model/do"
	"lina-plugin-linapro-live-manage/backend/internal/service/liveroom"
	"lina-plugin-linapro-live-manage/backend/internal/service/view"
)

const (
	// liveDateFormat is the date-only wire format of live_date.
	liveDateFormat = "2006-01-02"
)

// List queries live content with pagination and filters.
func (s *serviceImpl) List(ctx context.Context, in ListInput) (*ListOutput, error) {
	liveColumns := dao.Live.Columns()

	m := tenantspi.ApplyPluginTableFilter(ctx, s.pluginTableFilter(), dao.Live.Ctx(ctx), "")

	// Apply filters before counting so pagination stays consistent.
	if title := strings.TrimSpace(in.Title); title != "" {
		m = m.WhereLike(liveColumns.Title, "%"+title+"%")
	}
	if in.RoomId > 0 {
		m = m.Where(liveColumns.RoomId, in.RoomId)
	}
	if in.State != nil {
		m = m.Where(liveColumns.State, *in.State)
	}
	if in.IsPublic != nil {
		m = m.Where(liveColumns.IsPublic, *in.IsPublic)
	}
	if dateStart := strings.TrimSpace(in.LiveDateStart); dateStart != "" {
		m = m.WhereGTE(liveColumns.LiveDate, dateStart)
	}
	if dateEnd := strings.TrimSpace(in.LiveDateEnd); dateEnd != "" {
		m = m.WhereLTE(liveColumns.LiveDate, dateEnd+" 23:59:59")
	}

	total, err := m.Count()
	if err != nil {
		return nil, err
	}

	list := make([]*LiveEntity, 0)
	err = m.Page(in.PageNum, in.PageSize).
		OrderDesc(liveColumns.Id).
		Scan(&list)
	if err != nil {
		return nil, err
	}

	roomNames, err := s.resolveRoomNameMap(ctx, list)
	if err != nil {
		return nil, err
	}
	userNameMap, err := s.resolveCreatorNameMap(ctx, list)
	if err != nil {
		return nil, err
	}
	viewStats, err := s.resolveViewStatsMap(ctx, list)
	if err != nil {
		return nil, err
	}

	items := make([]*ListItem, 0, len(list))
	for _, item := range list {
		stats := viewStats[item.Id]
		onlineCount, totalViews := int64(0), int64(0)
		if stats != nil {
			onlineCount, totalViews = stats.OnlineCount, stats.TotalViews
		}
		items = append(items, &ListItem{
			LiveEntity:    item,
			RoomName:      roomNames[item.RoomId],
			CreatedByName: userNameMap[item.CreatedBy],
			OnlineCount:   onlineCount,
			TotalViews:    totalViews,
		})
	}

	return &ListOutput{
		List:  items,
		Total: total,
	}, nil
}

// resolveViewStatsMap assembles the watch counters for the current page with
// exactly one grouped query. Empty pages skip the query entirely, and lives
// without sessions simply fall back to zero counters in the caller.
func (s *serviceImpl) resolveViewStatsMap(ctx context.Context, list []*LiveEntity) (map[int64]*view.Stats, error) {
	if len(list) == 0 || s.viewSvc == nil {
		return map[int64]*view.Stats{}, nil
	}
	liveIDs := make([]int64, 0, len(list))
	for _, item := range list {
		liveIDs = append(liveIDs, item.Id)
	}
	return s.viewSvc.BatchGet(ctx, liveIDs)
}

// GetById retrieves one live content record by ID.
func (s *serviceImpl) GetById(ctx context.Context, id int64) (*ListItem, error) {
	liveColumns := dao.Live.Columns()

	var record *LiveEntity
	err := tenantspi.ApplyPluginTableFilter(ctx, s.pluginTableFilter(), dao.Live.Ctx(ctx), "").
		Where(liveColumns.Id, id).
		Scan(&record)
	if err != nil {
		return nil, err
	}
	if record == nil {
		return nil, bizerr.NewCode(CodeLiveNotFound)
	}

	item := &ListItem{LiveEntity: record}
	roomNames, err := s.resolveRoomNameMap(ctx, []*LiveEntity{record})
	if err != nil {
		return nil, err
	}
	item.RoomName = roomNames[record.RoomId]
	if record.CreatedBy > 0 {
		names, err := s.resolveCreatorNameMap(ctx, []*LiveEntity{record})
		if err != nil {
			return nil, err
		}
		if name := names[record.CreatedBy]; name != "" {
			item.CreatedByName = name
		}
	}
	return item, nil
}

// Create creates one live content record.
func (s *serviceImpl) Create(ctx context.Context, in CreateInput) (int64, error) {
	if !liveStateValid(in.State) {
		return 0, bizerr.NewCode(CodeLiveStateInvalid)
	}
	if !liveVisibilityValid(in.IsPublic) {
		return 0, bizerr.NewCode(CodeLivePublicInvalid)
	}
	songList, err := normalizeSongList(in.SongList)
	if err != nil {
		return 0, err
	}
	liveDate, err := parseLiveDate(in.LiveDate)
	if err != nil {
		return 0, err
	}
	if err = s.ensureRoomAvailable(ctx, in.RoomId); err != nil {
		return 0, err
	}

	var (
		bizCtx    = s.bizCtxSvc.Current(ctx)
		createdBy = int64(bizCtx.UserID)
		tenantID  = s.tenantFilterContext(ctx).TenantID
	)

	// GoFrame auto-fills created_at and updated_at. The replay switch
	// defaults to true when the request omits it, matching the column
	// default so existing behavior is unchanged.
	replayEnabled := in.ReplayEnabled == nil || *in.ReplayEnabled
	id, err := dao.Live.Ctx(ctx).Data(do.Live{
		TenantId:         tenantID,
		RoomId:           in.RoomId,
		Title:            in.Title,
		Content:          in.Content,
		LiveDate:         liveDate,
		PageUrl:          in.PageUrl,
		PushUrl:          in.PushUrl,
		LiveUrl:          in.LiveUrl,
		CoverUrl:         in.CoverUrl,
		SongName:         in.SongName,
		SongList:         songList,
		LeadSinger:       in.LeadSinger,
		Accompaniment:    in.Accompaniment,
		Host:             in.Host,
		SermonTitle:      in.SermonTitle,
		Preacher:         in.Preacher,
		PreacherIdentity: in.PreacherIdentity,
		ScriptureRef:     in.ScriptureRef,
		ScriptureContent: in.ScriptureContent,
		Outline:          in.Outline,
		DeviceInfo:       in.DeviceInfo,
		Reception:        in.Reception,
		State:            in.State,
		IsPublic:         in.IsPublic,
		StartTime:        millisToTime(in.StartTime),
		ReplayEnabled:    replayEnabled,
		CreatedBy:        createdBy,
		UpdatedBy:        createdBy,
	}).InsertAndGetId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

// Update updates live-content information.
func (s *serviceImpl) Update(ctx context.Context, in UpdateInput) error {
	if in.State != nil && !liveStateValid(*in.State) {
		return bizerr.NewCode(CodeLiveStateInvalid)
	}
	if in.IsPublic != nil && !liveVisibilityValid(*in.IsPublic) {
		return bizerr.NewCode(CodeLivePublicInvalid)
	}
	var normalizedSongList *string
	if in.SongList != nil {
		normalized, err := normalizeSongList(*in.SongList)
		if err != nil {
			return err
		}
		normalizedSongList = &normalized
	}
	liveColumns := dao.Live.Columns()

	// Load the tenant-scoped record so the immutable room binding can be
	// enforced before any field update is applied.
	existing := &LiveEntity{}
	err := tenantspi.ApplyPluginTableFilter(ctx, s.pluginTableFilter(), dao.Live.Ctx(ctx), "").
		Where(liveColumns.Id, in.Id).
		Scan(existing)
	if err != nil {
		return err
	}
	if existing.Id == 0 {
		return bizerr.NewCode(CodeLiveNotFound)
	}

	var liveDate *time.Time
	if in.LiveDate != nil {
		parsed, parseErr := parseLiveDate(*in.LiveDate)
		if parseErr != nil {
			return parseErr
		}
		liveDate = &parsed
	}
	if in.RoomId != nil {
		// The room binding is immutable after creation: re-binding would bypass
		// the ongoing-live occupancy check and desync live-room status linkage.
		if *in.RoomId != existing.RoomId {
			return bizerr.NewCode(CodeRoomImmutable)
		}
		if err = s.ensureRoomAvailable(ctx, *in.RoomId); err != nil {
			return err
		}
	}

	var (
		bizCtx    = s.bizCtxSvc.Current(ctx)
		updatedBy = int64(bizCtx.UserID)
		tenantID  = s.tenantFilterContext(ctx).TenantID
	)

	data := buildLiveUpdateData(in, liveDate, normalizedSongList, updatedBy)

	_, err = dao.Live.Ctx(ctx).
		OmitNilData().
		Where(tenantspi.TenantFilterColumn, tenantID).
		Where(liveColumns.Id, in.Id).
		Data(data).
		Update()
	return err
}

// buildLiveUpdateData projects the optional update fields onto the DO carrier:
// only fields present in the request are populated so omitted fields keep
// their stored values under OmitNilData.
func buildLiveUpdateData(
	in UpdateInput,
	liveDate *time.Time,
	normalizedSongList *string,
	updatedBy int64,
) do.Live {
	data := do.Live{UpdatedBy: updatedBy}
	if in.RoomId != nil {
		data.RoomId = *in.RoomId
	}
	if in.Title != nil {
		data.Title = *in.Title
	}
	if in.Content != nil {
		data.Content = *in.Content
	}
	if liveDate != nil {
		data.LiveDate = *liveDate
	}
	if in.PageUrl != nil {
		data.PageUrl = *in.PageUrl
	}
	if in.PushUrl != nil {
		data.PushUrl = *in.PushUrl
	}
	if in.LiveUrl != nil {
		data.LiveUrl = *in.LiveUrl
	}
	if in.CoverUrl != nil {
		data.CoverUrl = *in.CoverUrl
	}
	if in.SongName != nil {
		data.SongName = *in.SongName
	}
	if normalizedSongList != nil {
		data.SongList = *normalizedSongList
	}
	if in.LeadSinger != nil {
		data.LeadSinger = *in.LeadSinger
	}
	if in.Accompaniment != nil {
		data.Accompaniment = *in.Accompaniment
	}
	if in.Host != nil {
		data.Host = *in.Host
	}
	if in.SermonTitle != nil {
		data.SermonTitle = *in.SermonTitle
	}
	if in.Preacher != nil {
		data.Preacher = *in.Preacher
	}
	if in.PreacherIdentity != nil {
		data.PreacherIdentity = *in.PreacherIdentity
	}
	if in.ScriptureRef != nil {
		data.ScriptureRef = *in.ScriptureRef
	}
	if in.ScriptureContent != nil {
		data.ScriptureContent = *in.ScriptureContent
	}
	if in.Outline != nil {
		data.Outline = *in.Outline
	}
	if in.DeviceInfo != nil {
		data.DeviceInfo = *in.DeviceInfo
	}
	if in.Reception != nil {
		data.Reception = *in.Reception
	}
	if in.State != nil {
		data.State = *in.State
	}
	if in.IsPublic != nil {
		data.IsPublic = *in.IsPublic
	}
	if in.StartTime != nil {
		data.StartTime = millisToTime(in.StartTime)
	}
	if in.ReplayEnabled != nil {
		data.ReplayEnabled = *in.ReplayEnabled
	}
	return data
}

// Delete soft-deletes live content records by IDs.
func (s *serviceImpl) Delete(ctx context.Context, ids []int64) error {
	idList := normalizePositiveInt64IDs(ids)
	if len(idList) == 0 {
		return bizerr.NewCode(CodeLiveDeleteRequired)
	}

	liveColumns := dao.Live.Columns()
	_, err := tenantspi.ApplyPluginTableFilter(ctx, s.pluginTableFilter(), dao.Live.Ctx(ctx), "").
		WhereIn(liveColumns.Id, idList).
		Delete()
	return err
}

// Start transitions one not-started live to ongoing with room linkage.
func (s *serviceImpl) Start(ctx context.Context, id int64) error {
	return s.transitionLiveState(ctx, id, liveStateActionStart)
}

// Stop transitions one ongoing live to finished with room restoration.
func (s *serviceImpl) Stop(ctx context.Context, id int64) error {
	return s.transitionLiveState(ctx, id, liveStateActionStop)
}

// transitionLiveState validates and executes one explicit state action under
// the tenant filter, linking the live update and the live-room status change
// in one transaction so partial transitions never persist.
func (s *serviceImpl) transitionLiveState(ctx context.Context, id int64, action liveStateAction) error {
	record, err := s.loadTenantLive(ctx, id)
	if err != nil {
		return err
	}
	targetState, ok := liveStateTransition(record.State, action)
	if !ok {
		return bizerr.NewCode(CodeLiveStateTransition)
	}
	if action == liveStateActionStart {
		if err = s.ensureRoomAvailable(ctx, record.RoomId); err != nil {
			return err
		}
		if err = s.ensureRoomFreeForLive(ctx, record.RoomId, id); err != nil {
			return err
		}
	}

	var (
		bizCtx    = s.bizCtxSvc.Current(ctx)
		updatedBy = int64(bizCtx.UserID)
		tenantID  = s.tenantFilterContext(ctx).TenantID
	)

	return dao.Live.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		liveColumns := dao.Live.Columns()
		data := do.Live{State: targetState, UpdatedBy: updatedBy}
		if action == liveStateActionStart {
			now := time.Now()
			data.StartTime = &now
		}
		if _, err := dao.Live.Ctx(ctx).
			OmitNilData().
			Where(tenantspi.TenantFilterColumn, tenantID).
			Where(liveColumns.Id, id).
			Data(data).
			Update(); err != nil {
			return err
		}
		return s.syncRoomLiveStatus(ctx, tenantID, updatedBy, record.RoomId, action)
	})
}

// loadTenantLive loads one tenant-scoped live record for state actions and
// reports out-of-scope or absent records as not found.
func (s *serviceImpl) loadTenantLive(ctx context.Context, id int64) (*LiveEntity, error) {
	liveColumns := dao.Live.Columns()
	var record *LiveEntity
	err := tenantspi.ApplyPluginTableFilter(ctx, s.pluginTableFilter(), dao.Live.Ctx(ctx), "").
		Where(liveColumns.Id, id).
		Scan(&record)
	if err != nil {
		return nil, err
	}
	if record == nil {
		return nil, bizerr.NewCode(CodeLiveNotFound)
	}
	return record, nil
}

// ensureRoomFreeForLive verifies no other live in the same room is ongoing so
// one room hosts at most one ongoing live at a time.
func (s *serviceImpl) ensureRoomFreeForLive(ctx context.Context, roomID int64, excludeLiveID int64) error {
	liveColumns := dao.Live.Columns()
	count, err := tenantspi.ApplyPluginTableFilter(ctx, s.pluginTableFilter(), dao.Live.Ctx(ctx), "").
		Where(liveColumns.RoomId, roomID).
		Where(liveColumns.State, liveStateOngoing).
		WhereNot(liveColumns.Id, excludeLiveID).
		Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return bizerr.NewCode(CodeRoomBusy)
	}
	return nil
}

// syncRoomLiveStatus flips the associated live room inside the state-action
// transaction: starting marks the room live, while stopping restores an idle
// status only when the room is still marked live.
func (s *serviceImpl) syncRoomLiveStatus(
	ctx context.Context,
	tenantID int,
	updatedBy int64,
	roomID int64,
	action liveStateAction,
) error {
	roomColumns := dao.Room.Columns()
	targetStatus := liveroom.RoomStatusIdle
	if action == liveStateActionStart {
		targetStatus = liveroom.RoomStatusLive
	}
	m := tenantspi.ApplyPluginTableFilter(ctx, s.pluginTableFilter(), dao.Room.Ctx(ctx), "")
	if action == liveStateActionStop {
		// Only an occupied room needs restoration; other statuses stay untouched.
		m = m.Where(roomColumns.Status, liveroom.RoomStatusLive)
	}
	_, err := m.
		Where(roomColumns.Id, roomID).
		Data(do.Room{Status: targetStatus, UpdatedBy: updatedBy}).
		Update()
	return err
}

// ensureRoomAvailable verifies the live room exists in the current tenant and
// is not disabled. Out-of-scope rooms are invisible and reported as unavailable
// without leaking their existence.
func (s *serviceImpl) ensureRoomAvailable(ctx context.Context, roomID int64) error {
	roomColumns := dao.Room.Columns()
	var room *liveroom.RoomEntity
	err := tenantspi.ApplyPluginTableFilter(ctx, s.pluginTableFilter(), dao.Room.Ctx(ctx), "").
		Fields(roomColumns.Id, roomColumns.Status).
		Where(roomColumns.Id, roomID).
		Scan(&room)
	if err != nil {
		return err
	}
	if room == nil || room.Status == liveroom.RoomStatusDisabled {
		return bizerr.NewCode(CodeRoomUnavailable)
	}
	return nil
}

// resolveRoomNameMap resolves live-room display names for the current page with
// one batch query, keeping list assembly free of per-row lookups.
func (s *serviceImpl) resolveRoomNameMap(ctx context.Context, records []*LiveEntity) (map[int64]string, error) {
	names := make(map[int64]string)
	roomIDs := make([]int64, 0, len(records))
	seen := make(map[int64]struct{}, len(records))
	for _, record := range records {
		if record == nil || record.RoomId <= 0 {
			continue
		}
		if _, ok := seen[record.RoomId]; ok {
			continue
		}
		seen[record.RoomId] = struct{}{}
		roomIDs = append(roomIDs, record.RoomId)
	}
	if len(roomIDs) == 0 {
		return names, nil
	}

	roomColumns := dao.Room.Columns()
	rooms := make([]*liveroom.RoomEntity, 0, len(roomIDs))
	err := tenantspi.ApplyPluginTableFilter(ctx, s.pluginTableFilter(), dao.Room.Ctx(ctx), "").
		Fields(roomColumns.Id, roomColumns.RoomName).
		WhereIn(roomColumns.Id, roomIDs).
		Scan(&rooms)
	if err != nil {
		return nil, err
	}
	for _, room := range rooms {
		if room == nil {
			continue
		}
		names[room.Id] = room.RoomName
	}
	return names, nil
}

// normalizeSongList validates the song-list payload and returns the canonical
// trimmed JSON text. Empty input stays empty; otherwise it must be a valid JSON
// array.
func normalizeSongList(songList string) (string, error) {
	trimmed := strings.TrimSpace(songList)
	if trimmed == "" {
		return "", nil
	}
	if !json.Valid([]byte(trimmed)) {
		return "", bizerr.NewCode(CodeSongListInvalid)
	}
	if strings.HasPrefix(trimmed, "[") {
		return trimmed, nil
	}
	return "", bizerr.NewCode(CodeSongListInvalid)
}

// parseLiveDate converts the date-only wire format into a stored date value.
// Empty input defaults to today. Formatting must use time.Now with the Go
// layout: gtime.Format interprets PHP-style tokens and would emit the layout
// digits literally, storing the constant "2006-01-02" instead of today.
func parseLiveDate(liveDate string) (time.Time, error) {
	trimmed := strings.TrimSpace(liveDate)
	if trimmed == "" {
		trimmed = time.Now().Format(liveDateFormat)
	}
	parsed, err := time.ParseInLocation(liveDateFormat, trimmed, time.Local)
	if err != nil {
		return time.Time{}, bizerr.NewCode(CodeLiveDateInvalid)
	}
	return parsed, nil
}

// millisToTime converts Unix milliseconds to a time pointer for storage.
func millisToTime(millis *int64) *time.Time {
	if millis == nil || *millis <= 0 {
		return nil
	}
	value := time.UnixMilli(*millis)
	return &value
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
func (s *serviceImpl) resolveCreatorNameMap(ctx context.Context, records []*LiveEntity) (map[int64]string, error) {
	names := make(map[int64]string)
	userIDs := creatorDomainIDs(records)
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

// creatorDomainIDs converts plugin-owned creator storage values to user-domain
// IDs for host capability calls while de-duplicating the current page.
func creatorDomainIDs(records []*LiveEntity) []usercap.UserID {
	ids := make([]usercap.UserID, 0, len(records))
	seen := make(map[int64]struct{}, len(records))
	for _, record := range records {
		if record == nil || record.CreatedBy <= 0 {
			continue
		}
		if _, ok := seen[record.CreatedBy]; ok {
			continue
		}
		seen[record.CreatedBy] = struct{}{}
		ids = append(ids, usercap.UserID(strconv.FormatInt(record.CreatedBy, 10)))
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
