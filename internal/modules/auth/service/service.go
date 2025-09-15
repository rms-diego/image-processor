package authService

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rms-diego/image-processor/internal/config"
	authRepository "github.com/rms-diego/image-processor/internal/modules/auth/repository"
	"github.com/rms-diego/image-processor/internal/utils/exception"
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
		return exception.New(err.Error(), http.StatusInternalServerError)
	}

	if userFound != nil {
		return exception.New("User already exists", http.StatusBadRequest)
	}

	passwordBytes := []byte(payload.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)
	if err != nil {
		return exception.New(err.Error(), http.StatusInternalServerError)
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
		return nil, exception.New(err.Error(), http.StatusInternalServerError)
	}

	if userFound == nil {
		return nil, exception.New("user not found", http.StatusNotFound)
	}

	err = bcrypt.CompareHashAndPassword([]byte(userFound.Password), []byte(payload.Password))
	if err != nil {
		return nil, exception.New("Invalid credentials", http.StatusUnauthorized)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       userFound.ID,
		"username": userFound.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenStr, err := token.SignedString([]byte(config.Env.JWT_SECRET))
	if err != nil {
		return nil, exception.New(err.Error(), http.StatusInternalServerError)
	}

	return &tokenStr, nil
}
