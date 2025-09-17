package imageservice

import (
	"encoding/json"
	"mime/multipart"
	"net/http"

	"github.com/google/uuid"
	repository "github.com/rms-diego/image-processor/internal/modules/image/image_repository"
	"github.com/rms-diego/image-processor/internal/utils/exception"
	"github.com/rms-diego/image-processor/internal/utils/parse"
	"github.com/rms-diego/image-processor/internal/validations"
	"github.com/rms-diego/image-processor/pkg/gateway"
)

type imageService struct {
	s3Gateway  gateway.S3GatewayInterface
	sqsGateway gateway.SqsGatewayInterface
	repository repository.ImageRepositoryInterface
}

type ImageServiceInterface interface {
	UploadImage(userID string, file *multipart.FileHeader) error
	GetImageById(imageId string) (*string, error)
	GetImages(limit, page string) (*validations.ListImagesResponse, error)
	TransformImage(imageId string, payload *validations.TransformImageReqBody) error
}

func NewService(s3Gateway gateway.S3GatewayInterface, sqsGateway gateway.SqsGatewayInterface, repository repository.ImageRepositoryInterface) ImageServiceInterface {
	return &imageService{
		s3Gateway:  s3Gateway,
		sqsGateway: sqsGateway,
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
	_, err := uuid.Parse(imageId)
	if err != nil {
		return nil, exception.New("invalid image id", http.StatusBadRequest)
	}

	image, err := s.repository.GetImageById(imageId)
	if err != nil {
		return nil, err
	}

	if image == nil {
		return nil, exception.New("image not found", http.StatusNotFound)
	}

	return &image.URL, nil
}

func (s *imageService) GetImages(limit, page string) (*validations.ListImagesResponse, error) {
	l, err := func() (*int, error) {
		if limit == "" {
			r := int(10)
			return &r, nil
		}

		parsedLimit, err := parse.StringToInt(limit)
		if err != nil {
			return nil, exception.New("limit must be a number", http.StatusBadRequest)
		}

		r := int(parsedLimit)
		return &r, nil
	}()

	if err != nil {
		return nil, err
	}

	p, err := func() (*int, error) {
		if page == "" {
			r := int(1)
			return &r, nil
		}

		parsedPage, err := parse.StringToInt(page)
		if err != nil {
			return nil, exception.New("page must be a number", http.StatusBadRequest)
		}

		if parsedPage <= 0 {
			return nil, exception.New("page must be greater than 0", http.StatusBadRequest)
		}

		r := int((parsedPage - 1) * *l)
		return &r, nil
	}()

	if err != nil {
		return nil, err
	}

	images, count, err := s.repository.GetImages(l, p)
	if err != nil {
		return nil, err
	}

	return &validations.ListImagesResponse{
		TotalImages: *count,
		Data:        *images,
	}, nil
}

func (s *imageService) TransformImage(imageId string, payload *validations.TransformImageReqBody) error {
	_, err := uuid.Parse(imageId)
	if err != nil {
		return exception.New("invalid image id", http.StatusBadRequest)
	}

	image, err := s.repository.GetImageById(imageId)
	if err != nil {
		return err
	}

	if image == nil {
		return exception.New("image not found", http.StatusNotFound)
	}

	queueJson := validations.TransformMessageQueue{
		S3Key:   image.S3Key,
		Payload: *payload,
	}

	jm, err := json.Marshal(queueJson)
	if err != nil {
		return err
	}

	queueMessage := string(jm)
	if err := s.sqsGateway.SendMessage(&queueMessage); err != nil {
		return err
	}

	return nil
}
