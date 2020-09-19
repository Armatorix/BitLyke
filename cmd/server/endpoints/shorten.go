package endpoints

import (
	"errors"
	"net/http"

	"github.com/Armatorix/BitLyke/pkg/model"
	"github.com/Armatorix/BitLyke/pkg/pg"
	"github.com/labstack/echo/v4"
)

func (h *Handler) GetAllShorten(c echo.Context) error {
	ls, err := h.db.GetLinkShortens()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, ls)
}

type postShortenRequest struct {
	model.ShortenLink
}

func (h *Handler) CreateShorten(c echo.Context) error {
	var req postShortenRequest
	if err := c.Bind(&req); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	ls, err := h.db.InsertLinkShorten(&req.ShortenLink)
	if err != nil {
		if errors.Is(err, pg.ErrAlreadyInUse) {
			return c.NoContent(http.StatusConflict)
		}
		return nil
	}
	return c.JSON(http.StatusOK, ls)
}

func (h *Handler) GetShorten(c echo.Context) error {

	return nil
}

func (h *Handler) DeleteShorten(c echo.Context) error {

	return nil
}
