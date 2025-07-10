package middleware

import (
	"github.com/entl/evolyte-api-gateway/internal/errors"
	"github.com/entl/evolyte-api-gateway/internal/services"
	"github.com/labstack/echo/v4"
	"strings"
)

func JWTAuth(jwtSvc *services.JwtService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			auth := c.Request().Header.Get("Authorization")
			if auth == "" {
				return errors.ErrHeaderRequired
			}
			token := strings.TrimPrefix(auth, "Bearer ")
			valid, err := jwtSvc.ValidateToken(token)
			if err != nil || !valid {
				return errors.ErrInvalidToken
			}
			return next(c)
		}
	}
}
