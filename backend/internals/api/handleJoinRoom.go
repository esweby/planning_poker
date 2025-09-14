package api

import (
	"backend/internals/poker"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	
	// IMPORTANT: in production you should check origin properly
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleJoinRoom(ctx *gin.Context) {
	roomId := ctx.Param("roomId")
	username := ctx.Query("username")
	role := ctx.Query("role")
	seed := ctx.Query("seed")

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upgrade connection"})
		return
	}

	user := poker.User{
		Username: poker.Username(username),
		Role: poker.Role(role),
		Seed: seed,
	}

	roomManager.AddUserToRoom(roomId, user, conn)
}
