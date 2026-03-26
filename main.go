package main

import (
	"CesarRodriguezPardo/template-go/config"
	"citiaps/golang-backend-template/docs"
	"citiaps/golang-backend-template/mailer"
	"citiaps/golang-backend-template/middleware"
	"citiaps/golang-backend-template/repositories"
	"citiaps/golang-backend-template/routes"
	"citiaps/golang-backend-template/services"
	"citiaps/golang-backend-template/utils"
	"fmt"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func initApp() *gin.Engine {
	app := gin.Default()

	return app
}

func initSwagger(app *gin.Engine) {
	docs.SwaggerInfo.BasePath = "/api/v1"
	app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func initCorsMiddleware(app *gin.Engine) {
	app.Use(middleware.CorsMiddleware())
}

func setupApp() error {
	mode := config.Cfg.Server.Mode
	gin.SetMode(mode)

	app := initApp()

	initCorsMiddleware(app)

	initSwagger(app)

	routes.InitRoutes(app)

	port := config.Cfg.Server.Port
	err := app.Run(":" + port)

	if err != nil {
		return fmt.Errorf("Could not initialize app: %w.", err)
	}

	return nil
}

func main() {
	config.LoadConfig()

	utils.InitLogger()

	repositories.InitConnections()

	services.InitRepositories()

	services.InitIndexes()

	mailer.InitMailer()

	err := setupApp()

	if err != nil {
		utils.Fatal(err)
	}
}
