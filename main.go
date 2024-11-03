package main

import (
	"github.com/gin-gonic/gin"
)

const keyDir = "./keystore"

func main() {
	router := gin.New()
	router.Use(gin.Recovery(), gin.Logger())

	router.Run(":8080")
}
