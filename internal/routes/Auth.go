package routes

import (
	"CesarRodriguezPardo/template-go/internal/controllers"
	"CesarRodriguezPardo/template-go/internal/middleware"

	"github.com/gin-gonic/gin"
)

func InitAuthRoutes(r *gin.RouterGroup) {
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/login", controllers.LoginFunc)
		authGroup.POST("/refresh", middleware.GetJWTAuth().MiddlewareFunc(), controllers.RefreshToken)
		authGroup.POST("/logout", middleware.GetJWTAuth().MiddlewareFunc(), controllers.Logout)
	}
}
