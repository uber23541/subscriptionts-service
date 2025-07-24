// swag init -g server/main.go -d . --output docs

package main

import (
	"log"

	"subscriptions_service/internal/config"
	"subscriptions_service/internal/models/entities"
	"subscriptions_service/internal/router"

	docs "subscriptions_service/docs"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Subscriptions Service API
// @version         1.0
// @description     REST‑сервис для учёта онлайн‑подписок пользователей.
// @BasePath        /
func main() {
	docs.SwaggerInfo.BasePath = "/"
	cfg := config.LoadConfig()

	db, err := gorm.Open(postgres.Open(cfg.DBDsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("db open error: %v", err)
	}

	if err := db.AutoMigrate(&entities.Subscription{}); err != nil {
		log.Fatalf("automigrate error: %v", err)
	}

	r := router.SetupRouter(db)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if err := r.Run(":" + cfg.HTTPPort); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
