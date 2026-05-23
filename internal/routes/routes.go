package routes

import (
	"CesarRodriguezPardo/template-go/internal/middleware"

	"github.com/gin-gonic/gin"
)

func InitRoutes(r *gin.Engine) {
	group := r.Group("/api/v1")
	group.Use(middleware.RateLimitMiddleware(100, 150))
	InitAuthRoutes(group)
	InitUserRoutes(group)
}
