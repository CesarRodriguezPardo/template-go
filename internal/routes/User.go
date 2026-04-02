package routes

import (
	"CesarRodriguezPardo/template-go/internal/controllers"

	"github.com/gin-gonic/gin"
)

func InitUserRoutes(r *gin.RouterGroup) {
	userGroup := r.Group("/user")
	{
		userGroup.POST("/", controllers.CreateUser)
		//userGroup.GET("/postgres", middleware.SetRoles(models.ALL), middleware.LoadJWTAuth().MiddlewareFunc(), controllers.GetAllUsersControllerPostgres)
	}
}
