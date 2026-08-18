package models

import "time"

const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

type User struct {
	UserID       string    `gorm:"primaryKey;size:100" json:"userid"`
	Email        string    `gorm:"uniqueIndex;not null;size:255" json:"email"`
	PasswordHash string    `gorm:"not null;size:255" json:"-"`
	Role         string    `gorm:"not null;size:20" json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
