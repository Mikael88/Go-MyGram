package models

import "time"

type SocialMedia struct {
	ID             uint         `gorm:"primaryKey" json:"id"`
    Name           string       `gorm:"not null" json:"name" validate:"required"`
    SocialMediaURL string       `gorm:"not null" json:"social_media_url" validate:"required"`
    UserID         uint         `json:"user_id"`
    User           User         `json:"user"`
    CreatedAt      time.Time    `json:"created_at"`
    UpdatedAt      time.Time    `json:"updated_at"`
}