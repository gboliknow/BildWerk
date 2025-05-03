package models

import (
	"time"
)

type User struct {
	ID        string    `gorm:"type:uuid;primaryKey;" json:"id"`
	FirstName string    `gorm:"type:varchar(255);not null"`
	LastName  string    `gorm:"type:varchar(255);not null"`
	Email     string    `gorm:"uniqueIndex;not null" json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type UserResponse struct {
	ID        string    `gorm:"type:uuid;primaryKey;" json:"id"`
	FirstName string    `gorm:"type:varchar(255);not null"`
	LastName  string    `gorm:"type:varchar(255);not null"`
	Email     string    `gorm:"uniqueIndex;not null" json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type ImageResponse struct {
	ID         string    `json:"id"`
	Filename   string    `json:"filename"`
	URL        string    `json:"url"`
	Format     string    `json:"format"`
	Width      int       `json:"width"`
	Height     int       `json:"height"`
	Size       int64     `json:"size"`
	OriginalID *string   `json:"original_id,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

type Response struct {
	StatusCode   int         `json:"statusCode"`
	IsSuccessful bool        `json:"isSuccessful"`
	Message      string      `json:"message"`
	Data         interface{} `json:"data,omitempty"`
}

type AuthResponse struct {
	User  UserResponse `json:"user"`
	Token string       `json:"token"`
}

type Image struct {
	ID        string    `gorm:"type:uuid;primaryKey;" json:"id"`
	UserID    string    `gorm:"not null;index" json:"user_id"`
	Path      string    `gorm:"not null" json:"path"`     // URL to the image (on S3/Cloudinary)
	Filename  string    `gorm:"not null" json:"filename"` // Original filename
	Width     int       `json:"width"`                    // Optional: store image width
	Height    int       `json:"height"`                   // Optional: store image height
	Size      int64     `json:"size"`                     // Optional: image size in bytes
	Format    string    `json:"format"`                   // Optional: jpeg/png/webp, etc.
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
