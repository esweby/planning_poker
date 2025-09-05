package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func handleCheckRoom(ctx *gin.Context) {
	roomID := ctx.Param("roomId")

	if len(roomID) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "no room id provided"})
		return
	}

	_, err := roomManager.GetRoom(roomID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusFound, gin.H{"success": true})
}
