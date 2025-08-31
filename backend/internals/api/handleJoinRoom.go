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
	var req struct {
		User   poker.User   `json:"user"`
		RoomID poker.RoomID `json:"roomId"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	if !roomManager.RoomExists(req.RoomID) {
        ctx.JSON(http.StatusNotFound, gin.H{"error": "Room not found"})
        return
    }

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upgrade connection"})
		return
	}

	roomManager.AddUserToRoom(req.RoomID, req.User, conn)
}
