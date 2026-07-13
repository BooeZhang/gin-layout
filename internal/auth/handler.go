package auth

import (
	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Login(c *gin.Context, req LoginReq) (*LoginRes, error) {
	return h.svc.Login(c.Request.Context(), req)
}

func (h *Handler) RefreshToken(c *gin.Context, req RefreshTokenReq) (*RefreshTokenRes, error) {
	return h.svc.RefreshToken(c.Request.Context(), req)
}

func (h *Handler) Logout(c *gin.Context, req LogoutReq) (*LogoutRes, error) {
	return h.svc.Logout(c.Request.Context(), req)
}
