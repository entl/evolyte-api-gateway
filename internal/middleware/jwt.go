package middleware

import (
	"strconv"
	"strings"

	"github.com/entl/evolyte-api-gateway/internal/errors"
	"github.com/entl/evolyte-api-gateway/internal/services"
	"github.com/labstack/echo/v4"
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

			claims, err := jwtSvc.ExtractClaims(token)
			if err != nil {
				return errors.ErrInvalidToken
			}

			c.Set("role", claims["role"])
			c.Set("user_id", claims["user_id"])

			c.Request().Header.Set("X-User-ID", strconv.FormatFloat(claims["user_id"].(float64), 'f', -1, 64))
			return next(c)
		}
	}
}
