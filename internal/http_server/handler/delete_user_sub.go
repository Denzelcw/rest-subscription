package handler

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"
	"subscription/internal/lib/api/er"
	"subscription/internal/lib/api/resp"
	"subscription/internal/lib/logger/sl"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type DeleteResponse struct {
	Id      int    `json:"id"`
	Message string `json:"message"`
}

// DeleteUserSubscriptionHandler godoc
// @Summary      Delete user subscription
// @Description  Deletes a user subscription by ID
// @Tags Subscription
// @Param        id   path      int  true  "User subscription ID"
// @Success      200  {object}  DeleteResponse
// @Failure      400  {object}  resp.ErrorResponse "Invalid ID"
// @Failure      404  {object}  resp.ErrorResponse "User ubscription not found"
// @Failure      500  {object}  resp.ErrorResponse "Server error"
// @Router       /subscriptions/{id} [delete]
func (h *UserSubscriptionHandler) DeleteUserSubscriptionHandler(w http.ResponseWriter, r *http.Request) {
	const op = "handler.DeleteUserSubscriptionHandler"

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

	err = h.service.DeleteById(ctx, id)
	if err != nil {
		log.Error("failed to delete user subscription")
		if msg, code, ok := er.MapErrorToStatus(err); ok {
			resp.Error(w, msg, code)
			return
		}
		resp.Error(w, "failed to delete user subscription", http.StatusInternalServerError)
		return
	}

	response := DeleteResponse{
		Id:      id,
		Message: "user subscription successfully deleted",
	}

	resp.ResponseOk(w, response, http.StatusOK)
}
