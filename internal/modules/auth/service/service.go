package authService

import (
	"net/http"

	authRepository "github.com/rms-diego/image-processor/internal/modules/auth/repository"
	"github.com/rms-diego/image-processor/internal/utils/exception"
	"github.com/rms-diego/image-processor/internal/validations"
	"golang.org/x/crypto/bcrypt"
)

type AuthServiceInterface interface {
	Register(user *validations.AuthRequest) error
}

type authService struct {
	Repository authRepository.AuthRepositoryInterface
}

func NewService(repository authRepository.AuthRepositoryInterface) AuthServiceInterface {
	return &authService{Repository: repository}
}

func (s *authService) Register(payload *validations.AuthRequest) error {
	userFound, err := s.Repository.FindByUsername(payload.Username)
	if err != nil {
		return exception.New(err.Error(), http.StatusInternalServerError, &err)
	}

	if userFound != nil {
		return exception.New("User already exists", http.StatusBadRequest, &err)
	}

	passwordBytes := []byte(payload.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)
	if err != nil {
		return exception.New(err.Error(), http.StatusInternalServerError, &err)
	}

	err = s.Repository.Register(&validations.AuthRequest{
		Username: payload.Username,
		Password: string(hashedPassword),
	})

	if err != nil {
		return exception.New("User already exists", http.StatusBadRequest, &err)
	}

	return nil
}
