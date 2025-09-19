package imagerepository

import (
	"github.com/doug-martin/goqu/v9"
	"github.com/rms-diego/image-processor/internal/validations"
)

type ImageRepositoryInterface interface {
	UploadImage(userID, fileUrl, s3Key *string) error
	UpdateImage(imageID, s3Key, url *string) error
	GetImageById(imageID string) (*validations.Image, error)
	GetImages(limit, page *int) (*validations.ManyImages, *int, error)
}

type imageRepository struct {
	database *goqu.Database
}

func NewImageRepository(database *goqu.Database) ImageRepositoryInterface {
	return &imageRepository{
		database: database,
	}
}

func (r *imageRepository) UploadImage(userID, fileUrl, s3Key *string) error {
	query := goqu.Record{
		"url":     *fileUrl,
		"user_id": *userID,
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

func (r *imageRepository) GetImageById(imageID string) (*validations.Image, error) {
	var image validations.Image

	found, err := r.database.From("images").
		Select("*").
		Where(goqu.Ex{"id": imageID}).
		ScanStruct(&image)

	if err != nil {
		return nil, err
	}

	if !found {
		return nil, nil
	}

	return &image, nil
}

func (r *imageRepository) GetImages(limit, page *int) (*validations.ManyImages, *int, error) {
	var images validations.ManyImages

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

func (r *imageRepository) UpdateImage(imageID, s3Key, url *string) error {
	query := goqu.Record{
		"s3_key": *s3Key,
		"url":    *url,
	}

	_, err := r.database.From("images").
		Where(goqu.Ex{"id": imageID}).
		Update().
		Set(query).
		Executor().
		Exec()

	if err != nil {
		return err
	}

	return nil
}
