package proxy

import (
	"fmt"
	"log/slog"
	"net/http/httputil"
	"net/url"
	"path"
	"strings"

	"github.com/labstack/echo/v4"
)

// NewProxy returns an Echo handler that forwards requests to targetHost
func NewProxy(targetHost, fromPath, toPath string) echo.HandlerFunc {
	target, _ := url.Parse(targetHost)
	proxy := httputil.NewSingleHostReverseProxy(target)

	return func(c echo.Context) error {
		reqPath := c.Request().URL.Path

		// If your fromPath ends with "/*", do wildcard replacement
		if strings.HasSuffix(fromPath, "/*") && strings.HasSuffix(toPath, "/*") {
			fromPrefix := strings.TrimSuffix(fromPath, "*")
			toPrefix := strings.TrimSuffix(toPath, "*")

			// Get what comes after the prefix
			suffix := strings.TrimPrefix(reqPath, fromPrefix)
			suffix = strings.TrimPrefix(suffix, "/") // clean leading slash

			// Join to backend path
			newPath := path.Join(toPrefix, suffix)
			c.Request().URL.Path = newPath
			slog.Info(fmt.Sprintf("Proxying: %s => %s to %s", reqPath, newPath, targetHost))
		} else {
			// No wildcard: force toPath as fixed path
			c.Request().URL.Path = toPath
			slog.Info(fmt.Sprintf("Proxying: %s => %s to %s", reqPath, toPath, targetHost))
		}

		proxy.ServeHTTP(c.Response(), c.Request())
		return nil
	}
}
