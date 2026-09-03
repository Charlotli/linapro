// play_impl.go implements the anonymous viewer play-info resolution: explicit
// tenant validation, point queries on tenant-prefixed indexes, and the
// play-priority selection (ongoing, replay, preview). The authoritative
// visibility boundary is the is_public flag combined with the live state
// machine and the room live status, as documented in the OpenSpec change
// add-live-h5-player.

package play

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"

	"lina-core/pkg/bizerr"
	"lina-core/pkg/plugin/capability/tenantcap"
	"lina-plugin-linapro-live-manage/backend/internal/dao"
	"lina-plugin-linapro-live-manage/backend/internal/model/entity"
	"lina-plugin-linapro-live-manage/backend/internal/service/liveroom"
)

// Play enumeration values mirror the plugin-wide dictionaries. They repeat the
// values owned by the sibling service packages because this package stays
// independent of admin-facing services.
const (
	// playLiveStateNotStarted matches the plugin_live_state value 0.
	playLiveStateNotStarted = 0
	// playLiveStateOngoing matches the plugin_live_state value 1.
	playLiveStateOngoing = 1
	// playLiveStateFinished matches the plugin_live_state value 2.
	playLiveStateFinished = 2
	// playVisibilityPublic matches the plugin_live_public value 1.
	playVisibilityPublic = 1
	// playPlatformTenantID is the platform tenant owning single-tenant data.
	playPlatformTenantID = 0
)

// tenantStatusActive mirrors the provider-owned tenant lifecycle status that
// allows normal tenant access. The host fallback and the tenant-core provider
// both store this value for active tenants.
const tenantStatusActive = "active"

// playStateOrder lists the candidate live states in play priority order: the
// ongoing live wins, the latest finished live serves as the replay, and the
// latest not-started live renders as the preview.
var playStateOrder = []int{playLiveStateOngoing, playLiveStateFinished, playLiveStateNotStarted}

// Get resolves one playable live for the requested room and tenant.
func (s *serviceImpl) Get(ctx context.Context, in GetInput) (*PlayInfo, error) {
	tenantID, err := s.resolveTenantID(ctx, in.TenantId)
	if err != nil {
		return nil, err
	}

	room, err := loadRoom(ctx, tenantID, in.RoomCode)
	if err != nil {
		return nil, err
	}
	if room == nil {
		return nil, bizerr.NewCode(CodePlayNotFound)
	}

	live, err := s.selectPlayableLive(ctx, tenantID, room)
	if err != nil {
		return nil, err
	}
	if live == nil {
		return nil, bizerr.NewCode(CodePlayNotFound)
	}
	return buildPlayInfo(room, live), nil
}

// resolveTenantID maps the explicit tenant parameter onto the query tenant. It
// forces the platform tenant when the tenant capability is disabled, requires
// the parameter when enabled, and validates existence plus lifecycle status.
func (s *serviceImpl) resolveTenantID(ctx context.Context, tenantID *int) (int, error) {
	if s.tenantSvc == nil || !s.tenantSvc.Available(ctx) {
		// Without the tenant capability every plugin row lives under the
		// platform tenant, so the parameter is ignored by design.
		return playPlatformTenantID, nil
	}
	if tenantID == nil {
		return 0, bizerr.NewCode(CodePlayTenantRequired)
	}
	if *tenantID < 0 {
		return 0, bizerr.NewCode(CodePlayTenantInvalid)
	}
	info, err := s.tenantSvc.Directory().Get(ctx, tenantcap.TenantID(*tenantID))
	if err != nil {
		return 0, bizerr.WrapCode(err, CodePlayTenantInvalid)
	}
	if info == nil || info.Status != tenantStatusActive {
		return 0, bizerr.NewCode(CodePlayTenantInvalid)
	}
	return *tenantID, nil
}

// loadRoom resolves one tenant-scoped live room by its external code. The
// (tenant_id, room_code) unique index turns this into a point query.
func loadRoom(ctx context.Context, tenantID int, roomCode string) (*entity.Room, error) {
	cols := dao.Room.Columns()
	var room *entity.Room
	err := dao.Room.Ctx(ctx).
		Fields(cols.Id, cols.RoomCode, cols.RoomName, cols.Status).
		Where(cols.TenantId, tenantID).
		Where(cols.RoomCode, strings.TrimSpace(roomCode)).
		Scan(&room)
	if err != nil {
		return nil, gerror.Wrap(err, "load live room for public play")
	}
	return room, nil
}

// selectPlayableLive walks the play priority order and returns the first
// presentable public live of the room. At most len(playStateOrder) point
// queries run; the bound is a fixed small constant by design.
func (s *serviceImpl) selectPlayableLive(ctx context.Context, tenantID int, room *entity.Room) (*entity.Live, error) {
	for _, state := range playStateOrder {
		if !presentable(state, room) {
			continue
		}
		live, err := probeLive(ctx, tenantID, room.Id, state)
		if err != nil {
			return nil, err
		}
		if live != nil {
			return live, nil
		}
	}
	return nil, nil
}

// presentable reports whether a live in the given state may be shown for the
// room. Ongoing lives additionally require the room to be in live status, so
// disabling a room immediately stops its ongoing stream.
func presentable(state int, room *entity.Room) bool {
	if state != playLiveStateOngoing {
		return true
	}
	return room != nil && room.Status == liveroom.RoomStatusLive
}

// playURLVisible reports whether the play URL may leave the service for the
// given state. Not-started lives keep their stream address private until the
// admin starts the live.
func playURLVisible(state int) bool {
	return state == playLiveStateOngoing || state == playLiveStateFinished
}

// probeLive loads the newest public live of one room in the requested state.
func probeLive(ctx context.Context, tenantID int, roomID int64, state int) (*entity.Live, error) {
	cols := dao.Live.Columns()
	var live *entity.Live
	err := dao.Live.Ctx(ctx).
		Fields(
			cols.Id, cols.RoomId, cols.Title, cols.CoverUrl, cols.LiveUrl,
			cols.LiveDate, cols.StartTime, cols.State, cols.Host, cols.SongName,
			cols.SongList, cols.LeadSinger, cols.Accompaniment, cols.SermonTitle,
			cols.Preacher, cols.PreacherIdentity, cols.ScriptureRef,
			cols.ScriptureContent, cols.Outline,
		).
		Where(cols.TenantId, tenantID).
		Where(cols.RoomId, roomID).
		Where(cols.State, state).
		Where(cols.IsPublic, playVisibilityPublic).
		OrderDesc(cols.Id).
		Limit(1).
		Scan(&live)
	if err != nil {
		return nil, gerror.Wrapf(err, "probe public live state=%d", state)
	}
	return live, nil
}

// buildPlayInfo projects the room and live rows into the viewer-facing
// payload. The play URL is withheld for not-started lives and the push URL is
// never exposed.
func buildPlayInfo(room *entity.Room, live *entity.Live) *PlayInfo {
	playURL := ""
	if playURLVisible(live.State) {
		playURL = live.LiveUrl
	}
	return &PlayInfo{
		RoomId:           room.Id,
		RoomCode:         room.RoomCode,
		RoomName:         room.RoomName,
		LiveId:           live.Id,
		Title:            live.Title,
		CoverUrl:         live.CoverUrl,
		LiveUrl:          playURL,
		LiveDate:         live.LiveDate,
		StartTime:        live.StartTime,
		State:            live.State,
		Host:             live.Host,
		SongName:         live.SongName,
		SongList:         live.SongList,
		LeadSinger:       live.LeadSinger,
		Accompaniment:    live.Accompaniment,
		SermonTitle:      live.SermonTitle,
		Preacher:         live.Preacher,
		PreacherIdentity: live.PreacherIdentity,
		ScriptureRef:     live.ScriptureRef,
		ScriptureContent: live.ScriptureContent,
		Outline:          live.Outline,
	}
}
