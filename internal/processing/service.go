package processing

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"io"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/gboliknow/bildwerk/internal/bucket"
	"github.com/gboliknow/bildwerk/internal/models"
	"github.com/gboliknow/bildwerk/internal/store"
	"github.com/gboliknow/bildwerk/internal/utility"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

type ImageService struct {
	store  store.Store
	logger zerolog.Logger
	bucket *bucket.CloudinaryService
}

func NewImageService(s store.Store, logger zerolog.Logger, b *bucket.CloudinaryService) *ImageService {
	return &ImageService{store: s, logger: logger, bucket: b}
}

func (s *ImageService) UploadImage(file multipart.File, filename string, userID string, size int64) (*models.ImageResponse, *utility.AppError) {
	const maxSize int64 = 5 * 1024 * 1024
	allowedFormats := map[string]bool{
		"jpeg": true,
		"png":  true,
		"webp": true,
	}

	if size > maxSize {
		return nil, &utility.AppError{
			Message:    "File size exceeds 5MB limit",
			StatusCode: http.StatusBadRequest,
		}
	}

	data, err := io.ReadAll(file)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to read uploaded file")
		return nil, &utility.AppError{
			Message:    "Failed to process image",
			StatusCode: http.StatusInternalServerError,
		}
	}

	img, imgFormat, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to decode image")
		return nil, &utility.AppError{
			Message:    "Invalid or unsupported image format",
			StatusCode: http.StatusBadRequest,
		}
	}

	if !allowedFormats[strings.ToLower(imgFormat)] {
		return nil, &utility.AppError{
			Message:    "Only JPEG, PNG, and WebP formats are allowed",
			StatusCode: http.StatusBadRequest,
		}
	}

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	base64Str := base64.StdEncoding.EncodeToString(data)
	dataURI := fmt.Sprintf("data:image/%s;base64,%s", imgFormat, base64Str)

	fileUrl, err := s.bucket.UploadBase64(dataURI)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to upload image to bucket")
		return nil, &utility.AppError{
			Message:    "Failed to upload image",
			StatusCode: http.StatusInternalServerError,
		}
	}

	imageID := uuid.New().String()
	imageRecord := &models.Image{
		ID:       imageID,
		UserID:   userID,
		Filename: filename,
		Path:     fileUrl,
		Format:   imgFormat,
		Width:    width,
		Height:   height,
		Size:     size,
	}

	image, err := s.store.CreateImage(imageRecord)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to save image metadata to database")
		return nil, &utility.AppError{
			Message:    "Failed to save image metadata",
			StatusCode: http.StatusInternalServerError,
		}
	}

	return &models.ImageResponse{
		ID:        imageID,
		Filename:  filename,
		URL:       fileUrl,
		Format:    imgFormat,
		Width:     width,
		Height:    height,
		Size:      size,
		CreatedAt: image.CreatedAt,
	}, nil
}
