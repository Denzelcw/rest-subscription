package handler

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"task_manager/internal/lib/api/resp"
	"task_manager/internal/lib/logger/sl"
	"task_manager/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (h *SubscriptionHandler) GetSubscriptionHandler(w http.ResponseWriter, r *http.Request) {
	const op = "handler.GetSubscriptionHandler"

	log := h.log.With(
		slog.String("op", op),
		slog.String("request_url", middleware.GetReqID(r.Context())),
	)

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Error("failed to parse id", sl.Err(err))

		resp.Error(w, "invalid subscription ID", http.StatusBadRequest)
		return
	}

	subscription, err := h.service.GetById(id)
	if err != nil {
		log.Error("failed to get subscription", sl.Err(err))
		if errors.Is(err, storage.ErrNotFound) {
			resp.Error(w, "subscription not found", http.StatusNotFound)
		} else {
			resp.Error(w, "failed to get subscription", http.StatusInternalServerError)
		}
		return
	}

	resp.ResponseOk(w, subscription, http.StatusOK)
}
