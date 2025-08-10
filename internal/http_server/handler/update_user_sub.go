package handler

import (
	"context"
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
)

// UpdateSubscriptionHandler godoc
// @Summary      Обновление подписки пользователя
// @Description  Обновляет данные подписки пользователя по её ID
// @Tags         subscriptions
// @Accept       json
// @Produce      json
// @Param        id   path      int                  true  "ID подписки"
// @Param        body body      dto.UpdateUserSubDTO true  "Данные для обновления подписки"
// @Success      201  {object}  domain.UserSubscription
// @Failure      400  {object}  resp.ErrorResponse "Некорректный ID или тело запроса"
// @Failure      500  {object}  resp.ErrorResponse "Ошибка при обновлении подписки"
// @Router       /subscriptions/{id} [put]
func (h *UserSubscriptionHandler) UpdateSubscriptionHandler(w http.ResponseWriter, r *http.Request) {
	const op = "handler.AddUserSubscriptionHandler"

	ctx, cancel := context.WithTimeout(r.Context(), h.timeOut)
	defer cancel()

	log := h.log.With(
		slog.String("op", op),
		slog.String("request_url", middleware.GetReqID(r.Context())),
	)

	var req dto.UpdateUserSubDTO

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

	sub, err := h.service.UpdateById(ctx, req)
	if err != nil {
		log.Error("failed to update subscription", sl.Err(err))

		resp.Error(w, "failed to update subscription", http.StatusInternalServerError)
		return
	}

	resp.ResponseOk(w, sub, http.StatusCreated)
}
