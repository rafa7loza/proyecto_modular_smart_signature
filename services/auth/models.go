package auth

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FirstName  string `gorm:"size:256" json:"first_name"`
	LastName   string `gorm:"size:256" json:"last_name"`
	Email      string `gorm:"uniqueIndex; size:256" json:"email"`
	Password   string `gorm:"size:1024" json:"password"`
	IsVerified bool   `gorm:"DEFAULT: false" json:"is_verified"`
	Token      Token
}

type Token struct {
	gorm.Model
	EncodedToken   string `gorm:"size:1024"`
	ExpirationDate *time.Time
	UserID         uint
}

type UserDTO struct {
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Email      string `json:"email"`
	IsVerified bool   `json:"is_verified"`
}
