package authRepository

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/rms-diego/image-processor/internal/validations"
)

type AuthRepositoryInterface interface {
	Register(user *validations.AuthRequest) error
	FindByUsername(username string) (*validations.UserFound, error)
}

type authRepository struct {
	database *goqu.Database
}

func NewRepository(database *goqu.Database) AuthRepositoryInterface {
	return &authRepository{
		database: database,
	}
}

func (r *authRepository) Register(user *validations.AuthRequest) error {
	query := goqu.Record{"username": user.Username, "password": user.Password}

	_, err := r.database.From("users").
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

	found, err := r.database.From("users").
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
