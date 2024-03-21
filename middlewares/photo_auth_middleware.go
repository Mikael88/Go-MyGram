package middlewares

import (
	"net/http"

	"github.com/Mikael88/go-mygram/config"
	"github.com/Mikael88/go-mygram/models"
	"github.com/gin-gonic/gin"
)

func AuthorizePhoto() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Mendapatkan ID pengguna dari context
		userId, exists := c.Get("userId")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Mendapatkan ID foto dari path parameter
		photoId := c.Param("photoId")

		// Mencari foto berdasarkan ID
		var photo models.Photo
		if err := config.DB.Where("id = ?", photoId).First(&photo).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Photo not found"})
			c.Abort()
			return
		}

		// Memeriksa apakah pengguna memiliki izin untuk mengubah atau menghapus foto
		if photo.UserID != userId {
			c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to perform this action"})
			c.Abort()
			return
		}

		c.Next()
	}
}
