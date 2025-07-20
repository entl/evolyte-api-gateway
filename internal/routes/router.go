package routes

import (
	"fmt"

	"github.com/entl/evolyte-api-gateway/internal"
	"github.com/entl/evolyte-api-gateway/internal/middleware"
	"github.com/entl/evolyte-api-gateway/internal/proxy"
	"github.com/entl/evolyte-api-gateway/internal/services"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

// NewRouter sets up all routes based on config
func NewRouter(cfg internal.Config) *echo.Echo {
	e := echo.New()
	jwtSvc := services.NewJwtService(cfg.JWT.JWTSecret)
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.RateLimitDB,
	})
	cacheSvc := services.NewCacheService(redisClient)

	for _, svc := range cfg.Gateway.Services {
		for _, r := range svc.Routes {
			handler := proxy.NewProxy(svc.Backend, r.FromPath, r.ToPath)
			if r.AuthRequired {
				e.Add(r.Method, r.FromPath, handler, middleware.JWTAuth(jwtSvc), middleware.RoleMiddleware(r.AllowedRoles), middleware.CacheMiddleware(cacheSvc))
			} else {
				e.Add(r.Method, r.FromPath, handler, middleware.CacheMiddleware(cacheSvc))
			}
		}
	}

	return e
}
