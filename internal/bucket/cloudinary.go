package bucket

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type CloudinaryConfig struct {
	Folder        string
	ResourceType  string
	MaxSizeMB     int
	DefaultTags   []string
	CloudinaryUrl string
}

type CloudinaryService struct {
	cld    *cloudinary.Cloudinary
	Config CloudinaryConfig
}

func NewCloudinaryService(config CloudinaryConfig) *CloudinaryService {
	cld, err := cloudinary.NewFromURL(config.CloudinaryUrl)
	if err != nil {
		log.Fatalf("Failed to initialize Cloudinary: %v", err)
	}
	return &CloudinaryService{cld: cld, Config: config}
}

func (c *CloudinaryService) UploadBase64(base64Str string) (string, error) {

	timestamp := time.Now().UnixNano()
	uniqueID := fmt.Sprintf("img_%d", timestamp)
	uploadParams := uploader.UploadParams{
		PublicID:     uniqueID,
		Folder:       "images",
		ResourceType: "image",
		Tags:         []string{"api_upload", "event"},
	}

	uploadResult, err := c.cld.Upload.Upload(context.Background(), base64Str, uploadParams)
	if err != nil {
		return "", fmt.Errorf("failed to upload image: %w", err)
	}
	return uploadResult.SecureURL, nil
}
