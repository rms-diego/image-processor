package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/rms-diego/image-processor/internal/database"
	authHandler "github.com/rms-diego/image-processor/internal/modules/auth/handler"
	authRepository "github.com/rms-diego/image-processor/internal/modules/auth/repository"
	authService "github.com/rms-diego/image-processor/internal/modules/auth/service"
)

func RouteInit(g *gin.RouterGroup) {
	r := authRepository.NewRepository(database.Db)
	s := authService.NewService(r)
	h := authHandler.NewHandler(s)

	g.POST("/register", h.Register)
	g.POST("/", h.Login)
}
