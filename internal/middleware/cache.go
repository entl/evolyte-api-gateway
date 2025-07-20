package middleware

import (
	"bytes"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/entl/evolyte-api-gateway/internal/services"
	"github.com/labstack/echo/v4"
)

func CacheMiddleware(cacheSvc services.CacheService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Request().Method != "GET" {
				// Skip caching for non-GET requests
				return next(c)
			}
			req := c.Request()
			key := "cache:" + req.URL.Path
			if c.Get("user_id") != nil {
				// Include user_id in cache key if available
				key += ":" + strconv.FormatFloat(c.Get("user_id").(float64), 'f', -1, 64)
			}

			// Try to get cached response
			cachedResponse, err := cacheSvc.Get(c.Request().Context(), key)
			if err != nil {
				slog.Error("Failed to get cached response for key %s: %v", key, err)
				return echo.NewHTTPError(500, "Failed to get cached response")
			}
			if cachedResponse != "" {
				// If cached response exists, return it
				slog.Info(fmt.Sprintf("Cache hit for key: %s", key))
				return c.JSONBlob(http.StatusOK, []byte(cachedResponse))
			}

			slog.Info(fmt.Sprintf("Cache miss for key: %s", key))
			rec := NewResponseRecorder(c.Response().Writer)
			c.Response().Writer = rec

			// If no cached response, proceed with the request
			if err := next(c); err != nil {
				return err
			}

			if (*rec.Result())["StatusCode"] == http.StatusOK {
				slog.Info(fmt.Sprintf("Caching response for key: %s", key))
				responseBody := string((*rec.Result())["Body"].([]byte))
				cacheSvc.Set(c.Request().Context(), key, responseBody, 60*time.Second)
			}

			return nil
		}
	}
}

type ResponseRecorder struct {
	http.ResponseWriter

	status       int
	body         bytes.Buffer
	headers      http.Header
	headerCopied bool
}

func NewResponseRecorder(w http.ResponseWriter) *ResponseRecorder {
	return &ResponseRecorder{
		ResponseWriter: w,
		headers:        make(http.Header),
	}
}

func (w *ResponseRecorder) Write(b []byte) (int, error) {
	w.copyHeaders()
	i, err := w.ResponseWriter.Write(b)
	if err != nil {
		return i, err
	}

	return w.body.Write(b[:i])
}

func (r *ResponseRecorder) copyHeaders() {
	if r.headerCopied {
		return
	}

	r.headerCopied = true
	copyHeaders(r.ResponseWriter.Header(), r.headers)
}

func (w *ResponseRecorder) WriteHeader(statusCode int) {
	w.copyHeaders()

	w.status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (r *ResponseRecorder) Result() *map[string]interface{} {
	r.copyHeaders()

	return &map[string]interface{}{
		"Header":     r.headers,
		"StatusCode": r.status,
		"Body":       r.body.Bytes(),
	}
}

func copyHeaders(src, dst http.Header) {
	for k, v := range src {
		for _, v := range v {
			dst.Set(k, v)
		}
	}
}
