package config

import "github.com/Mikael88/go-mygram/models"

func RunMigration() {
	DB.AutoMigrate(&models.User{}, &models.Photo{}, &models.Comment{}, &models.SocialMedia{})
}

