package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/entl/evolyte-api-gateway/internal"
	"github.com/entl/evolyte-api-gateway/internal/errors"
	"github.com/entl/evolyte-api-gateway/internal/routes"
	"github.com/labstack/echo/v4"
)

func main() {
	cfg, err := internal.LoadConfig(".env", "config.docker.yaml")
	if err != nil {
		log.Fatal("Failed to load config")
	}

	e := routes.NewRouter(*cfg)

	e.HTTPErrorHandler = func(err error, c echo.Context) {
		if appErr, ok := err.(*errors.AppError); ok {
			c.JSON(appErr.Code, map[string]string{"error_code": strconv.Itoa(appErr.Code), "message": appErr.Message})
			return
		}

		if he, ok := err.(*echo.HTTPError); ok {
			c.JSON(he.Code, map[string]interface{}{"error_code": he.Code, "message": he.Message})
			return
		}

		// fallback generic
		c.JSON(http.StatusInternalServerError, map[string]string{"error_code": "500", "message": "internal server error"})
	}

	addr := fmt.Sprintf(":%d", 8080)
	log.Printf("Starting API Gateway on %s", addr)
	log.Fatal(e.Start(addr))
}
