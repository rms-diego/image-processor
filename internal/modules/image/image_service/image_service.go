package imageservice

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"sync"

	"github.com/chai2010/webp"
	"github.com/disintegration/gift"
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
	GetImageById(imageID string) (*string, error)
	GetImages(limit, page string) (*validations.ListImagesResponse, error)
	TransformImage(imageID string, payload *validations.TransformImageReqBody) error
	ProcessImage(file *[]byte, data *validations.TransformMessageQueue) error
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

	fileKey := fmt.Sprintf("%v.%v", uuid.New().String(), fh.Filename)

	fileBytes, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	location, err := s.s3Gateway.Upload(&fileKey, &fileBytes)
	if err != nil {
		return err
	}

	if err := s.repository.UploadImage(&userID, location, &fileKey); err != nil {
		return err
	}

	return nil
}

func (s *imageService) GetImageById(imageID string) (*string, error) {
	_, err := uuid.Parse(imageID)
	if err != nil {
		return nil, exception.New("invalid image id", http.StatusBadRequest)
	}

	image, err := s.repository.GetImageById(imageID)
	if err != nil {
		return nil, err
	}

	if image == nil {
		return nil, exception.New("image not found", http.StatusNotFound)
	}

	return &image.URL, nil
}

func (s *imageService) GetImages(limit, page string) (*validations.ListImagesResponse, error) {
	var l, p int

	err := func() error {
		if limit == "" {
			l = 10
			return nil
		}

		parsedLimit, err := parse.StringToInt(limit)
		if err != nil {
			return exception.New("limit must be a number", http.StatusBadRequest)
		}

		l = int(parsedLimit)
		return nil
	}()

	if err != nil {
		return nil, err
	}

	err = func() error {
		if page == "" {
			p = int(1)
			return nil
		}

		parsedPage, err := parse.StringToInt(page)
		if err != nil {
			return exception.New("page must be a number", http.StatusBadRequest)
		}

		if parsedPage <= 0 {
			return exception.New("page must be greater than 0", http.StatusBadRequest)
		}

		p = int((parsedPage - 1) * l)
		return nil
	}()

	if err != nil {
		return nil, err
	}

	images, count, err := s.repository.GetImages(&l, &p)
	if err != nil {
		return nil, err
	}

	return &validations.ListImagesResponse{
		TotalImages: *count,
		Data:        *images,
	}, nil
}

func (s *imageService) TransformImage(imageID string, payload *validations.TransformImageReqBody) error {
	_, err := uuid.Parse(imageID)
	if err != nil {
		return exception.New("invalid image id", http.StatusBadRequest)
	}

	image, err := s.repository.GetImageById(imageID)
	if err != nil {
		return err
	}

	if image == nil {
		return exception.New("image not found", http.StatusNotFound)
	}

	queueStruct := validations.TransformMessageQueue{
		S3Key:   image.S3Key,
		ImageID: image.ID,
		Payload: *payload,
	}

	j, err := json.Marshal(queueStruct)
	if err != nil {
		return err
	}

	m := string(j)
	if err := s.sqsGateway.SendMessage(&m); err != nil {
		return err
	}

	return nil
}

func (s *imageService) ProcessImage(file *[]byte, data *validations.TransformMessageQueue) error {
	src, _, err := image.Decode(bytes.NewReader(*file))
	if err != nil {
		fmt.Println("Error decoding image:", err)
		return err
	}

	var giftFilters []gift.Filter
	quality := 100
	format := "jpeg"

	if data.Payload.Quality != nil {
		quality = *data.Payload.Quality
	}

	if data.Payload.Resize != nil {
		giftFilters = append(
			giftFilters,
			gift.Resize(
				data.Payload.Resize.Width,
				data.Payload.Resize.Height,
				gift.LanczosResampling,
			),
		)
	}

	switch {
	case data.Payload.Filters != nil && data.Payload.Filters.Grayscale:
		giftFilters = append(giftFilters, gift.Grayscale())

	case data.Payload.Filters != nil && data.Payload.Filters.Sepia:
		giftFilters = append(giftFilters, gift.Sepia(100))
	}

	if data.Payload.Rotate != nil {
		giftFilters = append(
			giftFilters,
			gift.Rotate(
				float32(*data.Payload.Rotate),
				color.Transparent,
				gift.CubicInterpolation,
			),
		)
	}

	if data.Payload.Crop != nil {
		cropX := data.Payload.Crop.X
		cropY := data.Payload.Crop.Y
		cropW := data.Payload.Crop.Width
		cropH := data.Payload.Crop.Height

		rect := image.Rect(
			cropX, cropY,
			cropX+cropW, cropY+cropH,
		)

		giftFilters = append(
			giftFilters,
			gift.Crop(rect),
		)
	}

	g := gift.New(giftFilters...)
	dst := image.NewRGBA(g.Bounds(src.Bounds()))
	g.Draw(dst, src)

	if data.Payload.Format != nil {
		format = *data.Payload.Format
	}

	var fileBuf bytes.Buffer
	switch format {
	case "jpeg", "jpg":
		err = jpeg.Encode(&fileBuf, dst, &jpeg.Options{Quality: quality})

	case "png":
		encoder := png.Encoder{CompressionLevel: png.BestCompression}
		err = encoder.Encode(&fileBuf, dst)

	case "webp":
		err = webp.Encode(&fileBuf, dst, &webp.Options{
			Lossless: false,
			Quality:  float32(quality),
		})
	}

	if err != nil {
		fmt.Println("Error encoding image: ", err)
		return err
	}

	imageBytes, err := io.ReadAll(&fileBuf)
	if err != nil {
		fmt.Println("Error reading image bytes: ", err)
		return err
	}

	str := strings.Split(data.S3Key, ".")
	s3Id := str[0]
	filename := str[1]

	newS3Key := fmt.Sprintf("%v.%v.%v", s3Id, filename, format)

	wg := sync.WaitGroup{}
	wg.Add(2)

	newLocationFile := make(chan string, 1)
	errUpload := make(chan error, 1)
	errRemoveS3 := make(chan error, 1)

	go func() {
		defer wg.Done()
		location, err := gateway.S3Gateway.Upload(&newS3Key, &imageBytes)
		if err != nil {
			fmt.Println("Error uploading image to S3: ", err)

			newLocationFile <- ""
			errUpload <- err
			return
		}
		newLocationFile <- *location
		errUpload <- nil
	}()

	go func() {
		defer wg.Done()
		err := gateway.S3Gateway.RemoveObject(&data.S3Key)
		if err != nil {
			fmt.Println("Error removing image from S3: ", err)
			errRemoveS3 <- err
			return
		}

		errRemoveS3 <- nil
	}()

	wg.Wait()

	if err := <-errUpload; err != nil {
		return err
	}

	if err := <-errRemoveS3; err != nil {
		return err
	}

	url := <-newLocationFile
	if err := s.repository.UpdateImage(&data.ImageID, &newS3Key, &url); err != nil {
		return err
	}

	return nil
}
