package middleware

import (
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CorsMiddleware(corsUrl string) gin.HandlerFunc {
	configCors := cors.DefaultConfig()

	configCors.AllowMethods = append(configCors.AllowMethods, "DELETE", "OPTIONS", "POST", "GET", "PUT")
	configCors.AllowHeaders = append(configCors.AllowHeaders, "Authorization", "Pagination-Count")
	configCors.ExposeHeaders = append(configCors.ExposeHeaders, "Pagination-Count")
	configCors.AllowOrigins = strings.Split(corsUrl, ",")
	//config.AllowAllOrigins = true
	configCors.AllowCredentials = false

	return cors.New(configCors)
}
