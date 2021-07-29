package endpoints

import (
	"errors"
	"net/http"
	"sync"

	"github.com/Armatorix/BitLyke/pkg/pg"
	"github.com/Armatorix/BitLyke/pkg/schema"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	db      *pg.DB
	counter map[string]int64
	sync.RWMutex
}

func NewHandler(db *pg.DB) *Handler {
	return &Handler{db: db, counter: make(map[string]int64)}
}

var statusOK = map[string]string{"status": "ok"}

func Healthcheck(c echo.Context) error {
	c.Logger().Info("healthcheck")
	return c.JSON(http.StatusOK, statusOK)
}

func (h *Handler) GetAllShorts(c echo.Context) error {
	ls, err := h.db.GetLinkShorts()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, ls)
}

type postShortRequest struct {
	ShortPath string `json:"short_path" validate:"required,ne=api,ne=counts,alphanum"`
	RealURL   string `json:"real_url" validate:"required,uri"`
}

func (h *Handler) CreateShort(c echo.Context) error {
	var req postShortRequest
	if err := c.Bind(&req); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	if err := c.Validate(req); err != nil {
		return c.NoContent(http.StatusBadRequest)
	}
	ls, err := h.db.InsertShort(&schema.ShortLink{
		ShortPath: req.ShortPath,
		RealURL:   req.RealURL,
	})
	if err != nil {
		if errors.Is(err, pg.ErrDuplicatedEntry) {
			return c.NoContent(http.StatusConflict)
		}
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.JSON(http.StatusCreated, ls)
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
	h.incOccurance(req.Link)
	return c.Redirect(http.StatusTemporaryRedirect, l.RealURL)
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

	err = h.db.DeleteShort(req.Link)
	if err != nil {
		if errors.Is(err, pg.ErrNotFound) {
			return c.NoContent(http.StatusNotFound)
		}
		return c.NoContent(http.StatusInternalServerError)
	}
	return nil
}

func (h *Handler) GetCounts(c echo.Context) error {
	h.RLock()
	defer h.RUnlock()
	return c.JSON(http.StatusOK, h.counter)
}

func (h *Handler) incOccurance(linkID string) {
	h.Lock()
	defer h.Unlock()
	h.counter[linkID]++
}
