package middleware

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func RoleMiddleware(allowedRoles []string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userRole := c.Get("role")
			for _, role := range allowedRoles {
				if userRole == strings.ToUpper(role) {
					slog.Info("Role check passed", "user_role", userRole, "allowed_roles", allowedRoles)
					c.Request().Header.Set("X-User-Role", userRole.(string))
					return next(c)
				}
			}
			slog.Error("Forbidden - insufficient role", "user_role", userRole, "allowed_roles", allowedRoles)
			return echo.NewHTTPError(http.StatusForbidden, "forbidden - insufficient role")
		}
	}
}
