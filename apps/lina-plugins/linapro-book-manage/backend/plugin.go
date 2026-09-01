// Package backend wires the linapro-book-manage source plugin into the host
// plugin registry.
package backend

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"

	"lina-core/pkg/plugin/pluginhost"
	bookmanage "lina-plugin-linapro-book-manage"
	bookcontroller "lina-plugin-linapro-book-manage/backend/internal/controller/book"
	borrowcontroller "lina-plugin-linapro-book-manage/backend/internal/controller/borrow"
	booksvc "lina-plugin-linapro-book-manage/backend/internal/service/book"
	borrowsvc "lina-plugin-linapro-book-manage/backend/internal/service/borrow"
)

// linapro-book-manage plugin constants.
const (
	// pluginID is the immutable identifier published by the embedded source plugin.
	pluginID = "linapro-book-manage"
)

// init registers the linapro-book-manage source plugin and its host callbacks.
func init() {
	plugin := pluginhost.NewDeclarations(pluginID)
	plugin.Assets().UseEmbeddedFiles(bookmanage.EmbeddedFiles)
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

// registerRoutes binds book-management routes through the published host
// middleware set.
func registerRoutes(ctx context.Context, registrar pluginhost.HTTPRegistrar) error {
	var (
		routes      = registrar.Routes()
		middlewares = routes.Middlewares()
		services    = registrar.Services()
	)
	if services == nil {
		return gerror.New("linapro-book-manage routes require host bizctx, tenant-filter, and user capability services")
	}
	tenantSvc := services.Tenant()
	if services.BizCtx() == nil ||
		tenantSvc == nil ||
		tenantSvc.Filter() == nil ||
		services.Users() == nil {
		return gerror.New("linapro-book-manage routes require host bizctx, tenant-filter, and user capability services")
	}
	bookSvc := booksvc.New(
		services.BizCtx(),
		tenantSvc,
		services.Users(),
	)
	borrowSvc := borrowsvc.New(
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
				group.Bind(bookcontroller.NewV1(bookSvc))
				group.Bind(borrowcontroller.NewV1(borrowSvc))
			})
		})
	})

	return nil
}
