package handlers

import (
	"log/slog"
	"strconv"
	"strings"

	"github.com/entl/evolyte-api-gateway/internal/errors"
	"github.com/entl/evolyte-api-gateway/internal/services"
	"github.com/labstack/echo/v4"
)

type JWTHandler struct {
	jwtSvc *services.JwtService
}

func NewJWTHandler(jwtSvc *services.JwtService) *JWTHandler {
	return &JWTHandler{jwtSvc: jwtSvc}
}

func (h *JWTHandler) Authorize(c echo.Context) (userID string, role string, err error) {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		slog.Warn("Missing or malformed Authorization header")
		return "", "", errors.ErrHeaderRequired
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	valid, err := h.jwtSvc.ValidateToken(token)
	if err != nil || !valid {
		slog.Warn("Invalid JWT token", "error", err)
		return "", "", errors.ErrInvalidToken
	}

	claims, err := h.jwtSvc.ExtractClaims(token)
	if err != nil {
		slog.Warn("Failed to extract JWT claims", "error", err)
		return "", "", errors.ErrInvalidToken
	}

	roleVal, _ := claims["role"].(string)
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		slog.Warn("Missing or invalid user_id claim")
		return "", "", errors.ErrInvalidToken
	}
	userID = strconv.FormatFloat(userIDFloat, 'f', -1, 64)

	// Inject headers for downstream services
	c.Request().Header.Set("X-User-ID", userID)
	c.Request().Header.Set("X-User-Role", roleVal)
	c.Set("role", roleVal)
	c.Set("user_id", userID)

	slog.Info("JWT validated and claims extracted", "role", roleVal, "user_id", userID)
	return userID, roleVal, nil
}
