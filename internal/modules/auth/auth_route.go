package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/rms-diego/image-processor/internal/database"
	authhandler "github.com/rms-diego/image-processor/internal/modules/auth/auth_handler"
	authrepository "github.com/rms-diego/image-processor/internal/modules/auth/auth_repository"
	authservice "github.com/rms-diego/image-processor/internal/modules/auth/auth_service"
)

func RouteInit(g *gin.RouterGroup) {
	r := authrepository.NewRepository(database.DB)
	s := authservice.NewService(r)
	h := authhandler.NewHandler(s)

	g.POST("/register", h.Register)
	g.POST("/", h.Login)
}
