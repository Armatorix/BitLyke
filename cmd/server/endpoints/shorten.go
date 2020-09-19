package endpoints

import (
	"errors"
	"net/http"

	"github.com/Armatorix/BitLyke/pkg/model"
	"github.com/Armatorix/BitLyke/pkg/pg"
	"github.com/labstack/echo/v4"
)

func (h *Handler) GetAllShorts(c echo.Context) error {
	ls, err := h.db.GetLinkShorts()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, ls)
}

type postShortRequest struct {
	model.ShortLink
}

func (h *Handler) CreateShort(c echo.Context) error {
	var req postShortRequest
	if err := c.Bind(&req); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	ls, err := h.db.InsertShort(&req.ShortLink)
	if err != nil {
		if errors.Is(err, pg.ErrAlreadyInUse) {
			return c.NoContent(http.StatusConflict)
		}
		return nil
	}
	return c.JSON(http.StatusOK, ls)
}

type getShortedRequest struct {
	Link string `params:"link" json:"-" validate:"required"`
}

func (h *Handler) GetShort(c echo.Context) error {
	var req getShortedRequest
	err := c.Bind(&req)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	l, err := h.db.GetDestinationLink(req.Link)
	if err != nil {
		if errors.Is(err, pg.ErrNotFound) {
			return c.NoContent(http.StatusNotFound)
		}
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.Redirect(http.StatusTemporaryRedirect, l.RealUrl)
}

type deleteShortedRequest struct {
	Link string `params:"link" json:"-" validate:"required"`
}

func (h *Handler) DeleteShort(c echo.Context) error {
	var req deleteShortedRequest
	err := c.Bind(&req)
	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	l, err := h.db.DeleteShort(req.Link)
	if err != nil {
		if errors.Is(err, pg.ErrNotFound) {
			return c.NoContent(http.StatusNotFound)
		}
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, l)
}
