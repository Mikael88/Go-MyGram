package routes

import (
	"github.com/Mikael88/go-mygram/controllers"
	"github.com/Mikael88/go-mygram/middlewares"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.POST("/register", controllers.RegisterUser)
	r.POST("/login", controllers.LoginUser)

	api := r.Group("/api")
	api.Use(middlewares.AuthMiddleware()) // Terapkan middleware AuthMiddleware pada grup api

	api.POST("/photos", controllers.CreatePhoto)
	api.PUT("/photos/:photoId", middlewares.AuthorizePhoto(), controllers.UpdatePhoto)
	api.DELETE("/photos/:photoId", middlewares.AuthorizePhoto(), controllers.DeletePhoto)

	api.POST("/comments", controllers.CreateComment)
	api.PUT("/comments/:commentId", middlewares.AuthorizeComment(), controllers.UpdateComment)
	api.DELETE("/comments/:commentId", middlewares.AuthorizeComment(), controllers.DeleteComment)

	api.POST("/socialmedias", controllers.CreateSocialMedia)
	api.PUT("/socialmedias/:socialMediaId", middlewares.AuthorizeSocialMedia(), controllers.UpdateSocialMedia)
	api.DELETE("/socialmedias/:socialMediaId", middlewares.AuthorizeSocialMedia(), controllers.DeleteSocialMedia)
}