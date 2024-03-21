package controllers

import (
	"net/http"

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

	c.JSON(http.StatusCreated, gin.H{"data": comment})
}
// UpdateComment mengelola proses pembaruan komentar.
func UpdateComment(c *gin.Context) {
	commentID := c.Param("commentId")

	var comment models.Comment
	if err := config.DB.Where("id = ?", commentID).First(&comment).Error; err != nil {
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

	c.JSON(http.StatusOK, gin.H{"data": comment})
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
