package routes

import (
	"citiaps/golang-backend-template/controllers"
	"citiaps/golang-backend-template/middleware"

	"github.com/gin-gonic/gin"
)

func InitAuthRoutes(r *gin.RouterGroup) {
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/login", controllers.LoginFunc)
		authGroup.POST("/refresh", middleware.LoadJWTAuth().MiddlewareFunc(), controllers.RefreshToken)
		authGroup.POST("/logout", middleware.LoadJWTAuth().MiddlewareFunc(), controllers.Logout)
	}
}
