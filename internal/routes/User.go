package routes

import (
	"citiaps/golang-backend-template/controllers"
	"citiaps/golang-backend-template/middleware"
	"citiaps/golang-backend-template/models"

	"github.com/gin-gonic/gin"
)

func InitUserRoutes(r *gin.RouterGroup) {
	userGroup := r.Group("/user")
	{
		// postgres
		userGroup.POST("/postgres", controllers.CreateUserControllerPostgres)
		userGroup.GET("/postgres", middleware.SetRoles(models.ALL), middleware.LoadJWTAuth().MiddlewareFunc(), controllers.GetAllUsersControllerPostgres)
	}
}
