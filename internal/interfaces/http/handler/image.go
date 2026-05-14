package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	appimage "github.com/udevs/ai-chat/internal/application/image"
	"github.com/udevs/ai-chat/internal/interfaces/http/dto"
)

type ImageHandler struct {
	svc *appimage.Service
}

func NewImageHandler(svc *appimage.Service) *ImageHandler {
	return &ImageHandler{svc: svc}
}

// Generate godoc
// @Summary      Generate an image
// @Description  Generates one image from a text prompt using the configured OpenAI image model.
// @Tags         images
// @Accept       json
// @Produce      json
// @Param        body  body      dto.GenerateImageRequest  true  "Image prompt"
// @Success      200   {object}  dto.GenerateImageResponse
// @Failure      400   {object}  dto.ErrorResponse
// @Failure      500   {object}  dto.ErrorResponse
// @Router       /v1/images/generate [post]
func (h *ImageHandler) Generate(c *gin.Context) {
	var req dto.GenerateImageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{Error: err.Error()})
		return
	}

	out, err := h.svc.Generate(c.Request.Context(), appimage.GenerateInput{
		Prompt:  req.Prompt,
		Model:   req.Model,
		Size:    req.Size,
		Quality: req.Quality,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return
	}

	dataURL := ""
	if out.B64JSON != "" {
		dataURL = "data:image/png;base64," + out.B64JSON
	}
	c.JSON(http.StatusOK, dto.GenerateImageResponse{
		B64JSON: out.B64JSON,
		URL:     out.URL,
		DataURL: dataURL,
	})
}
