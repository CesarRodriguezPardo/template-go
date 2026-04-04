package routes

import (
	"CesarRodriguezPardo/template-go/internal/controllers"
	"CesarRodriguezPardo/template-go/internal/middleware"
	"CesarRodriguezPardo/template-go/internal/models"

	"github.com/gin-gonic/gin"
)

func InitUserRoutes(r *gin.RouterGroup) {
	userGroup := r.Group("/user")
	{
		userGroup.POST("/", controllers.CreateUser)
		userGroup.GET("/", middleware.SetRoles(models.ADMIN, models.WORKER), middleware.GetJWTAuth().MiddlewareFunc(), controllers.GetAllUsers)
	}
}
