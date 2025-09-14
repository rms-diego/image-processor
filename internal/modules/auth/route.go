package auth

import (
	"github.com/gin-gonic/gin"
	authHandler "github.com/rms-diego/image-processor/internal/modules/auth/handler"
	authRepository "github.com/rms-diego/image-processor/internal/modules/auth/repository"
	authService "github.com/rms-diego/image-processor/internal/modules/auth/service"
)

func RouteInit(r *gin.RouterGroup) {
	ur := authRepository.NewRepository()
	s := authService.NewService(ur)
	h := authHandler.NewHandler(s)

	r.POST("/register", h.Register)
}
