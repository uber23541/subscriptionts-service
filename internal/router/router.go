package router

import (
	"subscriptions_service/internal/handlers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	h := handlers.NewHandler(db)

	r.POST("/subscriptions", h.CreateSubscription)
	r.GET("/subscriptions", h.ListSubscriptions)

	r.GET("/subscriptions/:user_id/:service_name", h.GetSubscription)
	r.PUT("/subscriptions/:user_id/:service_name", h.UpdateSubscription)
	r.DELETE("/subscriptions/:user_id/:service_name", h.DeleteSubscription)

	r.GET("/subscriptions/summary", h.SumSubscriptions)

	return r
}
