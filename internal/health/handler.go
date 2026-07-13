package health

import (
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

type CheckRes struct {
	Status  string       `json:"status"`
	Message HealthStatus `json:"message"`
}

func (h *Handler) Check(c *gin.Context, _ struct{}) (CheckRes, error) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	healthStatus, err := h.svc.Check(ctx)
	if err != nil {
		return CheckRes{}, fmt.Errorf("health check: %w", err)
	}
	return CheckRes{Status: "ok", Message: healthStatus}, nil
}
