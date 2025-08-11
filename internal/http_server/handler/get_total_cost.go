package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"task_manager/internal/http_server/dto"
	"task_manager/internal/lib/api/er"
	"task_manager/internal/lib/api/resp"
	valid "task_manager/internal/lib/api/valid"
	"task_manager/internal/lib/logger/sl"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
)

type TotalCostResponse struct {
	TotalCost int64 `json:"total_cost"`
}

// GetTotalCostHandler godoc
// @Summary      Get total user subscription cost
// @Description  Returns the total cost of a user's subscriptions for the specified period
// @Tags Total Cost
// @Accept       json
// @Produce      json
// @Param        request body dto.TotalCost true "Request data"
// @Success      200 {object} TotalCostResponse
// @Failure      400 {object} resp.ErrorResponse "Invalid request"
// @Failure      500 {object} resp.ErrorResponse "Server error"
// @Router       /subscriptions/total_cost [get]
func (h *UserSubscriptionHandler) GetTotalCostHandler(w http.ResponseWriter, r *http.Request) {
	const op = "handler.GetTotalCostHandler"

	ctx, cancel := context.WithTimeout(r.Context(), h.timeOut)
	defer cancel()

	log := h.log.With(
		slog.String("op", op),
		slog.String("request_url", middleware.GetReqID(ctx)),
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

	totalCost, err := h.service.TotalCost(ctx, req)
	if err != nil {
		log.Error("failed to get total cost", sl.Err(err))
		if msg, code, ok := er.MapErrorToStatus(err); ok {
			resp.Error(w, msg, code)
			return
		}

		resp.Error(w, "failed to get total cost", http.StatusInternalServerError)
		return
	}

	response := TotalCostResponse{TotalCost: totalCost}

	resp.ResponseOk(w, response, http.StatusOK)
}
