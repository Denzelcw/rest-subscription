package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"task_manager/internal/http_server/dto"
	"task_manager/internal/lib/api/resp"
	valid "task_manager/internal/lib/api/valid"
	"task_manager/internal/lib/logger/sl"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type UpdateResponse struct {
	ServiceName string    `json:"service_name"`
	Price       int       `json:"price"`
	UserID      uuid.UUID `json:"user_id"`
	StartDate   string    `json:"start_date"`
	EndDate     string    `json:"end_date"`
}

func (h *SubscriptionHandler) UpdateSubscriptionHandler(w http.ResponseWriter, r *http.Request) {
	const op = "handler.AddSubscriptionHandler"

	log := h.log.With(
		slog.String("op", op),
		slog.String("request_url", middleware.GetReqID(r.Context())),
	)

	var req dto.UpdateSubDTO

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Error("failed to parse id", sl.Err(err))

		resp.Error(w, "invalid subscription ID", http.StatusBadRequest)
		return
	}
	req.ID = id

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error("failed to decode request", sl.Err(err))

		resp.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err = valid.ValidateDates(req.StartDate, req.EndDate)
	if err != nil {
		log.Error("invalid request body", sl.Err(err))

		resp.Error(w, fmt.Sprintf("invalid request body: %s", err), http.StatusBadRequest)
		return
	}

	validWithOpts := validator.New(validator.WithRequiredStructEnabled())
	if err := validWithOpts.Struct(req); err != nil {
		var validateErr validator.ValidationErrors
		errors.As(err, &validateErr)
		log.Error("invalid request", sl.Err(err))

		resp.Error(w, fmt.Sprintf("invalid request: %s", valid.ValidationError(validateErr, req)), http.StatusBadRequest)
		return
	}

	sub, err := h.service.UpdateById(req)
	if err != nil {
		log.Error("failed to update subscription", sl.Err(err))

		resp.Error(w, "failed to update subscription", http.StatusInternalServerError)
		return
	}

	resp.ResponseOk(w, sub, http.StatusCreated)
}
