package imageservice

import (
	"mime/multipart"

	repository "github.com/rms-diego/image-processor/internal/modules/image/image_repository"
	"github.com/rms-diego/image-processor/pkg/gateway"
)

type imageService struct {
	s3Gateway  gateway.S3GatewayInterface
	repository repository.ImageRepositoryInterface
}

type ImageServiceInterface interface {
	UploadImage(userID string, file *multipart.FileHeader) error
	GetImageById(imageId string) (*string, error)
}

func NewService(s3Gateway gateway.S3GatewayInterface, repository repository.ImageRepositoryInterface) ImageServiceInterface {
	return &imageService{
		s3Gateway:  s3Gateway,
		repository: repository,
	}
}

func (s *imageService) UploadImage(userID string, fh *multipart.FileHeader) error {
	f, err := fh.Open()
	if err != nil {
		return err
	}

	defer f.Close()

	location, s3Key, err := s.s3Gateway.Upload(fh, &f)
	if err != nil {
		return err
	}

	if err := s.repository.UploadImage(&userID, location, s3Key); err != nil {
		return err
	}

	return nil
}

func (s *imageService) GetImageById(imageId string) (*string, error) {
	image, err := s.repository.GetImageById(imageId)
	if err != nil {
		return nil, err
	}

	return image, nil
}
