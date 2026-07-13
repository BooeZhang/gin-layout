package user

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

func (h *Handler) List(c *gin.Context, req ListUserReq) (domain.PageResult[UserItem], error) {
	return h.svc.List(c.Request.Context(), req)
}

func (h *Handler) Create(c *gin.Context, req CreateUserReq) (CreateUserRes, error) {
	return h.svc.Create(c.Request.Context(), req)
}

func (h *Handler) GetDetails(c *gin.Context, _ struct{}) (UserItem, error) {
	return h.svc.GetDetails(c.Request.Context())
}

func (h *Handler) Update(c *gin.Context, req UpdateUserReq) (UpdateUserRes, error) {
	id := cast.ToInt64(c.Param("id"))
	if id == 0 {
		return UpdateUserRes{}, domain.ErrInvalidUserID
	}
	req.UserID = id
	return h.svc.Update(c.Request.Context(), req)
}

func (h *Handler) Delete(c *gin.Context, _ struct{}) (struct{}, error) {
	id := cast.ToInt64(c.Param("id"))
	if id == 0 {
		return struct{}{}, domain.ErrInvalidUserID
	}
	return struct{}{}, h.svc.Delete(c.Request.Context(), id)
}

func (h *Handler) GetMenus(c *gin.Context, _ struct{}) ([]domain.MenuItem, error) {
	return h.svc.GetCurrentUserMenus(c.Request.Context())
}
