package main

import (
	"log"

	"github.com/Armatorix/BitLyke/pkg/config"
	"github.com/Armatorix/BitLyke/pkg/endpoints"
	"github.com/Armatorix/BitLyke/pkg/pg"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type V struct {
	Validator *validator.Validate
}

func (cv *V) Validate(i interface{}) error {
	return cv.Validator.Struct(i)
}

func main() {
	cfg, err := config.FromEnv()
	if err != nil {
		log.Fatalf("failed env 
	}

	db, err := pg.New(cfg.Postgres)
	if err != nil {
		log.Fatalf("failed connection test: %v", err)
	}

	e := echo.New()

	e.Use(middleware.Logger())
	e.Logger.SetLevel(cfg.Server.LogLevel)
	e.Validator = &V{validator.New()}
	e.Use(middleware.CORS())
	e.GET("/public/health-check", endpoints.Healthcheck)

	h := endpoints.NewHandler(db)
	e.GET("/counts", h.GetCounts)

	api := e.Group("/api")
	api.GET("", h.GetAllShorts)
	api.POST("", h.CreateShort)
	api.DELETE("/:link", h.DeleteShort)
	e.GET("/:link", h.GetShort)

	e.Logger.Fatal(e.Start(cfg.Server.Address).Error())
}
