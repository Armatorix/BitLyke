package main

import (
	"fmt"
	"html/template"
	"log"

	"github.com/Armatorix/BitLyke/pkg/config"
	"github.com/Armatorix/BitLyke/pkg/endpoints"
	"github.com/Armatorix/BitLyke/pkg/pg"
	"github.com/Armatorix/BitLyke/renderer"
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
		log.Fatalf("failed env init: %v", err)
	}

	db, err := pg.New(cfg.Postgres)
	if err != nil {
		log.Fatalf("failed connection test: %v", err)
	}
	if err := db.InitModels(); err != nil {
		log.Fatalf("failed model inits %v", err)
	}

	e := echo.New()
	e.Renderer = renderer.New(template.Must(template.ParseFiles("templates/index.html")))

	e.Use(middleware.Logger())
	e.Logger.SetLevel(cfg.Server.LogLevel)
	e.Validator = &V{validator.New()}
	e.Use(middleware.CORS())
	e.GET("/public/health-check", endpoints.Healthcheck)

	h := endpoints.NewHandler(db)
	e.GET("/", h.Mainpage)
	e.GET("/index.html", h.Mainpage)
	e.GET("/counts", h.GetCounts)

	api := e.Group("/api")
	api.GET("", h.GetAllShorts)
	api.POST("", h.CreateShort)
	api.DELETE("/:link", h.DeleteShort)
	e.GET("/:link", h.GetShort)

	e.Static("/assets", "public")
	e.File("/auth_config.json", "./auth_config.json")

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", cfg.Server.Port)).Error())
}
