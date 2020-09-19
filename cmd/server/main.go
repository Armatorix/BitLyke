package main

import (
	"log"

	"github.com/Armatorix/BitLyke/cmd/server/endpoints"
	"github.com/Armatorix/BitLyke/pkg/config"
	"github.com/Armatorix/BitLyke/pkg/pg"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	cfg, err := config.FromEnv()
	if err != nil {
		log.Fatalf("failed env init: %v", err)
	}

	db := pg.New(cfg.Postgres)
	if err := db.TestRequest(); err != nil {
		log.Fatalf("failed connection test: %v", err)
	}

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.CORS())

	e.GET("/public/health-check", endpoints.Healthcheck)

	shortens := e.Group("/short")

	h := endpoints.New(db)
	shortens.GET("", h.GetAllShorten)
	shortens.POST("", h.CreateShorten)
	shortens.GET("/:link", h.GetShorten)
	shortens.DELETE("/:link", h.DeleteShorten)

}