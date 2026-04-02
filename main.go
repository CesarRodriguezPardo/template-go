package main

import (
	"CesarRodriguezPardo/template-go/config"
	"CesarRodriguezPardo/template-go/docs"
	mailer "CesarRodriguezPardo/template-go/infra/mailer"
	"CesarRodriguezPardo/template-go/internal/middleware"
	"CesarRodriguezPardo/template-go/internal/repositories"
	"CesarRodriguezPardo/template-go/internal/routes"
	"CesarRodriguezPardo/template-go/internal/services"
	"context"
	"fmt"
	"log"

	"CesarRodriguezPardo/template-go/infra/database"
	logger "CesarRodriguezPardo/template-go/infra/logger"

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
	if err := app.Run(":" + port); err != nil {
		return fmt.Errorf("Could not initialize app: %w", err)
	}

	return nil
}

func main() {
	if err := config.LoadConfig(); err != nil {
		log.Fatal(err)
	}

	logger.InitLogger()

	repositories.InitConnections()

	ctx := context.Background()
	db, err := database.Connect(ctx)

	if err != nil {
		logger.Fatal("could not connect to postgres", err)
	}

	services.InitRepositories(db)

	mailer.InitMailer()

	if err := setupApp(); err != nil {
		logger.Fatal("could not setup app", err)
	}
}
