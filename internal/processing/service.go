package processing

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gboliknow/bildwerk/internal/bucket"
	"github.com/gboliknow/bildwerk/internal/models"
	"github.com/gboliknow/bildwerk/internal/store"
	"github.com/gboliknow/bildwerk/internal/utility"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"golang.org/x/image/webp"
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
	defer file.Close()

	maxSizeStr := os.Getenv("MAX_IMAGE_SIZE_MB")
	maxSizeMB, err := strconv.ParseInt(maxSizeStr, 10, 64)
	if err != nil || maxSizeMB <= 0 {
		maxSizeMB = 5
	}
	maxSize := maxSizeMB * 1024 * 1024
	if size > maxSize {
		return nil, &utility.AppError{
			Message:    "File size exceeds 5MB limit",
			StatusCode: http.StatusBadRequest,
		}
	}

	filename = sanitizeFilename(filename)

	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil || n < 8 {
		s.logger.Error().Err(err).Int("bytes_read", n).Msg("Failed to read file header")
		return nil, &utility.AppError{
			Message:    "Failed to read image content",
			StatusCode: http.StatusInternalServerError,
		}
	}

	fileHeader := fmt.Sprintf("% x", buffer[:8])
	contentType := http.DetectContentType(buffer)
	s.logger.Debug().
		Str("content_type", contentType).
		Str("file_header", fileHeader).
		Str("filename", filename).
		Msg("File signature detected")

	if !isValidImageSignature(buffer) {
		s.logger.Error().
			Str("content_type", contentType).
			Str("file_header", fileHeader).
			Msg("Invalid image signature")
		return nil, &utility.AppError{
			Message:    "Invalid image file signature",
			StatusCode: http.StatusBadRequest,
		}
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return nil, &utility.AppError{
			Message:    "Failed to reset file reader",
			StatusCode: http.StatusInternalServerError,
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

	var img image.Image
	var imgFormat string

	switch {
	case bytes.HasPrefix(buffer, []byte{0x89, 0x50, 0x4E, 0x47}): // PNG
		img, err = png.Decode(bytes.NewReader(data))
		imgFormat = "png"
	case bytes.HasPrefix(buffer, []byte{0xFF, 0xD8}): // JPEG
		img, err = jpeg.Decode(bytes.NewReader(data))
		imgFormat = "jpeg"
	case bytes.HasPrefix(buffer, []byte("RIFF")) && len(buffer) > 8 && bytes.HasPrefix(buffer[8:], []byte("WEBP")): // WEBP
		img, err = webp.Decode(bytes.NewReader(data))
		imgFormat = "webp"
	default:
		img, imgFormat, err = image.Decode(bytes.NewReader(data))
	}

	if err != nil {
		s.logger.Error().
			Err(err).
			Str("content_type", contentType).
			Str("file_header", fileHeader).
			Msg("Failed to decode image")
		return nil, &utility.AppError{
			Message:    "Invalid or unsupported image format",
			StatusCode: http.StatusBadRequest,
		}
	}

	allowedTypes := map[string]string{
		"png":  "image/png",
		"jpeg": "image/jpeg",
		"webp": "image/webp",
	}
	if expectedContentType, ok := allowedTypes[imgFormat]; !ok || expectedContentType != contentType {
		s.logger.Error().
			Str("detected_format", imgFormat).
			Str("content_type", contentType).
			Msg("Format/content-type mismatch")
		return nil, &utility.AppError{
			Message:    "Image format doesn't match file content",
			StatusCode: http.StatusBadRequest,
		}
	}

	bounds := img.Bounds()
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
		Width:    bounds.Dx(),
		Height:   bounds.Dy(),
		Size:     size,
	}

	if _, err := s.store.CreateImage(imageRecord); err != nil {
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
		Width:     bounds.Dx(),
		Height:    bounds.Dy(),
		Size:      size,
		CreatedAt: imageRecord.CreatedAt,
	}, nil
}

func isValidImageSignature(header []byte) bool {
	if bytes.HasPrefix(header, []byte{0x89, 0x50, 0x4E, 0x47}) {
		return true
	}
	if bytes.HasPrefix(header, []byte{0xFF, 0xD8}) {
		return true
	}
	if len(header) > 12 && bytes.HasPrefix(header, []byte("RIFF")) && bytes.HasPrefix(header[8:], []byte("WEBP")) {
		return true
	}
	return false
}

func sanitizeFilename(filename string) string {
	filename = strings.ReplaceAll(filename, "..", "")
	filename = strings.ReplaceAll(filename, "/", "")
	filename = strings.ReplaceAll(filename, "\\", "")
	filename = strings.TrimSpace(filename)
	return filename
}
