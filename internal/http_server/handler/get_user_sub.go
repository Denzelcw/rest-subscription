package handler

import (
	"context"
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

// GetUserSubscriptionHandler godoc
// @Summary      Получение подписки пользователя
// @Description  Возвращает информацию о подписке пользователя по её ID
// @Tags         subscriptions
// @Param        id   path      int  true  "ID подписки"
// @Success      200  {object}  domain.UserSubscription
// @Failure      400  {object}  resp.ErrorResponse "Неверный ID"
// @Failure      404  {object}  resp.ErrorResponse "Подписка не найдена"
// @Failure      500  {object}  resp.ErrorResponse "Ошибка получения данных"
// @Router       /subscriptions/{id} [get]
func (h *UserSubscriptionHandler) GetUserSubscriptionHandler(w http.ResponseWriter, r *http.Request) {
	const op = "handler.GetUserSubscriptionHandler"

	ctx, cancel := context.WithTimeout(r.Context(), h.timeOut)
	defer cancel()

	log := h.log.With(
		slog.String("op", op),
		slog.String("request_url", middleware.GetReqID(ctx)),
	)

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Error("failed to parse id", sl.Err(err))

		resp.Error(w, "invalid user subscription ID", http.StatusBadRequest)
		return
	}

	subscription, err := h.service.GetById(ctx, id)
	if err != nil {
		log.Error("failed to get user subscription", sl.Err(err))
		if errors.Is(err, storage.ErrNotFound) {
			resp.Error(w, "user subscription not found", http.StatusNotFound)
		} else {
			resp.Error(w, "failed to get user subscription", http.StatusInternalServerError)
		}
		return
	}

	resp.ResponseOk(w, subscription, http.StatusOK)
}
