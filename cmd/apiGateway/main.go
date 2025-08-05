package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/entl/evolyte-api-gateway/internal"
	"github.com/entl/evolyte-api-gateway/internal/errors"
	"github.com/entl/evolyte-api-gateway/internal/routes"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	ctx := context.Background()
	cfg, err := internal.LoadConfig(".env", "config.docker.yaml")
	if err != nil {
		log.Fatal("Failed to load config")
	}

	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		cfg.Postgres.Username,
		cfg.Postgres.Password,
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.Database,
	)
	conn, err := pgx.Connect(context.Background(), connectionString)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer conn.Close(ctx)

	e := routes.NewRouter(*cfg)
	e.Use(echoprometheus.NewMiddleware("evolyte_gateway"))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
	}))
	e.GET("/metrics", echoprometheus.NewHandler())

	e.HTTPErrorHandler = func(err error, c echo.Context) {
		if appErr, ok := err.(*errors.AppError); ok {
			c.JSON(appErr.Code, map[string]string{"error_code": strconv.Itoa(appErr.Code), "message": appErr.Message})
			return
		}

		if he, ok := err.(*echo.HTTPError); ok {
			c.JSON(he.Code, map[string]any{"error_code": he.Code, "message": he.Message})
			return
		}

		// fallback generic
		c.JSON(http.StatusInternalServerError, map[string]string{"error_code": "500", "message": "internal server error"})
	}

	addr := fmt.Sprintf(":%d", 8080)
	log.Printf("Starting API Gateway on %s", addr)
	log.Fatal(e.Start(addr))
}
