package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID 			uint 		`gorm:"primaryKey" json:"id"`
	Username 	string 		`gorm:"unique;not null" json:"username" validate:"required"`
	Email 		string 		`gorm:"unique;not null" json:"email" validate:"required,email"`
	Password 	string 		`gorm:"not null" json:"password" validate:"required,min=6"`
	Age 		int 		`gorm:"not null" json:"age" validate:"required,min=8"`
	CreatedAt 	time.Time 	`json:"created_at"`
	UpdateAt 	time.Time 	`json:"updated_at"`
	Photos		[]Photo 	`json:"photos"`
	Comments 	[]Comment 	`json:"comments"`
	SocialMedias []SocialMedia `json:"social_medias"`
}

type UserResponse struct {
    Age      int    `json:"age"`
    Email    string `json:"email"`
    ID       uint   `json:"id"`
    Username string `json:"username"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}