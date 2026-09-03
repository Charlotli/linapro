// Package backend wires the linapro-live-manage source plugin into the host
// plugin registry. Admin routes run behind auth, tenancy, and permission
// middlewares; the anonymous viewer H5 routes compose only the identity-free
// published middlewares, as documented in the OpenSpec change
// add-live-h5-player.
package backend

import (
	"context"
	"io/fs"
	"mime"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	"lina-core/pkg/plugin/pluginhost"
	livemanage "lina-plugin-linapro-live-manage"
	livecontroller "lina-plugin-linapro-live-manage/backend/internal/controller/live"
	liveroomcontroller "lina-plugin-linapro-live-manage/backend/internal/controller/liveroom"
	playcontroller "lina-plugin-linapro-live-manage/backend/internal/controller/play"
	subscribecontroller "lina-plugin-linapro-live-manage/backend/internal/controller/subscribe"
	calendarsvc "lina-plugin-linapro-live-manage/backend/internal/service/calendar"
	livesvc "lina-plugin-linapro-live-manage/backend/internal/service/live"
	liveroomsvc "lina-plugin-linapro-live-manage/backend/internal/service/liveroom"
	playsvc "lina-plugin-linapro-live-manage/backend/internal/service/play"
)

// linapro-live-manage plugin constants.
const (
	// pluginID is the immutable identifier published by the embedded source plugin.
	pluginID = "linapro-live-manage"

	// h5AssetRoot is the embedded frontend subdirectory hosting the viewer H5 assets.
	h5AssetRoot = "frontend/h5"

	// h5RoutePrefix is the plugin-public URL prefix serving the viewer H5 assets.
	h5RoutePrefix = "/h5"

	// h5IndexFile is served when the H5 root or a directory path is requested.
	h5IndexFile = "index.html"
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
// middleware set and the anonymous viewer play routes through an
// identity-free subset of the same published middlewares.
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
	// The viewer play service intentionally consumes only the tenant
	// capability: anonymous viewers carry no user identity, and the tenant is
	// an explicit request parameter validated against the tenant lifecycle.
	// The calendar subscription service shares the same anonymous contract.
	playInfoSvc := playsvc.New(tenantSvc)
	calendarSvc := calendarsvc.New(tenantSvc)
	routes.Group(routes.APIPrefix(), func(group pluginhost.RouteGroup) {
		registerViewerRoutes(group, middlewares, playInfoSvc, calendarSvc)
		registerAdminRoutes(group, middlewares, roomSvc, liveSvc)
	})

	return nil
}

// registerViewerRoutes binds the anonymous viewer routes: one JSON endpoint
// for play info, one raw endpoint for the ICS calendar stream, and one static
// group serving the embedded H5 assets. All groups omit auth, tenancy, and
// permission middlewares because viewers carry no JWT; the services derive
// their tenant from the explicit request parameter instead of the tenancy
// middleware.
func registerViewerRoutes(
	group pluginhost.RouteGroup,
	middlewares pluginhost.RouteMiddlewares,
	playInfoSvc playsvc.Service,
	calendarSvc calendarsvc.Service,
) {
	group.Group("/api/v1", func(group pluginhost.RouteGroup) {
		group.Middleware(
			middlewares.NeverDoneCtx(),
			middlewares.CORS(),
			middlewares.RequestBodyLimit(),
			middlewares.Ctx(),
		)
		group.Group("/", func(group pluginhost.RouteGroup) {
			group.Middleware(middlewares.HandlerResponse())
			group.Bind(playcontroller.NewV1(playInfoSvc))
		})
		// The subscribe endpoint answers with a raw RFC 5545 stream that
		// calendar clients cannot parse inside the JSON envelope, so it binds
		// outside the HandlerResponse wrapper per the documented design.
		group.Bind(subscribecontroller.NewV1(calendarSvc))
	})
	group.Group(h5RoutePrefix, func(group pluginhost.RouteGroup) {
		group.Middleware(
			middlewares.NeverDoneCtx(),
			middlewares.CORS(),
			middlewares.Ctx(),
		)
		// The wildcard route also matches the bare /h5 and /h5/ forms with an
		// empty any parameter; the handler falls back to the index asset.
		group.GET("/*any", newH5AssetHandler(""))
	})
}

// registerAdminRoutes binds the authenticated management controllers through
// the full host governance chain.
func registerAdminRoutes(
	group pluginhost.RouteGroup,
	middlewares pluginhost.RouteMiddlewares,
	roomSvc liveroomsvc.Service,
	liveSvc livesvc.Service,
) {
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
}

// newH5AssetHandler serves one embedded viewer asset. A wildcard match on the
// bare /h5 or /h5/ form carries an empty path and resolves to the H5 index;
// any other wildcard value resolves the cleaned request path inside the
// embedded H5 subtree.
func newH5AssetHandler(exactName string) ghttp.HandlerFunc {
	return func(r *ghttp.Request) {
		assetPath := exactName
		if assetPath == "" {
			assetPath = normalizeH5AssetPath(r.Get("any").String())
			if assetPath == "" {
				assetPath = h5IndexFile
			}
		}
		content, err := fs.ReadFile(livemanage.EmbeddedFiles, path.Join(h5AssetRoot, assetPath))
		if err != nil {
			if !os.IsNotExist(err) {
				g.Log().Errorf(r.Context(), "read viewer H5 asset %s: %v", assetPath, err)
			}
			r.Response.WriteStatus(http.StatusNotFound)
			return
		}
		contentType := mime.TypeByExtension(path.Ext(assetPath))
		if contentType == "" {
			contentType = http.DetectContentType(content)
		}
		r.Response.Header().Set("Content-Type", contentType)
		r.Response.Write(content)
	}
}

// normalizeH5AssetPath cleans one wildcard request path into a safe embedded
// asset path, rejecting traversal and empty targets.
func normalizeH5AssetPath(raw string) string {
	cleaned := path.Clean("/" + strings.ReplaceAll(strings.TrimSpace(raw), "\\", "/"))
	cleaned = strings.TrimPrefix(cleaned, "/")
	if cleaned == "" || cleaned == "." || strings.HasPrefix(cleaned, "../") {
		return ""
	}
	return cleaned
}
