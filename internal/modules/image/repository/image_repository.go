package imageRepository

import (
	"github.com/doug-martin/goqu/v9"
)

type ImageRepositoryInterface interface {
	UploadImage(userId, fileUrl, s3Key *string) error
}

type imageRepository struct {
	database *goqu.Database
}

func NewImageRepository(database *goqu.Database) ImageRepositoryInterface {
	return &imageRepository{
		database: database,
	}
}

func (r *imageRepository) UploadImage(userId, fileUrl, s3Key *string) error {
	query := goqu.Record{
		"url":     *fileUrl,
		"user_id": *userId,
		"s3_key":  *s3Key,
	}

	_, err := r.database.From("images").
		Insert().
		Rows(query).
		Executor().
		Exec()

	if err != nil {
		return err
	}

	return nil
}
