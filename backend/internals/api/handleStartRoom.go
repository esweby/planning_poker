package api

import (
	"backend/internals/poker"
	"log"
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

	roomID, token, err := roomManager.CreateNewRoom(req.Owner)
	if err != nil {
		ctx.JSON(joinErrorStatus(err), gin.H{"error": err.Error()})
		return
	}

	log.Println("Room started by ", req.Owner.Username)

	ctx.JSON(http.StatusOK, gin.H{
		"roomNumber": roomID,
		"token":      token,
	})
}
