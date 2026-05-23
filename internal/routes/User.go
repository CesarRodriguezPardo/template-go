package routes

import (
	"CesarRodriguezPardo/template-go/internal/controllers"
	"CesarRodriguezPardo/template-go/internal/middleware"
	"CesarRodriguezPardo/template-go/internal/models"

	"github.com/gin-gonic/gin"
)

func InitUserRoutes(r *gin.RouterGroup) {
	userGroup := r.Group("/user")

	adminOnly := []models.Role{models.ADMIN}
	allRoles := []models.Role{models.ADMIN, models.WORKER}

	userGroup.POST("/", middleware.SetRoles(adminOnly...), middleware.GetJWTAuth().MiddlewareFunc(), controllers.CreateUser)
	userGroup.GET("/", middleware.SetRoles(allRoles...), middleware.GetJWTAuth().MiddlewareFunc(), controllers.GetAllUsers)
	userGroup.GET("/:id", middleware.GetJWTAuth().MiddlewareFunc(), controllers.GetUserByID)
	userGroup.PUT("/:id", middleware.GetJWTAuth().MiddlewareFunc(), controllers.UpdateUser)
	userGroup.DELETE("/:id", middleware.GetJWTAuth().MiddlewareFunc(), controllers.DeleteUser)
}
