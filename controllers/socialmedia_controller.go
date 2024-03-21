package controllers

import (
	"net/http"

	"github.com/Mikael88/go-mygram/config"
	"github.com/Mikael88/go-mygram/models"

	"github.com/gin-gonic/gin"
)

// CreateSocialMediaInput adalah struktur untuk validasi input saat membuat data sosial media
type CreateSocialMediaInput struct {
	Name           string `json:"name" binding:"required"`
	SocialMediaURL string `json:"social_media_url" binding:"required"`
}

// CreateSocialMedia menambahkan data sosial media baru
func CreateSocialMedia(c *gin.Context) {
	var input CreateSocialMediaInput

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

	// Buat objek data sosial media
	socialMedia := models.SocialMedia{
		Name:           input.Name,
		SocialMediaURL: input.SocialMediaURL,
		UserID:         userId.(uint), // Konversi userId menjadi uint
	}

	// Simpan data sosial media ke database
	if err := config.DB.Create(&socialMedia).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": socialMedia})
}
// GetSocialMedias mengambil daftar media sosial
func GetSocialMedias(c *gin.Context) {
    var socialMedias []models.SocialMedia

    // Ambil semua media sosial dari database
    if err := config.DB.Find(&socialMedias).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Kembalikan daftar media sosial dalam respons
    c.JSON(http.StatusOK, gin.H{"social_medias": socialMedias})
}
// UpdateSocialMediaInput adalah struktur untuk validasi input saat memperbarui data sosial media
type UpdateSocialMediaInput struct {
	Name           string `json:"name" binding:"required"`
	SocialMediaURL string `json:"social_media_url" binding:"required"`
}
// UpdateSocialMedia mengelola proses pembaruan informasi data sosial media
func UpdateSocialMedia(c *gin.Context) {
	socialMediaID := c.Param("socialMediaId")

	var socialMedia models.SocialMedia
	if err := config.DB.Where("id = ?", socialMediaID).First(&socialMedia).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Social media not found"})
		return
	}

	// Periksa apakah pengguna yang meminta pembaruan adalah pemilik data sosial media
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if socialMedia.UserID != userId.(uint) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var updateInput UpdateSocialMediaInput
	if err := c.ShouldBindJSON(&updateInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validasi input menggunakan objek validator
	if err := validate.Struct(updateInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	socialMedia.Name = updateInput.Name
	socialMedia.SocialMediaURL = updateInput.SocialMediaURL

	if err := config.DB.Save(&socialMedia).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update social media"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": socialMedia})
}
// DeleteSocialMedia mengelola proses penghapusan data sosial media
func DeleteSocialMedia(c *gin.Context) {
	socialMediaID := c.Param("socialMediaId")

	var socialMedia models.SocialMedia
	if err := config.DB.Where("id = ?", socialMediaID).First(&socialMedia).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Social media not found"})
		return
	}

	// Periksa apakah pengguna yang meminta penghapusan adalah pemilik data sosial media
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if socialMedia.UserID != userId.(uint) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if err := config.DB.Delete(&socialMedia).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete social media"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Social media deleted successfully"})
}
