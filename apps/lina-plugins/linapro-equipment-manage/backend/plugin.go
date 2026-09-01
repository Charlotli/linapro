// Package backend wires the linapro-equipment-manage source plugin into the
// host plugin registry.
package backend

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"

	"lina-core/pkg/plugin/pluginhost"
	equipmentmanage "lina-plugin-linapro-equipment-manage"
	equipmentcontroller "lina-plugin-linapro-equipment-manage/backend/internal/controller/equipment"
	maintenancecontroller "lina-plugin-linapro-equipment-manage/backend/internal/controller/maintenance"
	equipmentsvc "lina-plugin-linapro-equipment-manage/backend/internal/service/equipment"
	maintenancesvc "lina-plugin-linapro-equipment-manage/backend/internal/service/maintenance"
)

// linapro-equipment-manage plugin constants.
const (
	// pluginID is the immutable identifier published by the embedded source plugin.
	pluginID = "linapro-equipment-manage"
)

// init registers the linapro-equipment-manage source plugin and its host callbacks.
func init() {
	plugin := pluginhost.NewDeclarations(pluginID)
	plugin.Assets().UseEmbeddedFiles(equipmentmanage.EmbeddedFiles)
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

// registerRoutes binds equipment-management routes through the published host
// middleware set.
func registerRoutes(ctx context.Context, registrar pluginhost.HTTPRegistrar) error {
	var (
		routes      = registrar.Routes()
		middlewares = routes.Middlewares()
		services    = registrar.Services()
	)
	if services == nil {
		return gerror.New("linapro-equipment-manage routes require host bizctx, tenant-filter, and user capability services")
	}
	tenantSvc := services.Tenant()
	if services.BizCtx() == nil ||
		tenantSvc == nil ||
		tenantSvc.Filter() == nil ||
		services.Users() == nil {
		return gerror.New("linapro-equipment-manage routes require host bizctx, tenant-filter, and user capability services")
	}
	equipmentSvc := equipmentsvc.New(
		services.BizCtx(),
		tenantSvc,
		services.Users(),
	)
	maintenanceSvc := maintenancesvc.New(
		services.BizCtx(),
		tenantSvc,
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
				group.Bind(equipmentcontroller.NewV1(equipmentSvc))
				group.Bind(maintenancecontroller.NewV1(maintenanceSvc))
			})
		})
	})

	return nil
}
