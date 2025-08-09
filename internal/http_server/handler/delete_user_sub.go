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

type DeleteResponse struct {
	Id      int    `json:"id"`
	Message string `json:"message"`
}

func (h *SubscriptionHandler) DeleteSubscriptionHandler(w http.ResponseWriter, r *http.Request) {
	const op = "handler.DeleteSubscriptionHandler"

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

	err = h.service.DeleteById(id)
	if err != nil {
		log.Error("failed to delete subscription", sl.Err(err))
		if errors.Is(err, storage.ErrNotFound) {
			resp.Error(w, "subscription not found", http.StatusNotFound)
		} else {
			resp.Error(w, "failed to delete subscription", http.StatusInternalServerError)
		}
		return
	}

	response := DeleteResponse{
		Id:      id,
		Message: "user subscription successfully deleted",
	}

	resp.ResponseOk(w, response, http.StatusOK)
}
