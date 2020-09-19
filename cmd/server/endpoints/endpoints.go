package endpoints

import (
	"net/http"

	"github.com/Armatorix/BitLyke/pkg/pg"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	db *pg.DB
}

func New(db *pg.DB) *Handler {
	return &Handler{db}
}

func Healthcheck(c echo.Context) error {
	c.Logger().Info("healthcheck")
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}
