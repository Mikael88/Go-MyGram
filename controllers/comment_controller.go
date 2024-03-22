package controllers

import (
	"net/http"
	"time"

	"github.com/Mikael88/go-mygram/config"
	"github.com/Mikael88/go-mygram/models"

	"github.com/gin-gonic/gin"
)

// Create Comment validation
type CreateCommentInput struct {
	Message string `json:"message" binding:"required"`
	PhotoID uint   `json:"photo_id" binding:"required"`
}
func CreateComment(c *gin.Context) {
	var input CreateCommentInput

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

	// Buat objek komentar
	comment := models.Comment{
		Message: input.Message,
		PhotoID: input.PhotoID,
		UserID:  userId.(uint), // Konversi userId menjadi uint
	}

	// Simpan komentar ke database
	if err := config.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response := map[string]interface{}{
		"id":        comment.ID,
		"message":   comment.Message,
		"photo_id":  comment.PhotoID,
		"user_id":   comment.UserID,
		"created_at": comment.CreatedAt.Format(time.RFC3339), // Format date-time
	}

	c.JSON(http.StatusCreated, gin.H{"data": response})
}
// GetComments mengambil daftar komentar
func GetComments(c *gin.Context) {
    var comments []models.Comment

    // Ambil semua komentar dari database
    if err := config.DB.Preload("User").Preload("Photo.User").Find(&comments).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Jika tidak ada komentar ditemukan, kembalikan respons kosong
    if len(comments) == 0 {
        c.JSON(http.StatusOK, gin.H{"message": "No comments found"})
        return
    }

    // Transformasi data komentar ke format yang diinginkan
    formattedComments := make([]map[string]interface{}, len(comments))
    for i, comment := range comments {
        formattedComment := map[string]interface{}{
            "id":         comment.ID,
            "message":    comment.Message,
            "photo_id":   comment.PhotoID,
            "user_id":    comment.UserID,
            "updated_at": comment.UpdatedAt,
            "created_at": comment.CreatedAt,
            "User": map[string]interface{}{
                "id":       comment.User.ID,
                "email":    comment.User.Email,
                "username": comment.User.Username,
            },
            "Photo": map[string]interface{}{
                "id":        comment.Photo.ID,
                "title":     comment.Photo.Title,
                "caption":   comment.Photo.Caption,
                "photo_url": comment.Photo.PhotoURL,
                "user_id":   comment.Photo.User.ID,
            },
        }
        formattedComments[i] = formattedComment
    }

    // Kembalikan daftar komentar dalam format yang diinginkan
    c.JSON(http.StatusOK, formattedComments)
}
// UpdateComment mengelola proses pembaruan komentar.
func UpdateComment(c *gin.Context) {
    commentID := c.Param("commentId")
    var comment models.Comment
    if err := config.DB.Where("id = ?", commentID).Preload("Photo").Preload("Photo.User").First(&comment).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
        return
    }

    // Periksa apakah pengguna yang meminta pembaruan adalah pemilik komentar
    userId, exists := c.Get("userId")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }
    if comment.UserID != userId.(uint) {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    var updateComment models.Comment
    if err := c.ShouldBindJSON(&updateComment); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    comment.Message = updateComment.Message
    if err := config.DB.Save(&comment).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update comment"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "id":         comment.Photo.ID,
        "title":      comment.Photo.Title,
        "caption":    comment.Photo.Caption,
        "photo_url":  comment.Photo.PhotoURL,
        "user_id":    comment.Photo.User.ID,
        "updated_at": comment.Photo.UpdatedAt,
    })
}
// DeleteComment mengelola proses penghapusan komentar.
func DeleteComment(c *gin.Context) {
	commentID := c.Param("commentId")

	var comment models.Comment
	if err := config.DB.Where("id = ?", commentID).First(&comment).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	// Periksa apakah pengguna yang meminta penghapusan adalah pemilik komentar
	userId, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if comment.UserID != userId.(uint) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if err := config.DB.Delete(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete comment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
}
