package main

import (
	"log"

	"github.com/Armatorix/BitLyke/cmd/server/endpoints"
	"github.com/Armatorix/BitLyke/cmd/server/validator"
	"github.com/Armatorix/BitLyke/pkg/config"
	"github.com/Armatorix/BitLyke/pkg/pg"
	"github.com/avast/retry-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	cfg, err := config.FromEnv()
	if err != nil {
		log.Fatalf("failed env init: %v", err)
	}

	db := pg.New(cfg.Postgres)
	err = retry.Do(db.TestRequest)
	if err != nil {
		log.Fatalf("failed connection test: %v", err)
	}

	v, err := validator.New()
	if err != nil {
		log.Fatalf("failed validator creation: %v", err)
	}

	e := echo.New()

	e.Use(middleware.Logger())
	e.Logger.SetLevel(cfg.Server.LogLevel)
	e.Validator = v
	e.Use(middleware.CORS())
	e.GET("/public/health-check", endpoints.Healthcheck)

	api := e.Group("/api")

	h := endpoints.New(db)
	api.GET("", h.GetAllShorts)
	api.POST("", h.CreateShort)
	api.DELETE("/:link", h.DeleteShort)
	e.GET("/:link", h.GetShort)

	e.Logger.Fatal(e.Start(cfg.Server.Address).Error())
}
