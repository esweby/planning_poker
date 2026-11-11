package api

import (
	"backend/internals/poker"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	roomManager = poker.GetRoomManager()
)

// defaultAllowedOrigins covers the Vite dev server. Set ALLOWED_ORIGINS to a
// comma-separated list to override it in other environments.
const defaultAllowedOrigins = "http://localhost:5173,http://127.0.0.1:5173"

func allowedOrigins() []string {
	raw := os.Getenv("ALLOWED_ORIGINS")
	if raw == "" {
		raw = defaultAllowedOrigins
	}

	origins := make([]string, 0)
	for _, origin := range strings.Split(raw, ",") {
		if origin = strings.TrimSpace(origin); origin != "" {
			origins = append(origins, origin)
		}
	}
	return origins
}

// newUpgrader only accepts websocket connections from browsers on an allowed
// origin. Requests with no Origin header (non-browser clients) are allowed,
// since origin checks only protect against other sites in a browser.
func newUpgrader(origins []string) *websocket.Upgrader {
	allowed := make(map[string]bool, len(origins))
	for _, origin := range origins {
		allowed[origin] = true
	}

	return &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			origin := r.Header.Get("Origin")
			return origin == "" || allowed[origin]
		},
	}
}

// NewRouter builds the HTTP API.
func NewRouter() *gin.Engine {
	origins := allowedOrigins()
	upgrader := newUpgrader(origins)

	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = origins
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	config.ExposeHeaders = []string{"Authorization"}

	router.Use(cors.New(config))

	api := router.Group("/api")

	api.POST("/", handleStartRoom)
	api.GET("/check/:roomId", handleCheckRoom)
	api.GET("/join/:roomId", func(ctx *gin.Context) { handleJoinRoom(ctx, upgrader) })

	return router
}

func Start() {
	if err := NewRouter().Run(":8080"); err != nil {
		log.Fatalf("%s", err)
	}
}
