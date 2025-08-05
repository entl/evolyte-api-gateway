package routes

import (
	"fmt"
	"time"

	"github.com/entl/evolyte-api-gateway/internal"
	"github.com/entl/evolyte-api-gateway/internal/handlers"
	"github.com/entl/evolyte-api-gateway/internal/proxy"
	"github.com/entl/evolyte-api-gateway/internal/services"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
)

// NewRouter sets up all routes based on config
func NewRouter(cfg internal.Config) *echo.Echo {
	e := echo.New()
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.RateLimitDB,
	})
	cacheSvc := services.NewCacheService(redisClient)
	cacheHandler := handlers.NewCacheHandler(cacheSvc, 1*time.Minute)

	jwtSvc := services.NewJwtService(cfg.JWT.JWTSecret)
	jwtHandler := handlers.NewJWTHandler(jwtSvc)

	e.Any("/api/v1/*", proxy.NewProxy(&cfg.Gateway, jwtHandler, cacheHandler))

	return e
}
