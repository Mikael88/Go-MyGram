package middlewares

import (
	"net/http"

	"github.com/Mikael88/go-mygram/config"
	"github.com/Mikael88/go-mygram/models"
	"github.com/gin-gonic/gin"
)

func AuthorizeComment() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Mendapatkan ID pengguna dari context
		userId, exists := c.Get("userId")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		// Mendapatkan ID komentar dari path parameter
		commentId := c.Param("commentId")

		// Mencari komentar berdasarkan ID
		var comment models.Comment
		if err := config.DB.Where("id = ?", commentId).First(&comment).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
			c.Abort()
			return
		}

		// Memeriksa apakah pengguna memiliki izin untuk mengubah atau menghapus komentar
		if comment.UserID != userId {
			c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to perform this action"})
			c.Abort()
			return
		}

		c.Next()
	}
}