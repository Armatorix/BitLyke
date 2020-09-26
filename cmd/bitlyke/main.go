package main

import (
	"log"

	"github.com/Armatorix/BitLyke/pkg/config"
	"github.com/Armatorix/BitLyke/pkg/endpoints"
	"github.com/Armatorix/BitLyke/pkg/pg"
	"github.com/Armatorix/BitLyke/pkg/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	cfg, err := config.FromEnv()
	if err != nil {
		log.Fatalf("failed env init: %v", err)
	}

	db, err := pg.New(cfg.Postgres)
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

	h := endpoints.NewHandler(db)
	api.GET("", h.GetAllShorts)
	api.POST("", h.CreateShort)
	api.DELETE("/:link", h.DeleteShort)
	e.GET("/:link", h.GetShort)

	e.Logger.Fatal(e.Start(cfg.Server.Address).Error())
}