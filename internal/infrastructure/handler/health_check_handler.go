package handler

import (
	"allmygigs/internal/application/usecase"
	"net/http"

	"github.com/gin-gonic/gin"
)

// HealthHandler é a estrutura da handler
type HealthHandler struct {
	useCase *usecase.HealthCheckUseCase
}

func NewHealthHandler(useCase *usecase.HealthCheckUseCase) *HealthHandler {
	return &HealthHandler{useCase: useCase}
}

func (h *HealthHandler) CheckHealth(ctx *gin.Context) {
	result, err := h.useCase.Check()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":  "Falha ao realizar o health check",
			"detail": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, result)
}
