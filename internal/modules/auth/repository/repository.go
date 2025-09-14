package authRepository

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/rms-diego/image-processor/internal/database"
	"github.com/rms-diego/image-processor/internal/validations"
)

type AuthRepositoryInterface interface {
	Register(user *validations.AuthRequest) error
	FindByUsername(username string) (*validations.UserFound, error)
}

type authRepository struct{}

func NewRepository() AuthRepositoryInterface {
	return &authRepository{}
}

func (r *authRepository) Register(user *validations.AuthRequest) error {
	query := goqu.Record{"username": user.Username, "password": user.Password}

	_, err := database.Db.From("users").
		Insert().
		Rows(query).
		Executor().
		Exec()

	if err != nil {
		return err
	}

	return nil
}

func (r *authRepository) FindByUsername(username string) (*validations.UserFound, error) {
	var user validations.UserFound

	found, err := database.Db.From("users").
		Select("*").
		Where(goqu.Ex{"username": username}).
		Limit(1).
		ScanStruct(&user)

	if err != nil {
		return nil, err
	}

	if !found {
		return nil, nil
	}

	return &user, nil
}
