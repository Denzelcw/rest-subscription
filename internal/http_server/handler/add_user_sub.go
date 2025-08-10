package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"task_manager/internal/http_server/dto"
	"task_manager/internal/lib/api/resp"
	valid "task_manager/internal/lib/api/valid"
	"task_manager/internal/lib/logger/sl"
	"task_manager/internal/storage"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
)

type CreateResponse struct {
	Id      int64  `json:"id"`
	Message string `json:"message"`
}

// AddUserSubscriptionHandler godoc
// @Summary Add user subscription
// @Description Adding user subscription to the database.
// @Accept json
// @Produce json
// @Param request body dto.CreateUserSubDTO true "Data for creating a user subscription"
// @Success 201 {object} CreateResponse "Subscription created successfully"
// @Failure 400 {object} resp.ErrorResponse "Invalid request"
// @Failure 409 {object} resp.ErrorResponse "User subscription already exists"
// @Failure 500 {object} resp.ErrorResponse "Server error"
// @Router /subscriptions [post]
func (h *UserSubscriptionHandler) AddUserSubscriptionHandler(w http.ResponseWriter, r *http.Request) {
	const op = "handler.AddUserSubscriptionHandler"

	ctx, cancel := context.WithTimeout(r.Context(), h.timeOut)
	defer cancel()

	log := h.log.With(
		slog.String("op", op),
		slog.String("request_url", middleware.GetReqID(ctx)),
	)

	var req dto.CreateUserSubDTO

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

	id, err := h.service.Add(ctx, req)
	if err != nil {
		log.Info("user subscription already exists")
		if errors.Is(err, storage.ErrUserSubExists) {
			resp.Error(w, "user subscription already exists", http.StatusConflict)
		} else {
			resp.Error(w, "failed to add user subscription", http.StatusInternalServerError)
		}
		return
	}

	response := CreateResponse{Id: id, Message: "User subscription created successfully"}

	resp.ResponseOk(w, response, http.StatusCreated)
}
