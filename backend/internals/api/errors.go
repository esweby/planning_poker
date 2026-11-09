package api

import (
	"backend/internals/poker"
	"errors"
	"net/http"
)

// joinErrorStatus maps a join failure to an HTTP status code.
func joinErrorStatus(err error) int {
	switch {
	case errors.Is(err, poker.ErrRoomNotFound), errors.Is(err, poker.ErrRoomClosed):
		return http.StatusNotFound
	case errors.Is(err, poker.ErrNameTaken):
		return http.StatusConflict
	case errors.Is(err, poker.ErrInvalidUsername):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
