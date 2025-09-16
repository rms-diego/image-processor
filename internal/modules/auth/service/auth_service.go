package authService

import (
	"net/http"

	authRepository "github.com/rms-diego/image-processor/internal/modules/auth/repository"
	"github.com/rms-diego/image-processor/internal/utils/exception"
	jwtutils "github.com/rms-diego/image-processor/internal/utils/jwt"
	"github.com/rms-diego/image-processor/internal/validations"
	"golang.org/x/crypto/bcrypt"
)

type AuthServiceInterface interface {
	Register(user *validations.AuthRequest) error
	Login(user *validations.AuthRequest) (*string, error)
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
		return err
	}

	if userFound != nil {
		return exception.New("User already exists", http.StatusBadRequest)
	}

	passwordBytes := []byte(payload.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	err = s.Repository.Register(&validations.AuthRequest{
		Username: payload.Username,
		Password: string(hashedPassword),
	})

	if err != nil {
		return exception.New("User already exists", http.StatusBadRequest)
	}

	return nil
}

func (s *authService) Login(payload *validations.AuthRequest) (*string, error) {
	userFound, err := s.Repository.FindByUsername(payload.Username)
	if err != nil {
		return nil, err
	}

	if userFound == nil {
		return nil, exception.New("user not found", http.StatusNotFound)
	}

	err = bcrypt.CompareHashAndPassword([]byte(userFound.Password), []byte(payload.Password))
	if err != nil {
		return nil, exception.New("Invalid credentials", http.StatusUnauthorized)
	}

	token, err := jwtutils.NewJwtUtils().GenerateToken(*userFound)
	if err != nil {
		return nil, err
	}

	return token, nil
}
