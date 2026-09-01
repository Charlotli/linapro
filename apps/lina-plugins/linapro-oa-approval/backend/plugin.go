// Package backend wires the linapro-oa-approval source plugin into the host
// plugin registry.
package backend

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"

	"lina-core/pkg/plugin/pluginhost"
	oaapproval "lina-plugin-linapro-oa-approval"
	approvalcontroller "lina-plugin-linapro-oa-approval/backend/internal/controller/approval"
	flowcontroller "lina-plugin-linapro-oa-approval/backend/internal/controller/flow"
	approvalsvc "lina-plugin-linapro-oa-approval/backend/internal/service/approval"
	flowsvc "lina-plugin-linapro-oa-approval/backend/internal/service/flow"
)

// linapro-oa-approval plugin constants.
const (
	// pluginID is the immutable identifier published by the embedded source plugin.
	pluginID = "linapro-oa-approval"
)

// init registers the linapro-oa-approval source plugin and its host callbacks.
func init() {
	plugin := pluginhost.NewDeclarations(pluginID)
	plugin.Assets().UseEmbeddedFiles(oaapproval.EmbeddedFiles)
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

// registerRoutes binds OA approval routes through the published host
// middleware set.
func registerRoutes(ctx context.Context, registrar pluginhost.HTTPRegistrar) error {
	var (
		routes      = registrar.Routes()
		middlewares = routes.Middlewares()
		services    = registrar.Services()
	)
	if services == nil {
		return gerror.New("linapro-oa-approval routes require host bizctx, tenant-filter, user, notification, and i18n capability services")
	}
	tenantSvc := services.Tenant()
	if services.BizCtx() == nil ||
		tenantSvc == nil ||
		tenantSvc.Filter() == nil ||
		services.Users() == nil ||
		services.Notifications() == nil ||
		services.I18n() == nil {
		return gerror.New("linapro-oa-approval routes require host bizctx, tenant-filter, user, notification, and i18n capability services")
	}
	flowSvc := flowsvc.New(
		services.BizCtx(),
		tenantSvc,
		services.Users(),
	)
	approvalSvc := approvalsvc.New(
		services.BizCtx(),
		tenantSvc,
		services.Users(),
		services.Notifications(),
		services.I18n(),
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
				group.Bind(flowcontroller.NewV1(flowSvc))
				group.Bind(approvalcontroller.NewV1(approvalSvc))
			})
		})
	})

	return nil
}
