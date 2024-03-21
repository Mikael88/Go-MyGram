package middlewares

import (
	"net/http"

	"github.com/Mikael88/go-mygram/config"
	"github.com/Mikael88/go-mygram/models"
	"github.com/gin-gonic/gin"
)

func AuthorizeSocialMedia() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Mendapatkan ID pengguna dari context
		userId, exists := c.Get("userId")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Mendapatkan ID media sosial dari path parameter
		socialMediaId := c.Param("socialMediaId")

		// Mencari media sosial berdasarkan ID
		var socialMedia models.SocialMedia
		if err := config.DB.Where("id = ?", socialMediaId).First(&socialMedia).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Social media not found"})
			c.Abort()
			return
		}

		// Memeriksa apakah pengguna memiliki izin untuk mengubah atau menghapus media sosial
		if socialMedia.UserID != userId {
			c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to perform this action"})
			c.Abort()
			return
		}

		c.Next()
	}
}