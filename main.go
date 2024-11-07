package main

import (
	"github.com/Guesstrain/EthBankok/config"
	"github.com/Guesstrain/EthBankok/controllers"
	"github.com/Guesstrain/EthBankok/services"
	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDatabase()
	merchantService := services.NewMerchantService()

	router := gin.New()
	router.Use(gin.Recovery(), gin.Logger())

	// Define routes for merchants
	router.POST("/merchants", func(c *gin.Context) {
		controllers.AddMerchantHandler(c, merchantService)
	})
	router.GET("/merchants/:id", func(c *gin.Context) {
		controllers.GetMerchantByIDHandler(c, merchantService)
	})
	router.GET("/merchants", func(c *gin.Context) {
		controllers.GetAllMerchantsHandler(c, merchantService)
	})

	router.Run(":8080")
}
