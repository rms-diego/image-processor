package imagehandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rms-diego/image-processor/internal/utils/exception"
	jwtutils "github.com/rms-diego/image-processor/internal/utils/jwt"

	service "github.com/rms-diego/image-processor/internal/modules/image/image_service"
)

type imageHandler struct {
	service service.ImageServiceInterface
}

type ImageHandlerInterface interface {
	UploadImage(c *gin.Context)
	GetImageById(c *gin.Context)
	GetImages(c *gin.Context)
}

func NewHandler(service service.ImageServiceInterface) ImageHandlerInterface {
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

func (h *imageHandler) GetImageById(c *gin.Context) {
	imageId := c.Param("imageId")
	if imageId == "" {
		c.Error(exception.New("imageId is required", http.StatusBadRequest))
		return
	}

	image, err := h.service.GetImageById(imageId)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": image})
}

func (h *imageHandler) GetImages(c *gin.Context) {
	limit := c.Query("limit")
	page := c.Query("page")

	images, err := h.service.GetImages(limit, page)
	if err != nil {
		c.Error(err)
		return
	}

	if images.Data == nil {
		c.JSON(http.StatusOK, gin.H{
			"total_images": images.TotalImages,
			"data":         []any{},
		})

		return
	}

	c.JSON(http.StatusOK, images)
}
