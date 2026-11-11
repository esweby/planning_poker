package api

import (
	"backend/internals/poker"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const closeWriteWait = time.Second

func handleJoinRoom(ctx *gin.Context, upgrader *websocket.Upgrader) {
	roomID := ctx.Param("roomId")
	token := ctx.Query("token")

	user := poker.User{
		Username: poker.Username(ctx.Query("username")),
		Role:     poker.Role(ctx.Query("role")),
		Seed:     ctx.Query("seed"),
	}

	// Reject what we can before upgrading, so the client gets a plain HTTP
	// error instead of a websocket that opens and immediately closes.
	room, err := roomManager.GetRoom(roomID)
	if err == nil {
		err = room.CheckJoin(user.Username, token)
	}
	if err != nil {
		ctx.JSON(joinErrorStatus(err), gin.H{"error": err.Error()})
		return
	}

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Printf("failed websocket upgrade: %v", err)
		return
	}

	// The checks above can race with another join or the room closing, so
	// JoinRoom validates again.
	if _, err := roomManager.JoinRoom(roomID, user, token, conn); err != nil {
		log.Printf("failed to join room: %v", err)
		closeMsg := websocket.FormatCloseMessage(websocket.ClosePolicyViolation, err.Error())
		conn.WriteControl(websocket.CloseMessage, closeMsg, time.Now().Add(closeWriteWait))
		conn.Close()
	}
}
