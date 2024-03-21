package controllers

import (
	"net/http"

	"github.com/Mikael88/go-mygram/config"
	"github.com/Mikael88/go-mygram/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

type CreatePhotoInput struct {
	Title    string `json:"title" binding:"required"`
	Caption  string `json:"caption"`
	PhotoURL string `json:"photo_url" binding:"required"`
}

func init() {
	// Buat objek validator
	validate = validator.New()
}
// Create menambahkan foto baru
func CreatePhoto(c *gin.Context) {
	var input CreatePhotoInput

	// Bind request body ke struct input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validasi input menggunakan objek validator
	if err := validate.Struct(input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Dapatkan ID pengguna dari konteks
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Buat objek foto
	photo := models.Photo{
		Title:    input.Title,
		Caption:  input.Caption,
		PhotoURL: input.PhotoURL,
		UserID:   userId.(uint), // Konversi userId menjadi uint
	}

	// Simpan foto ke database
	if err := config.DB.Create(&photo).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": photo})
}
// UpdatePhoto mengelola proses pembaruan informasi foto.
func UpdatePhoto(c *gin.Context) {
	photoID := c.Param("photoId")

	var photo models.Photo
	if err := config.DB.Where("id = ?", photoID).First(&photo).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Photo not found"})
		return
	}

	// Periksa apakah pengguna yang meminta pembaruan adalah pemilik foto
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if photo.UserID != userId.(uint) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var updatePhoto models.Photo
	if err := c.ShouldBindJSON(&updatePhoto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	photo.Title = updatePhoto.Title
	photo.Caption = updatePhoto.Caption
	photo.PhotoURL = updatePhoto.PhotoURL

	if err := config.DB.Save(&photo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update photo"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": photo})
}
// DeletePhoto mengelola proses penghapusan foto.
func DeletePhoto(c *gin.Context) {
	photoID := c.Param("photoId")

	var photo models.Photo
	if err := config.DB.Where("id = ?", photoID).First(&photo).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Photo not found"})
		return
	}

	// Periksa apakah pengguna yang meminta penghapusan adalah pemilik foto
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if photo.UserID != userId.(uint) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if err := config.DB.Delete(&photo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete photo"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Photo deleted successfully"})
}
