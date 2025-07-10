package main

import (
	"fmt"
	"github.com/entl/evolyte-api-gateway/internal"
	"github.com/entl/evolyte-api-gateway/internal/errors"
	"github.com/entl/evolyte-api-gateway/internal/routes"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
	"strconv"
)

func main() {
	cfg, err := internal.LoadConfig(".env", "config.yaml")
	if err != nil {
		log.Fatal("Failed to load config")
	}

	e := routes.NewRouter(*cfg)

	e.HTTPErrorHandler = func(err error, c echo.Context) {
		// Handle your custom AppError type
		if appErr, ok := err.(*errors.AppError); ok {
			c.JSON(appErr.Code, map[string]string{"error_code": strconv.Itoa(appErr.Code), "message": appErr.Message})
			return
		}

		// fallback for Echo built-in errors
		if he, ok := err.(*echo.HTTPError); ok {
			c.JSON(he.Code, map[string]interface{}{"error_code": he.Code, "message": he.Message})
			return
		}

		// fallback generic
		c.JSON(http.StatusInternalServerError, map[string]string{"error_code": "500", "message": "internal server error"})
	}

	addr := fmt.Sprintf(":%d", 1323)
	log.Printf("Starting API Gateway on %s", addr)
	log.Fatal(e.Start(addr))
}
