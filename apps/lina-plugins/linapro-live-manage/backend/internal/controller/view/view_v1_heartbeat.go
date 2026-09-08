// view_v1_heartbeat.go implements the controller method that serves the
// anonymous watch-heartbeat endpoint consumed by the plugin H5 viewer page.

package view

import (
	"context"

	v1 "lina-plugin-linapro-live-manage/backend/api/view/v1"
	viewsvc "lina-plugin-linapro-live-manage/backend/internal/service/view"
)

// Heartbeat records one viewer watch session and returns the aggregate
// counters.
func (c *ControllerV1) Heartbeat(ctx context.Context, req *v1.HeartbeatReq) (res *v1.HeartbeatRes, err error) {
	input := viewsvc.HeartbeatInput{
		RoomCode:   req.RoomCode,
		SessionKey: req.SessionKey,
	}
	if req.TenantId != nil {
		tenantID := *req.TenantId
		input.TenantId = &tenantID
	}
	stats, err := c.viewSvc.Heartbeat(ctx, input)
	if err != nil {
		return nil, err
	}
	if stats == nil {
		stats = &viewsvc.Stats{}
	}
	return &v1.HeartbeatRes{
		OnlineCount: stats.OnlineCount,
		TotalViews:  stats.TotalViews,
	}, nil
}
