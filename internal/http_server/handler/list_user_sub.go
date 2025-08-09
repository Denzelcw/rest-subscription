package handler

import (
	"errors"
	"log/slog"
	"net/http"
	"task_manager/internal/lib/api/resp"
	"task_manager/internal/lib/logger/sl"
	"task_manager/internal/storage"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
)

func (h *SubscriptionHandler) GetListSubscriptionHandler(w http.ResponseWriter, r *http.Request) {
	const op = "handler.GetSubscriptionHandler"

	log := h.log.With(
		slog.String("op", op),
		slog.String("request_url", middleware.GetReqID(r.Context())),
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

	subs, err := h.service.GetListByUUID(userId)
	if err != nil {
		log.Error("failed to get subscription", sl.Err(err))
		if errors.Is(err, storage.ErrUserNotFound) {
			resp.Error(w, "user not found", http.StatusNotFound)
		} else {
			resp.Error(w, "failed to get subscriptions list", http.StatusInternalServerError)
		}
		return
	}

	resp.ResponseOk(w, subs, http.StatusOK)
}
