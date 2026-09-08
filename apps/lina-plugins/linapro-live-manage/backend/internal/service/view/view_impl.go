// view_impl.go implements the watch-session heartbeat upsert and the batched
// per-live counter aggregation. The session table is a statistics ledger:
// rows carry no soft delete and are keyed by the unique session key so
// repeated heartbeats from one browser session update instead of inserting.

package view

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"lina-core/pkg/bizerr"
	"lina-plugin-linapro-live-manage/backend/internal/dao"
	"lina-plugin-linapro-live-manage/backend/internal/model/do"
	"lina-plugin-linapro-live-manage/backend/internal/service/play"
)

// maxSessionKeyLength mirrors the session_key column width.
const maxSessionKeyLength = 64

// Heartbeat records one watch session for the currently playable public live
// and returns the refreshed aggregate counters.
func (s *serviceImpl) Heartbeat(ctx context.Context, in HeartbeatInput) (*Stats, error) {
	sessionKey := strings.TrimSpace(in.SessionKey)
	if sessionKey == "" || len(sessionKey) > maxSessionKeyLength {
		return nil, bizerr.NewCode(CodeViewSessionKeyInvalid)
	}

	// Tenant resolution intentionally reuses the play service so the heartbeat
	// tenant contract can never drift from the play contract: tenant errors
	// propagate verbatim while unknown rooms and absent content surface as
	// zero counters, keeping invalid rooms indistinguishable from no
	// watchable live.
	playInfo, err := s.playSvc.Get(ctx, play.GetInput{
		TenantId: in.TenantId,
		RoomCode: in.RoomCode,
	})
	if err != nil {
		if bizerr.Is(err, play.CodePlayTenantRequired) || bizerr.Is(err, play.CodePlayTenantInvalid) {
			return nil, err
		}
		return &Stats{}, nil
	}
	if playInfo == nil || playInfo.LiveId <= 0 {
		return &Stats{}, nil
	}

	if err := upsertSession(ctx, playInfo.TenantId, playInfo.LiveId, sessionKey); err != nil {
		return nil, err
	}
	return s.statsForLive(ctx, playInfo.LiveId)
}

// upsertSession inserts the session row or refreshes the existing row when
// the unique session key already exists. Refreshing the live binding keeps a
// returning viewer attached to the currently playable live, and refreshing
// updated_at keeps repeated heartbeats inside the online window. OnDuplicate
// map values are rendered as identifiers in GoFrame, so literal values must
// go through gdb.Raw.
func upsertSession(ctx context.Context, tenantID int, liveID int64, sessionKey string) error {
	cols := dao.ViewSession.Columns()
	_, err := dao.ViewSession.Ctx(ctx).
		Data(do.ViewSession{
			TenantId:   tenantID,
			LiveId:     liveID,
			SessionKey: sessionKey,
		}).
		OnConflict(cols.SessionKey).
		OnDuplicate(g.Map{
			cols.LiveId:    gdb.Raw(strconv.FormatInt(liveID, 10)),
			cols.UpdatedAt: gdb.Raw("CURRENT_TIMESTAMP"),
		}).
		Save()
	if err != nil {
		return gerror.Wrap(err, "upsert watch session")
	}
	return nil
}

// statsForLive aggregates both counters for one live in a single grouped
// statement using the (live_id, updated_at) index.
func (s *serviceImpl) statsForLive(ctx context.Context, liveID int64) (*Stats, error) {
	statsMap, err := s.BatchGet(ctx, []int64{liveID})
	if err != nil {
		return nil, err
	}
	if stats, ok := statsMap[liveID]; ok {
		return stats, nil
	}
	return &Stats{}, nil
}

// BatchGet aggregates the online and total counters for the live ID set in
// one grouped query. The session rows are tenant-scoped at write time and
// live IDs are globally unique primary keys, so the aggregation needs no
// additional tenant filter. The online counter uses a cross-database CASE
// expression instead of a dialect-specific FILTER clause; the 60-second
// threshold is computed in Go and inlined as a standard single-quoted
// datetime literal because Fields expressions carry no parameter binding.
func (s *serviceImpl) BatchGet(ctx context.Context, liveIDs []int64) (map[int64]*Stats, error) {
	if len(liveIDs) == 0 {
		return map[int64]*Stats{}, nil
	}
	cols := dao.ViewSession.Columns()
	threshold := gtime.New(time.Now().Add(-onlineWindow)).Format("Y-m-d H:i:s")
	onlineCountExpr := fmt.Sprintf(
		"SUM(CASE WHEN %s > '%s' THEN 1 ELSE 0 END) AS online_count",
		cols.UpdatedAt, threshold,
	)
	type counterRow struct {
		LiveId      int64 `json:"liveId"`
		OnlineCount int64 `json:"onlineCount"`
		TotalViews  int64 `json:"totalViews"`
	}
	var rows []counterRow
	err := dao.ViewSession.Ctx(ctx).
		Fields(
			cols.LiveId,
			onlineCountExpr,
			"COUNT(*) AS total_views",
		).
		WhereIn(cols.LiveId, liveIDs).
		Group(cols.LiveId).
		Scan(&rows)
	if err != nil {
		return nil, gerror.Wrap(err, "aggregate watch session counters")
	}
	result := make(map[int64]*Stats, len(rows))
	for _, row := range rows {
		result[row.LiveId] = &Stats{
			OnlineCount: row.OnlineCount,
			TotalViews:  row.TotalViews,
		}
	}
	return result, nil
}
