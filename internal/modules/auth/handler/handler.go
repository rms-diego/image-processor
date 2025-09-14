package authHandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	usersService "github.com/rms-diego/image-processor/internal/modules/auth/service"
	"github.com/rms-diego/image-processor/internal/validations"
)

type AuthHandlerInterface interface {
	Register(c *gin.Context)
}

type authHandler struct {
	Service usersService.AuthServiceInterface
}

func NewHandler(service usersService.AuthServiceInterface) AuthHandlerInterface {
	return &authHandler{Service: service}
}

func (h *authHandler) Register(c *gin.Context) {
	var payload validations.RegisterRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Service.Register(&payload); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
