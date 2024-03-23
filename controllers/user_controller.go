package controllers

import (
	"net/http"
	"regexp"
	"time"

	"github.com/Mikael88/go-mygram/config"
	"github.com/Mikael88/go-mygram/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UpdateUserRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

type UpdateUserResponse struct {
    ID        uint      `json:"id"`
    Email     string    `json:"email"`
    Username  string    `json:"username"`
    Age       int       `json:"age"`
    UpdatedAt time.Time `json:"updated_at"`
}


// Register
func RegisterUser(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
    if !emailRegex.MatchString(user.Email) {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Format email tidak valid"})
        return
    }

	user.UpdateAt = time.Now()

	if err := config.DB.Create(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
	}

	userResponse := models.UserResponse{
        Age:      user.Age,
        Email:    user.Email,
        ID:       user.ID,
        Username: user.Username,
    }

	c.JSON(http.StatusCreated, gin.H{"data": userResponse})
}
// Login
func LoginUser(c *gin.Context) {
	var input struct {
		Email 	 string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	token, err := generateJWTToken(user.ID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
// Untuk update data user
func UpdateUser(c *gin.Context) {
    userID, exists := c.Get("userId")
    if !exists {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    var user models.User
    if err := config.DB.Where("id = ?", userID).First(&user).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    var req UpdateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    user.Email = req.Email
    if req.Password != "" {
        user.Password = req.Password
        config.DB.Model(&user).Update("password", user.Password)
    }

    if err := config.DB.Save(&user).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
        return
    }

    response := UpdateUserResponse{
        ID:        user.ID,
        Email:     user.Email,
        Username:  user.Username,
        Age:       user.Age,
        UpdatedAt: user.UpdateAt,
    }

    c.JSON(http.StatusOK, response)
}
// Untuk hapus user
func DeleteUser(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var user models.User
	if err := config.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if err := config.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
// Generate token
func generateJWTToken(userId uint) (string, error) {
	claims := jwt.MapClaims{
		"userId": userId,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte("your_secret_key"))
}