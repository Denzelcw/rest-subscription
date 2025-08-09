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

type CreateResponse struct {
	Id      int64  `json:"id"`
	Message string `json:"message"`
}

// @Summary Добавить новую подписку
// @Description Добавляет новую подписку в систему.
// @Tags Subscriptions
// @Accept json
// @Produce json
// @Param request body dto.CreateSubDTO true "Данные для создания подписки"
// @Success 201 {object} CreateResponse "Успешное создание подписки"
// @Failure 400 {object} resp.ErrorResponse "Некорректный запрос"
// @Failure 500 {object} resp.ErrorResponse "Ошибка сервера"
// @Router /subscriptions [post]
func (h *SubscriptionHandler) AddSubscriptionHandler(w http.ResponseWriter, r *http.Request) {
	const op = "handler.AddSubscriptionHandler"

	log := h.log.With(
		slog.String("op", op),
		slog.String("request_url", middleware.GetReqID(r.Context())),
	)

	var req dto.CreateSubDTO

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

	id, err := h.service.Add(req)
	if err != nil {
		log.Error("failed to add subscription", sl.Err(err))

		resp.Error(w, "failed to add subscription", http.StatusInternalServerError)
		return
	}

	response := CreateResponse{Id: id, Message: "Subscription created successfully"}

	resp.ResponseOk(w, response, http.StatusCreated)
}
