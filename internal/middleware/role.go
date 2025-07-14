package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func RoleMiddleware(allowedRoles []string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userRole := c.Get("role")
			fmt.Println("User role:", userRole)
			for _, role := range allowedRoles {
				if userRole == strings.ToUpper(role) {
					c.Request().Header.Set("X-User-Role", userRole.(string))
					return next(c)
				}
			}
			return echo.NewHTTPError(http.StatusForbidden, "forbidden - insufficient role")
		}
	}
}
