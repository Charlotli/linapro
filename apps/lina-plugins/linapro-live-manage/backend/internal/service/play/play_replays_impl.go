// play_impl_replays.go implements the public paginated replay-library query.
// It shares the play service visibility boundary: tenant-scoped, public
// records only, finished state, and the per-live replay switch enabled.

package play

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"

	"lina-plugin-linapro-live-manage/backend/internal/dao"
	"lina-plugin-linapro-live-manage/backend/internal/model/entity"
)

// Replay pagination bounds.
const (
	// replaysDefaultPageSize is the page size used when the request omits it.
	replaysDefaultPageSize = 10
	// replaysMaxPageSize clamps oversized page sizes to keep queries bounded.
	replaysMaxPageSize = 50
)

// Replays returns one page of the room's public replay library plus the
// total count. Unknown rooms and empty libraries both yield an empty result
// with a zero total, keeping room existence indistinguishable.
func (s *serviceImpl) Replays(ctx context.Context, in ReplaysInput) ([]ReplayItem, int, error) {
	tenantID, err := s.resolveTenantID(ctx, in.TenantId)
	if err != nil {
		return nil, 0, err
	}
	room, err := loadRoom(ctx, tenantID, in.RoomCode)
	if err != nil {
		return nil, 0, err
	}
	if room == nil {
		return nil, 0, nil
	}

	cols := dao.Live.Columns()
	base := dao.Live.Ctx(ctx).
		Where(cols.TenantId, tenantID).
		Where(cols.RoomId, room.Id).
		Where(cols.State, playLiveStateFinished).
		Where(cols.IsPublic, playVisibilityPublic).
		Where(cols.ReplayEnabled, true)

	total, err := base.Count()
	if err != nil {
		return nil, 0, gerror.Wrap(err, "count public replays")
	}
	if total == 0 {
		return nil, 0, nil
	}

	page, pageSize := normalizeReplayPage(in.Page, in.PageSize)
	var lives []*entity.Live
	err = base.
		Fields(
			cols.Id, cols.Title, cols.CoverUrl, cols.LiveUrl,
			cols.LiveDate, cols.StartTime,
		).
		OrderDesc(cols.LiveDate).
		OrderDesc(cols.Id).
		Page(page, pageSize).
		Scan(&lives)
	if err != nil {
		return nil, 0, gerror.Wrap(err, "list public replays")
	}

	items := make([]ReplayItem, 0, len(lives))
	for _, live := range lives {
		items = append(items, ReplayItem{
			LiveId:    live.Id,
			Title:     live.Title,
			CoverUrl:  live.CoverUrl,
			LiveUrl:   live.LiveUrl,
			LiveDate:  live.LiveDate,
			StartTime: live.StartTime,
		})
	}
	return items, total, nil
}

// normalizeReplayPage applies the documented pagination defaults and the
// hard page-size cap.
func normalizeReplayPage(page int, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = replaysDefaultPageSize
	}
	if pageSize > replaysMaxPageSize {
		pageSize = replaysMaxPageSize
	}
	return page, pageSize
}
