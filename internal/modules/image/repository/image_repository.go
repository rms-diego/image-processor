package imageRepository

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/rms-diego/image-processor/internal/database"
)

type ImageRepositoryInterface interface {
	UploadImage(userId, fileUrl string) error
}

type imageRepository struct{}

func NewImageRepository() ImageRepositoryInterface {
	return &imageRepository{}
}

func (r *imageRepository) UploadImage(userId, fileUrl string) error {
	query := goqu.Record{
		"url":     fileUrl,
		"user_id": userId,
	}

	_, err := database.Db.From("images").
		Insert().
		Rows(query).
		Executor().
		Exec()

	if err != nil {
		return err
	}

	return nil
}
