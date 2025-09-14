package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func handleCheckRoom(ctx *gin.Context) {
	roomID := ctx.Param("roomId")

	log.Print("- handleCheckRoom invoked")
	if len(roomID) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "no room id provided"})
		return
	}

	room, err := roomManager.GetRoom(roomID)
	if err != nil {
		log.Printf("- handleCheckRoom error: %s", err.Error())
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	log.Printf("- handleCheckRoom: Room found")
	ctx.JSON(http.StatusOK, gin.H{"success": true, "roomID": room.ID})
}
