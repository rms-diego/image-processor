package image

import (
	"github.com/gin-gonic/gin"
	"github.com/rms-diego/image-processor/internal/database"
	"github.com/rms-diego/image-processor/internal/middleware"
	imagehandler "github.com/rms-diego/image-processor/internal/modules/image/image_handler"
	imagerepository "github.com/rms-diego/image-processor/internal/modules/image/image_repository"
	imageservice "github.com/rms-diego/image-processor/internal/modules/image/image_service"
	"github.com/rms-diego/image-processor/pkg/gateway"
)

func RouteInit(g *gin.RouterGroup) {
	g.Use(middleware.AuthHandler())

	r := imagerepository.NewImageRepository(database.DB)
	s := imageservice.NewService(gateway.S3Gateway, gateway.SqsGateway, r)
	h := imagehandler.NewHandler(s)

	g.GET("/:imageId", h.GetImageById)
	g.GET("/", h.GetImages)
	g.POST("/", h.UploadImage)
	g.PUT("/:imageId/transform", h.TransformImage)
}
