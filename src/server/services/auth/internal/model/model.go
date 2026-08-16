package model

import "time"

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Username  string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"username"`
	Password  string    `gorm:"type:text;not null" json:"-"`
	Role      Role      `gorm:"type:varchar(20);not null;default:'user'" json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
