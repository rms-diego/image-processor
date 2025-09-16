package image

import (
	"github.com/gin-gonic/gin"
	"github.com/rms-diego/image-processor/internal/database"
	"github.com/rms-diego/image-processor/internal/middleware"
	imageHandler "github.com/rms-diego/image-processor/internal/modules/image/handler"
	imageRepository "github.com/rms-diego/image-processor/internal/modules/image/repository"
	imageService "github.com/rms-diego/image-processor/internal/modules/image/service"
	"github.com/rms-diego/image-processor/pkg/gateway"
)

func RouteInit(g *gin.RouterGroup) {
	g.Use(middleware.AuthHandler())

	r := imageRepository.NewImageRepository(database.Db)
	s := imageService.NewService(gateway.S3Gateway, r)
	h := imageHandler.NewHandler(s)

	g.POST("/", h.UploadImage)
}
