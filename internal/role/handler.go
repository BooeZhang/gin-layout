package role

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

func (h *Handler) List(c *gin.Context, req ListRoleReq) (domain.PageResult[RoleItem], error) {
	return h.svc.List(c.Request.Context(), req)
}

func (h *Handler) Create(c *gin.Context, req CreateRoleReq) (CreateRoleRes, error) {
	return h.svc.Create(c.Request.Context(), req)
}

func (h *Handler) GetOne(c *gin.Context, _ struct{}) (*RoleItem, error) {
	id := cast.ToInt64(c.Param("id"))
	if id == 0 {
		return nil, domain.ErrInvalidRoleID
	}
	return h.svc.GetOne(c.Request.Context(), id)
}

func (h *Handler) Update(c *gin.Context, req UpdateRoleReq) (UpdateRoleRes, error) {
	id := cast.ToInt64(c.Param("id"))
	if id == 0 {
		return UpdateRoleRes{}, domain.ErrInvalidRoleID
	}
	req.RoleID = id
	return h.svc.Update(c.Request.Context(), req)
}

func (h *Handler) Delete(c *gin.Context, _ struct{}) (struct{}, error) {
	id := cast.ToInt64(c.Param("id"))
	if id == 0 {
		return struct{}{}, domain.ErrInvalidRoleID
	}
	return struct{}{}, h.svc.Delete(c.Request.Context(), id)
}

func (h *Handler) GetAll(c *gin.Context, _ struct{}) ([]RoleItem, error) {
	return h.svc.GetAll(c.Request.Context())
}

func (h *Handler) UserAdd(c *gin.Context, req UserAddReq) (UserAddRes, error) {
	roleID := cast.ToInt64(c.Param("id"))
	if roleID == 0 {
		return UserAddRes{}, domain.ErrInvalidRoleID
	}
	req.RoleID = roleID
	return h.svc.UserAdd(c.Request.Context(), req)
}

func (h *Handler) UserRemove(c *gin.Context, req UserRemoveReq) (UserRemoveRes, error) {
	roleID := cast.ToInt64(c.Param("id"))
	if roleID == 0 {
		return UserRemoveRes{}, domain.ErrInvalidRoleID
	}
	req.RoleID = roleID
	return h.svc.UserRemove(c.Request.Context(), req)
}
