package image

import (
	"github.com/gin-gonic/gin"
	"github.com/rms-diego/image-processor/internal/database"
	"github.com/rms-diego/image-processor/internal/middleware"
	s3Gateway "github.com/rms-diego/image-processor/internal/modules/gateway/s3"
	imageHandler "github.com/rms-diego/image-processor/internal/modules/image/handler"
	imageRepository "github.com/rms-diego/image-processor/internal/modules/image/repository"
	imageService "github.com/rms-diego/image-processor/internal/modules/image/service"
)

func RouteInit(g *gin.RouterGroup) {
	g.Use(middleware.AuthHandler())

	s3Gateway := s3Gateway.NewService()

	r := imageRepository.NewImageRepository(database.Db)
	s := imageService.NewService(s3Gateway, r)
	h := imageHandler.NewHandler(s)

	g.POST("/", h.UploadImage)
}
