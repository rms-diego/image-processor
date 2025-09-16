package imageHandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rms-diego/image-processor/internal/utils/exception"
	jwtutils "github.com/rms-diego/image-processor/internal/utils/jwt"

	imageService "github.com/rms-diego/image-processor/internal/modules/image/service"
)

type imageHandler struct {
	service imageService.ImageServiceInterface
}

type ImageHandlerInterface interface {
	UploadImage(c *gin.Context)
}

func NewHandler(service imageService.ImageServiceInterface) ImageHandlerInterface {
	return &imageHandler{service: service}
}

func (h *imageHandler) UploadImage(c *gin.Context) {
	tokenDecoded, exists := c.Get("user")
	if !exists {
		c.Error(exception.New("unauthorized", http.StatusUnauthorized))
		return
	}

	file, err := c.FormFile("file")
	if err != nil || file == nil {
		c.Error(exception.New("file is required", http.StatusBadRequest))
		return
	}

	user := tokenDecoded.(jwtutils.JwtDecoded)
	if err := h.service.UploadImage(user.ID, file); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
