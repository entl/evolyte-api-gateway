package proxy

import (
	"bytes"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/entl/evolyte-api-gateway/internal"
	"github.com/entl/evolyte-api-gateway/internal/handlers"
	"github.com/labstack/echo/v4"
)

// NewProxy returns an Echo handler that forwards requests to targetHost
func NewProxy(cfg *internal.GatewayConfig, jwtHandler *handlers.JWTHandler, cacheHandler *handlers.CacheHandler) echo.HandlerFunc {
	serviceMap := make(map[string]internal.ServiceConfig, len(cfg.Services))
	for _, svc := range cfg.Services {
		serviceMap[svc.Name] = svc
	}

	return func(c echo.Context) error {
		const apiPrefix = "/api/v1/"
		req := c.Request()
		fullPath := req.URL.Path
		path := strings.TrimPrefix(fullPath, apiPrefix)

		parts := strings.SplitN(path, "/", 2)
		if len(parts) < 1 {
			return echo.NewHTTPError(400, "Invalid request path")
		}

		serviceName, routePath := parts[0], parts[1]
		svc, ok := serviceMap[serviceName]
		if !ok {
			return echo.NewHTTPError(404, "Service not found")
		}

		targetURL, err := url.Parse(svc.Backend)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "invalid backend URL")
		}

		if !isPublicRoute(svc.PublicRoutes, req.Method, "/"+routePath) {
			userID, role, err := jwtHandler.Authorize(c)
			if err != nil {
				slog.Error("JWT authentication failed", "error", err)
				return err
			}
			c.Set("user_id", userID)
			c.Set("role", role)
			req.Header.Set("X-User-ID", userID)
			req.Header.Set("X-User-Role", role)
			slog.Info("JWT authenticated", "user_id", userID, "role", role)
		}

		val, err := cacheHandler.TryCache(c)
		if err == nil && val != nil {
			// If cache hit, return cached response
			return c.JSONBlob(http.StatusOK, val)
		}

		// Create a response recorder to capture the response
		rec := &ResponseRecorder{
			ResponseWriter: c.Response().Writer,
			Body:           new(bytes.Buffer),
		}
		c.Response().Writer = rec

		// Prepare reverse proxy
		proxy := httputil.NewSingleHostReverseProxy(targetURL)
		req.URL.Scheme = targetURL.Scheme
		req.URL.Host = targetURL.Host
		req.URL.Path = apiPrefix + routePath
		req.Host = targetURL.Host

		slog.Info(fmt.Sprintf("Proxying request: service=%s, path=%s, to=%s", svc.Name, req.URL.Path, targetURL.String()))

		proxy.ServeHTTP(c.Response(), req)

		if rec.StatusCode == http.StatusOK {
			// Cache the response if it was successful
			responseBody := rec.Body.Bytes()
			if cacheErr := cacheHandler.CacheResponse(c, responseBody, http.StatusOK); cacheErr != nil {
				slog.Error(fmt.Sprintf("Failed to cache response: %v", cacheErr))
			}
		}
		return nil

	}
}

func isPublicRoute(routes []internal.RouteConfig, method, path string) bool {
	for _, route := range routes {
		if route.Method == method && strings.HasPrefix(path, route.PathPrefix) {
			return true
		}
	}
	return false
}

// ResponseRecorder captures the response for caching
type ResponseRecorder struct {
	http.ResponseWriter
	Body       *bytes.Buffer
	StatusCode int
}

func (r *ResponseRecorder) Write(b []byte) (int, error) {
	// Write to both the original writer and our buffer
	n, err := r.ResponseWriter.Write(b)
	if err != nil {
		return n, err
	}
	return r.Body.Write(b[:n])
}

func (r *ResponseRecorder) WriteHeader(statusCode int) {
	r.StatusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func (r *ResponseRecorder) Flush() {
	if flusher, ok := r.ResponseWriter.(http.Flusher); ok {
		flusher.Flush()
	}
}
