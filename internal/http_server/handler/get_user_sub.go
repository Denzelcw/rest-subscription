package handler

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"task_manager/internal/lib/api/er"
	"task_manager/internal/lib/api/resp"
	"task_manager/internal/lib/logger/sl"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// GetUserSubscriptionHandler godoc
// @Summary      Get user subscription
// @Description  Returns information about a user's subscription by its ID
// @Tags Subscription
// @Param        id   path      int  true  "Subscription ID"
// @Success      200  {object}  domain.UserSubscription "Request data"
// @Failure      400  {object}  resp.ErrorResponse "Invalid ID"
// @Failure      404  {object}  resp.ErrorResponse "User subscription not found"
// @Failure      500  {object}  resp.ErrorResponse "Server error"
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
		if msg, code, ok := er.MapErrorToStatus(err); ok {
			resp.Error(w, msg, code)
			return
		}

		resp.Error(w, "failed to get user subscription", http.StatusInternalServerError)
		return
	}

	resp.ResponseOk(w, subscription, http.StatusOK)
}
