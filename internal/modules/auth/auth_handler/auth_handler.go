package authhandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	service "github.com/rms-diego/image-processor/internal/modules/auth/auth_service"
	"github.com/rms-diego/image-processor/internal/utils/exception"
	"github.com/rms-diego/image-processor/internal/validations"
)

type AuthHandlerInterface interface {
	Register(c *gin.Context)
	Login(c *gin.Context)
}

type authHandler struct {
	Service service.AuthServiceInterface
}

func NewHandler(service service.AuthServiceInterface) AuthHandlerInterface {
	return &authHandler{Service: service}
}

func (h *authHandler) Register(c *gin.Context) {
	var payload validations.AuthRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.Error(exception.New(err.Error(), http.StatusBadRequest))
		return
	}

	if err := h.Service.Register(&payload); err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h *authHandler) Login(c *gin.Context) {
	var payload validations.AuthRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.Error(exception.New(err.Error(), http.StatusBadRequest))
		return
	}

	token, err := h.Service.Login(&payload)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
