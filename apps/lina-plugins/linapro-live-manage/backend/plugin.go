// Package backend wires the linapro-live-manage source plugin into the host
// plugin registry.
package backend

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"

	"lina-core/pkg/plugin/pluginhost"
	livemanage "lina-plugin-linapro-live-manage"
	livecontroller "lina-plugin-linapro-live-manage/backend/internal/controller/live"
	liveroomcontroller "lina-plugin-linapro-live-manage/backend/internal/controller/liveroom"
	livesvc "lina-plugin-linapro-live-manage/backend/internal/service/live"
	liveroomsvc "lina-plugin-linapro-live-manage/backend/internal/service/liveroom"
)

// linapro-live-manage plugin constants.
const (
	// pluginID is the immutable identifier published by the embedded source plugin.
	pluginID = "linapro-live-manage"
)

// init registers the linapro-live-manage source plugin and its host callbacks.
func init() {
	plugin := pluginhost.NewDeclarations(pluginID)
	plugin.Assets().UseEmbeddedFiles(livemanage.EmbeddedFiles)
	if err := plugin.HTTP().RegisterRoutes(
		pluginhost.ExtensionPointHTTPRouteRegister,
		pluginhost.CallbackExecutionModeBlocking,
		registerRoutes,
	); err != nil {
		panic(err)
	}
	if err := pluginhost.RegisterSourcePlugin(plugin); err != nil {
		panic(err)
	}
}

// registerRoutes binds live-management routes through the published host
// middleware set.
func registerRoutes(ctx context.Context, registrar pluginhost.HTTPRegistrar) error {
	var (
		routes      = registrar.Routes()
		middlewares = routes.Middlewares()
		services    = registrar.Services()
	)
	if services == nil {
		return gerror.New("linapro-live-manage routes require host bizctx, tenant-filter, and user capability services")
	}
	tenantSvc := services.Tenant()
	if services.BizCtx() == nil ||
		tenantSvc == nil ||
		tenantSvc.Filter() == nil ||
		services.Users() == nil {
		return gerror.New("linapro-live-manage routes require host bizctx, tenant-filter, and user capability services")
	}
	roomSvc := liveroomsvc.New(
		services.BizCtx(),
		tenantSvc,
		services.Users(),
	)
	liveSvc := livesvc.New(
		services.BizCtx(),
		tenantSvc,
		services.Users(),
	)
	routes.Group(routes.APIPrefix(), func(group pluginhost.RouteGroup) {
		group.Group("/api/v1", func(group pluginhost.RouteGroup) {
			group.Middleware(
				middlewares.NeverDoneCtx(),
				middlewares.HandlerResponse(),
				middlewares.CORS(),
				middlewares.RequestBodyLimit(),
				middlewares.Ctx(),
			)
			group.Group("/", func(group pluginhost.RouteGroup) {
				group.Middleware(
					middlewares.Auth(),
					middlewares.Tenancy(),
					middlewares.Permission(),
				)
				group.Bind(liveroomcontroller.NewV1(roomSvc))
				group.Bind(livecontroller.NewV1(liveSvc))
			})
		})
	})

	return nil
}
