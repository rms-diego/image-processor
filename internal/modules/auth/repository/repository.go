package authRepository

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/rms-diego/image-processor/internal/database"
	"github.com/rms-diego/image-processor/internal/validations"
)

type AuthRepositoryInterface interface {
	Register(user *validations.RegisterRequest) error
}

type authRepository struct{}

func NewRepository() AuthRepositoryInterface {
	return &authRepository{}
}

func (r *authRepository) Register(user *validations.RegisterRequest) error {
	query := goqu.Record{"username": user.Username, "password": user.Password}

	_, err := database.Db.From("users").
		Insert().
		Rows(query).
		Executor().
		Exec()

	if err != nil {
		return nil
	}

	return nil
}
