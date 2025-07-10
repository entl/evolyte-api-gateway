package routes

import (
	"github.com/entl/evolyte-api-gateway/internal"
	"github.com/entl/evolyte-api-gateway/internal/middleware"
	"github.com/entl/evolyte-api-gateway/internal/proxy"
	"github.com/entl/evolyte-api-gateway/internal/services"
	"github.com/labstack/echo/v4"
)

// NewRouter sets up all routes based on config
func NewRouter(cfg internal.Config) *echo.Echo {
	e := echo.New()
	jwtSvc := services.NewJwtService(cfg.JWT.JWTSecret)

	for _, svc := range cfg.Gateway.Services {
		for _, r := range svc.Routes {
			handler := proxy.NewProxy(svc.Backend, r.FromPath, r.ToPath)
			if r.AuthRequired {
				e.Add(r.Method, r.FromPath, handler, middleware.JWTAuth(jwtSvc))
			} else {
				e.Add(r.Method, r.FromPath, handler)
			}
		}
	}

	return e
}
