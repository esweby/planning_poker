package api

import (
	"backend/internals/poker"
	"net/http"

	"github.com/gin-gonic/gin"
)

func handleStartRoom(ctx *gin.Context) {
	var req struct {
		Owner poker.User `json:"owner"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	roomID := roomManager.CreateNewRoom(req.Owner)

	ctx.JSON(http.StatusOK, gin.H{
		"roomNumber": roomID,
	})
}
