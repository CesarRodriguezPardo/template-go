package routes

import (
	"github.com/gin-gonic/gin"
)

func InitRoutes(r *gin.Engine) {
	group := r.Group("/api/v1")
	InitAuthRoutes(group)
	InitUserRoutes(group)
}
