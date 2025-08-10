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

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
)

type CreateResponse struct {
	Id      int64  `json:"id"`
	Message string `json:"message"`
}

// AddUserSubscriptionHandler godoc
// @Summary Add user subscription
// @Description Adding user subsctiption to db.
// @Tags User subscriptions
// @Accept json
// @Produce json
// @Param request body dto.CreateUserSubDTO true "Данные для создания подписки"
// @Success 201 {object} CreateResponse "Успешное создание подписки"
// @Failure 400 {object} resp.ErrorResponse "Некорректный запрос"
// @Failure 500 {object} resp.ErrorResponse "Ошибка сервера"
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
		log.Error("failed to add user subscription", sl.Err(err))

		resp.Error(w, "failed to add user subscription", http.StatusInternalServerError)
		return
	}

	response := CreateResponse{Id: id, Message: "User subscription created successfully"}

	resp.ResponseOk(w, response, http.StatusCreated)
}
