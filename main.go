package main

import (
	"citiaps/golang-backend-template/config"
	"citiaps/golang-backend-template/docs"
	"citiaps/golang-backend-template/mailer"
	"citiaps/golang-backend-template/middleware"
	"citiaps/golang-backend-template/repositories"
	"citiaps/golang-backend-template/routes"
	"citiaps/golang-backend-template/services"
	"citiaps/golang-backend-template/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	// config
	config.LoadConfig()

	// logger
	utils.InitLogger(config.Cfg.Logger.Filepath, config.Cfg.Logger.Filename, config.Cfg.Logger.Level, config.Cfg.Logger.Tz)

	// incializacion de conexiones
	repositories.InitConnections()

	// incializacion de repositorios
	services.InitRepositories()

	// inicializacion de indices
	services.InitIndexes()

	// mailer
	mailer.InitMailer()

	// logs
	utils.Info("server up en: " + config.Cfg.Server.Port)

	// gin 
	app := gin.Default()
	routes.InitRoutes(app)
	app.Use(middleware.CorsMiddleware(config.Cfg.Cors.CorsUrl))

	// swagger
	docs.SwaggerInfo.BasePath = "/api/v1"
	app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	http.ListenAndServe(":"+config.Cfg.Server.Port, app)
}
