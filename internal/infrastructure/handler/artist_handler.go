package handler

import (
	"allmygigs/internal/application/usecase"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TopArtistsHandler struct {
	service *usecase.ArtistUsecase
}

func NewArtistHandler(service *usecase.ArtistUsecase) *TopArtistsHandler {
	return &TopArtistsHandler{
		service: service,
	}
}

func (h *TopArtistsHandler) GetQuery(ctx *gin.Context) {
	user := ctx.Query("user")
	period := ctx.Query("period")
	limit := ctx.Query("limit")

	artists, err := h.service.GetUserTopArtists(user, period, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("Erro ao obter top artistas: %v", err),
		})

		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"artists": artists,
	})
}
