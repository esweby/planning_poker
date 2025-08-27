package api

import (
	"backend/internals/poker"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	roomManager = poker.GetRoomManager()
)

func Start() {
	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	config.ExposeHeaders = []string{"Authorization"}

	router.Use(cors.New(config))

	api := router.Group("/api")

	api.POST("/", handleStartRoom)
	api.POST("/join", handleJoinRoom)

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("%s", err)
	}
}
