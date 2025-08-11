package er

import (
	"errors"
	"net/http"
	"task_manager/internal/storage"
)

func MapErrorToStatus(err error) (string, int, bool) {
	switch {
	case errors.Is(err, storage.ErrUserNotFound):
		return "user not found", http.StatusNotFound, true
	case errors.Is(err, storage.ErrUserSubExists):
		return "user subscription already exists", http.StatusConflict, true
	case errors.Is(err, storage.ErrOverlap):
		return "user subscription conflicts with existing record", http.StatusConflict, true
	default:
		return "", 0, false
	}
}
