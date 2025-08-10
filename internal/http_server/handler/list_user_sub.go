package handler

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"task_manager/internal/lib/api/resp"
	"task_manager/internal/lib/logger/sl"
	"task_manager/internal/storage"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
)

// GetListUserSubscriptionHandler godoc
// @Summary      Get list of user subscriptions
// @Description  Returns a list of a user's subscriptions by their UUID
// @Accept       json
// @Produce      json
// @Param        user_id query     string true "User UUID"
// @Success      200 {array}  domain.UserSubscription
// @Failure      400 {object} resp.ErrorResponse "Invalid UUID or missing parameter"
// @Failure      404 {object} resp.ErrorResponse "User not found"
// @Failure      500 {object} resp.ErrorResponse "Server error"
// @Router       /subscriptions [get]
func (h *UserSubscriptionHandler) GetListUserSubscriptionHandler(w http.ResponseWriter, r *http.Request) {
	const op = "handler.GetUserSubscriptionHandler"

	ctx, cancel := context.WithTimeout(r.Context(), h.timeOut)
	defer cancel()

	log := h.log.With(
		slog.String("op", op),
		slog.String("request_url", middleware.GetReqID(ctx)),
	)

	userIdStr := r.URL.Query().Get("user_id")
	if userIdStr == "" {
		log.Error("user_id is missing in query parameters")
		resp.Error(w, "user_id is required in query parameters", http.StatusBadRequest)
		return
	}

	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		log.Error("failed to parse user_id as UUID", sl.Err(err))
		resp.Error(w, "invalid user_id format (must be a valid UUID)", http.StatusBadRequest)
		return
	}

	subs, err := h.service.GetListByUUID(ctx, userId)
	if err != nil {
		log.Error("failed to get user subscriptions", sl.Err(err))
		if errors.Is(err, storage.ErrUserNotFound) {
			resp.Error(w, "user not found", http.StatusNotFound)
		} else {
			resp.Error(w, "failed to get user subscriptions list", http.StatusInternalServerError)
		}
		return
	}

	resp.ResponseOk(w, subs, http.StatusOK)
}
