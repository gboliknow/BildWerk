package models

import (
	"time"
)

type User struct {
	ID        string    `gorm:"type:uuid;primaryKey;" json:"id"`
	Name      string    `gorm:"type:varchar(100);not null" json:"name"`
	Email     string    `gorm:"uniqueIndex;not null" json:"email"`
	Password  string    `gorm:"not null" json:"-"`
	CreatedAt time.Time `json:"created_at"`
}


type Response struct {
	StatusCode   int         `json:"statusCode"`
	IsSuccessful bool        `json:"isSuccessful"`
	Message      string      `json:"message"`
	Data         interface{} `json:"data,omitempty"`
}
