package imageService

import (
	"mime/multipart"

	s3Gateway "github.com/rms-diego/image-processor/internal/modules/gateway/s3"
	imageRepository "github.com/rms-diego/image-processor/internal/modules/image/repository"
)

type imageService struct {
	s3Gateway  s3Gateway.S3GatewayServiceInterface
	repository imageRepository.ImageRepositoryInterface
}

type ImageServiceInterface interface {
	UploadImage(userID string, file *multipart.FileHeader) error
}

func NewService(s3Gateway s3Gateway.S3GatewayServiceInterface, repository imageRepository.ImageRepositoryInterface) ImageServiceInterface {
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

	location, err := s.s3Gateway.Upload(fh, &f)
	if err != nil {
		return err
	}

	if err := s.repository.UploadImage(userID, *location); err != nil {
		return err
	}

	return nil
}
