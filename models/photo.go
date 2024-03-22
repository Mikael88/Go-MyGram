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

type PhotoResponse struct {
    ID        uint      `json:"id"`
    Title     string    `json:"title"`
    Caption   string    `json:"caption"`
    PhotoURL  string    `json:"photo_url"`
    UserID    uint      `json:"user_id"`
    CreatedAt time.Time `json:"created_at"`
  }