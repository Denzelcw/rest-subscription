package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"task_manager/internal/http_server/dto"
	"task_manager/internal/lib/api/resp"
	valid "task_manager/internal/lib/api/valid"
	"task_manager/internal/lib/logger/sl"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
)

type TotalCostResponse struct {
	TotalCost int64 `json:"total_cost"`
}

func (h *SubscriptionHandler) GetTotalCostHandler(w http.ResponseWriter, r *http.Request) {
	const op = "handler.GetTotalCostHandler"

	log := h.log.With(
		slog.String("op", op),
		slog.String("request_url", middleware.GetReqID(r.Context())),
	)

	var req dto.TotalCost

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error("failed to decode request", sl.Err(err))

		resp.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	err := valid.ValidateDates(req.StartDate, req.EndDate)
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

	totalCost, err := h.service.TotalCost(req)
	if err != nil {
		log.Error("failed to get total cost", sl.Err(err))
		resp.Error(w, "subscription not found", http.StatusInternalServerError)

		return
	}

	response := TotalCostResponse{TotalCost: totalCost}

	resp.ResponseOk(w, response, http.StatusOK)
}
