package menu

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"

	"gin-layout/internal/domain"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) List(c *gin.Context, _ struct{}) ([]domain.MenuItem, error) {
	return h.svc.List(c.Request.Context())
}

func (h *Handler) Create(c *gin.Context, req CreateMenuReq) (CreateMenuRes, error) {
	return h.svc.Create(c.Request.Context(), req)
}

func (h *Handler) GetOne(c *gin.Context, _ struct{}) (*domain.MenuItem, error) {
	id := cast.ToInt64(c.Param("id"))
	if id == 0 {
		return nil, domain.ErrInvalidMenuID
	}
	return h.svc.GetOne(c.Request.Context(), id)
}

func (h *Handler) Update(c *gin.Context, req UpdateMenuReq) (UpdateMenuRes, error) {
	id := cast.ToInt64(c.Param("id"))
	if id == 0 {
		return UpdateMenuRes{}, domain.ErrInvalidMenuID
	}
	req.MenuID = id
	return h.svc.Update(c.Request.Context(), req)
}

func (h *Handler) Delete(c *gin.Context, _ struct{}) (struct{}, error) {
	id := cast.ToInt64(c.Param("id"))
	if id == 0 {
		return struct{}{}, domain.ErrInvalidMenuID
	}
	return struct{}{}, h.svc.Delete(c.Request.Context(), id)
}
