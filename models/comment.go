package models

import "time"

type Comment struct {
	ID        uint         `gorm:"primaryKey" json:"id"`
    UserID    uint         `json:"user_id"`
    User      User         `json:"user"`
    PhotoID   uint         `json:"photo_id"`
    Photo     Photo        `json:"photo"`
    Message   string       `gorm:"not null" json:"message" validate:"required"`
    CreatedAt time.Time    `json:"created_at"`
    UpdatedAt time.Time    `json:"updated_at"`
}