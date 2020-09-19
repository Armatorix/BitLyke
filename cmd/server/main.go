package main

import (
	"log"

	"github.com/Armatorix/BitLyke/cmd/server/endpoints"
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

	e := echo.New()

	e.Use(middleware.Logger())
	e.Logger.SetLevel(cfg.Server.LogLevel)
	e.Use(middleware.CORS())
	e.GET("/public/health-check", endpoints.Healthcheck)

	shorts := e.Group("/short")

	h := endpoints.New(db)
	shorts.GET("", h.GetAllShorts)
	shorts.POST("", h.CreateShort)
	shorts.GET("/:link", h.GetShort)
	shorts.DELETE("/:link", h.DeleteShort)

	e.Logger.Fatal(e.Start(cfg.Server.Address).Error())
}
