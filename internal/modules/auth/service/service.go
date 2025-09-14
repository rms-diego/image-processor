package authService

import (
	authRepository "github.com/rms-diego/image-processor/internal/modules/auth/repository"
	"github.com/rms-diego/image-processor/internal/validations"
	"golang.org/x/crypto/bcrypt"
)

type AuthServiceInterface interface {
	Register(user *validations.RegisterRequest) error
}

type authService struct {
	Repository authRepository.AuthRepositoryInterface
}

func NewService(repository authRepository.AuthRepositoryInterface) AuthServiceInterface {
	return &authService{Repository: repository}
}

func (s *authService) Register(user *validations.RegisterRequest) error {
	passwordBytes := []byte(user.Password)

	hashedPassword, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	err = s.Repository.Register(&validations.RegisterRequest{
		Username: user.Username,
		Password: string(hashedPassword),
	})

	if err != nil {
		return err
	}

	return nil
}
