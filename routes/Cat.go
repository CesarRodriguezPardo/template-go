package routes

import (
	"citiaps/golang-backend-template/controllers"
	"citiaps/golang-backend-template/middleware"
	"citiaps/golang-backend-template/models"

	"github.com/gin-gonic/gin"
)

func InitCatRoutes(r *gin.RouterGroup) {
	catGroup := r.Group("/cat")
	{
		catGroup.POST("/", middleware.LoadJWTAuth().MiddlewareFunc(), controllers.CreateCatController)
		catGroup.GET("/", middleware.SetRoles(models.ADMIN), middleware.LoadJWTAuth().MiddlewareFunc(), controllers.GetAllCatsController)

		// postgres
		catGroup.POST("/postgres", controllers.CreateCatControllerPostgres)
		catGroup.GET("/postgres", controllers.GetAllCatsControllerPostgres)
	}
}
