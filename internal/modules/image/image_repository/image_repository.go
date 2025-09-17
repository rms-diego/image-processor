package imagerepository

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/rms-diego/image-processor/internal/validations"
)

type ImageRepositoryInterface interface {
	UploadImage(userId, fileUrl, s3Key *string) error
	GetImageById(imageId string) (*string, error)
	GetImages(limit, page *int) (*validations.ImagesFound, *int, error)
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

func (r *imageRepository) GetImageById(imageId string) (*string, error) {
	var image string

	found, err := r.database.From("images").
		Select("url").
		Where(goqu.Ex{"id": imageId}).
		ScanVal(&image)

	if err != nil {
		return nil, err
	}

	if !found {
		return nil, nil
	}

	return &image, nil
}

func (r *imageRepository) GetImages(limit, page *int) (*validations.ImagesFound, *int, error) {
	var images validations.ImagesFound

	err := r.database.From("images").
		Select("id", "url", "created_at").
		Limit(uint(*limit)).
		Offset(uint(*page)).
		ScanStructs(&images)

	if err != nil {
		return nil, nil, err
	}

	var count int
	_, err = r.database.From("images").
		Select(goqu.COUNT("*").As("count")).
		ScanVal(&count)

	if err != nil {
		return nil, nil, err
	}

	return &images, &count, nil
}
