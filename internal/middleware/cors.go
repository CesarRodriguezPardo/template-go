package middleware

import (
	"CesarRodriguezPardo/template-go/config"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CorsMiddleware() gin.HandlerFunc {
	allowedCors := config.Cfg.Cors.AllowedCors

	configCors := cors.DefaultConfig()

	configCors.AllowMethods = append(configCors.AllowMethods, "DELETE", "OPTIONS", "POST", "GET", "PUT")
	configCors.AllowHeaders = append(configCors.AllowHeaders, "Authorization", "Pagination-Count")
	configCors.ExposeHeaders = append(configCors.ExposeHeaders, "Pagination-Count")
	configCors.AllowOrigins = strings.Split(allowedCors, ",")
	configCors.AllowCredentials = true

	return cors.New(configCors)
}
