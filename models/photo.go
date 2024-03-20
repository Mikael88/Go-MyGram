package models

import "time"

type Photo struct {
	ID        uint         `gorm:"primaryKey" json:"id"`
    Title     string       `gorm:"not null" json:"title" validate:"required"`
    Caption   string       `json:"caption"`
    PhotoURL  string       `gorm:"not null" json:"photo_url" validate:"required"`
    UserID    uint         `json:"user_id"`
    User      User         `json:"user"`
    CreatedAt time.Time    `json:"created_at"`
    UpdatedAt time.Time    `json:"updated_at"`
    Comments  []Comment    `json:"comments"`
}