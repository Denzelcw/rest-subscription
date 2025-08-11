package er

import (
	"errors"
	"fmt"
	"net/http"
	"subscription/internal/storage"
)

func MapErrorToStatus(err error) (string, int, bool) {
	fmt.Println("ERROR", err)
	switch {
	case errors.Is(err, storage.ErrNotFound):
		return "user_subscription not found", http.StatusNotFound, true
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
