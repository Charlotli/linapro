// subscribe_v1_get.go implements the controller method that serves the
// anonymous public ICS calendar stream consumed by viewer calendar clients.
// The response bypasses the JSON envelope by design: the handler writes the
// raw text/calendar document directly, as documented in the OpenSpec change
// add-live-calendar-subscription.

package subscribe

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	v1 "lina-plugin-linapro-live-manage/backend/api/subscribe/v1"
	calendarsvc "lina-plugin-linapro-live-manage/backend/internal/service/calendar"
)

// calendarContentMediaType is the RFC 5545 response media type.
const calendarContentMediaType = "text/calendar; charset=utf-8"

// viewerH5Path is the stable plugin-public viewer page path embedded into
// event descriptions. The /x/<plugin-id> namespace is a host-enforced stable
// convention, so links survive plugin upgrades.
const viewerH5Path = "/x/linapro-live-manage/h5"

// Subscribe renders the ICS calendar stream of one live room.
func (c *ControllerV1) Subscribe(ctx context.Context, req *v1.SubscribeReq) (res *v1.SubscribeRes, err error) {
	feed, err := c.calendarSvc.Feed(ctx, toServiceInput(req))
	if err != nil {
		return nil, err
	}

	r := g.RequestFromCtx(ctx)
	r.Response.Header().Set("Content-Type", calendarContentMediaType)
	r.Response.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="linapro-live-%s.ics"`, icsFilename(req.RoomCode)))
	r.Response.Write(calendarsvc.RenderICS(feed, watchURL(r, req)))
	return nil, nil
}

// toServiceInput converts the request DTO into the service input. The tenant
// parameter stays nil when absent so the service can degrade to an empty feed
// per its own presence rules.
func toServiceInput(req *v1.SubscribeReq) calendarsvc.FeedInput {
	input := calendarsvc.FeedInput{
		RoomCode: req.RoomCode,
	}
	if req.TenantId != nil {
		tenantID := *req.TenantId
		input.TenantId = &tenantID
	}
	return input
}

// watchURL builds the absolute viewer page URL embedded into event
// descriptions from the incoming request origin, so feeds rendered behind any
// host or reverse proxy carry reachable links. The tenant parameter is
// forwarded when the request carried one, mirroring the play endpoint's
// explicit-tenant contract.
func watchURL(r *ghttp.Request, req *v1.SubscribeReq) string {
	var builder strings.Builder
	builder.WriteString(r.GetSchema())
	builder.WriteString("://")
	builder.WriteString(r.GetHost())
	builder.WriteString(viewerH5Path)
	builder.WriteString("?room=")
	builder.WriteString(url.QueryEscape(req.RoomCode))
	if req.TenantId != nil {
		builder.WriteString("&tenant=")
		builder.WriteString(url.QueryEscape(fmt.Sprintf("%d", *req.TenantId)))
	}
	return builder.String()
}

// icsFilename maps the room code into a filename-safe fragment for the
// Content-Disposition header.
func icsFilename(roomCode string) string {
	var builder strings.Builder
	for _, symbol := range strings.TrimSpace(roomCode) {
		switch {
		case symbol >= 'a' && symbol <= 'z',
			symbol >= 'A' && symbol <= 'Z',
			symbol >= '0' && symbol <= '9',
			symbol == '-', symbol == '_', symbol == '.':
			builder.WriteRune(symbol)
		default:
			builder.WriteByte('-')
		}
	}
	name := builder.String()
	if name == "" {
		return "room"
	}
	return name
}
