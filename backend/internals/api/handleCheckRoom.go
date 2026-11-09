package api

import (
	"backend/internals/poker"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// handleCheckRoom reports whether the room exists and, when a username is
// given, whether that name (with the optional token) may join it.
func handleCheckRoom(ctx *gin.Context) {
	roomID := ctx.Param("roomId")

	room, err := roomManager.GetRoom(roomID)
	if err != nil {
		log.Printf("- handleCheckRoom error: %s", err.Error())
		ctx.JSON(joinErrorStatus(err), gin.H{"error": err.Error()})
		return
	}

	if username, ok := ctx.GetQuery("username"); ok {
		if err := room.CheckJoin(poker.Username(username), ctx.Query("token")); err != nil {
			ctx.JSON(joinErrorStatus(err), gin.H{"error": err.Error()})
			return
		}
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true, "roomID": room.ID})
}
