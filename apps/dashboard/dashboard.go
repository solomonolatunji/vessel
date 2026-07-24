package dashboard

import (
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/labstack/echo/v4"
)

const (
	errDashboardNotBuilt = "Dashboard not built. Please run 'npm run build' in the dashboard."
	headerCacheControl   = "Cache-Control"
	cacheNoStore         = "no-cache, no-store, must-revalidate"
	cacheImmutable       = "public, max-age=31536000, immutable"
	cacheStandard        = "public, max-age=3600"
)

func init() {
	mime.AddExtensionType(".js", "application/javascript")
	mime.AddExtensionType(".css", "text/css")
	mime.AddExtensionType(".svg", "image/svg+xml")
}

func RegisterHandlers(e *echo.Echo) {
	e.GET("/*", handleStaticFS)
}

func handleStaticFS(c echo.Context) error {
	reqPath := filepath.Clean(c.Request().URL.Path)
	reqPath = strings.TrimPrefix(reqPath, "/")
	if reqPath == "" || reqPath == "." {
		reqPath = "index.html"
	}

	content, err := DistFS.ReadFile("dist/" + reqPath)
	if err != nil {
		indexContent, err := DistFS.ReadFile("dist/index.html")

		if err != nil {
			return c.String(http.StatusNotFound, errDashboardNotBuilt)
		}

		c.Response().Header().Set(headerCacheControl, cacheNoStore)
		return c.HTMLBlob(http.StatusOK, indexContent)
	}

	contentType := mime.TypeByExtension(filepath.Ext(reqPath))

	if contentType == "" {
		contentType = http.DetectContentType(content)
	}

	if strings.HasPrefix(reqPath, "assets/") {
		c.Response().Header().Set(headerCacheControl, cacheImmutable)
	} else {
		c.Response().Header().Set(headerCacheControl, cacheStandard)
	}

	return c.Blob(http.StatusOK, contentType, content)
}
