package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/entl/evolyte-api-gateway/internal/services"
	"github.com/labstack/echo/v4"
)

type CacheHandler struct {
	cacheSvc services.CacheService
	ttl      time.Duration
}

func NewCacheHandler(cacheSvc services.CacheService, ttl time.Duration) *CacheHandler {
	if ttl == 0 {
		ttl = 60 * time.Second // default TTL
	}
	return &CacheHandler{
		cacheSvc: cacheSvc,
		ttl:      ttl,
	}
}

// TryCache attempts to retrieve cached response, returns true if cache hit
func (ch *CacheHandler) TryCache(c echo.Context) ([]byte, error) {
	if c.Request().Method != "GET" {
		// Skip caching for non-GET requests
		return nil, fmt.Errorf("caching only supported for GET requests")
	}

	key := ch.generateCacheKey(c)

	// Try to get cached response
	cachedResponse, err := ch.cacheSvc.Get(c.Request().Context(), key)
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to get cached response for key %s: %v", key, err))
		return nil, fmt.Errorf("failed to get cached response: %w", err)
	}

	if cachedResponse != "" {
		// If cached response exists, return it
		slog.Info(fmt.Sprintf("Cache hit for key: %s", key))
		return []byte(cachedResponse), nil
	}

	slog.Info(fmt.Sprintf("Cache miss for key: %s", key))
	return nil, nil
}

// CacheResponse caches the response body with the generated key
func (ch *CacheHandler) CacheResponse(c echo.Context, responseBody []byte, statusCode int) error {
	if c.Request().Method != "GET" || statusCode != http.StatusOK {
		// Only cache successful GET requests
		return nil
	}

	key := ch.generateCacheKey(c)
	slog.Info(fmt.Sprintf("Caching response for key: %s", key))

	return ch.cacheSvc.Set(c.Request().Context(), key, string(responseBody), ch.ttl)
}

// generateCacheKey creates a cache key based on request path, query, and user context
func (ch *CacheHandler) generateCacheKey(c echo.Context) string {
	req := c.Request()
	key := "cache:" + req.URL.Path + req.URL.RawQuery

	if val := c.Get("user_id"); val != nil {
		if userID, ok := val.(string); ok {
			key += ":" + userID
		} else {
			// fallback for legacy float64 cases
			if userIDFloat, ok := val.(float64); ok {
				key += ":" + strconv.FormatFloat(userIDFloat, 'f', -1, 64)
			}
		}
	}

	return key
}
