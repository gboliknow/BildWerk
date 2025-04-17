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
