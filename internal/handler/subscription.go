package handler

import (
	"context"
	"errors"
	"github.com/denyshuzovskyi/nimbus-notify/internal/dto"
	commonerrors "github.com/denyshuzovskyi/nimbus-notify/internal/error"
	"github.com/go-playground/validator/v10"
	"log/slog"
	"net/http"
)

type SubscriptionService interface {
	Subscribe(context.Context, dto.SubscriptionRequest) error
	Confirm(context.Context, string) error
	Unsubscribe(context.Context, string) error
}

type SubscriptionHandler struct {
	subscriptionService SubscriptionService
	validator           *validator.Validate
	log                 *slog.Logger
}

func NewSubscriptionHandler(subscriptionService SubscriptionService, validator *validator.Validate, log *slog.Logger) *SubscriptionHandler {
	return &SubscriptionHandler{
		subscriptionService: subscriptionService,
		validator:           validator,
		log:                 log,
	}
}

func (h *SubscriptionHandler) Subscribe(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		h.log.Error("error parsing form", "error", err)
		return
	}

	var subscriptionReq dto.SubscriptionRequest
	subscriptionReq.Email = r.FormValue("email")
	subscriptionReq.City = r.FormValue("city")
	subscriptionReq.Frequency = r.FormValue("frequency")

	if err = h.validator.Struct(subscriptionReq); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		h.log.Error("error validating data", "error", err)
		return
	}

	if err = h.subscriptionService.Subscribe(r.Context(), subscriptionReq); err != nil {
		if errors.Is(err, commonerrors.LocationNotFound) {
			http.Error(w, "invalid input", http.StatusBadRequest)
			h.log.Error("couldn't validate city", "error", err)
			return
		} else if errors.Is(err, commonerrors.SubscriptionAlreadyExists) {
			http.Error(w, "email already subscribed", http.StatusConflict)
			h.log.Error("subscription already exists", "error", err)
			return
		}

		http.Error(w, "", http.StatusInternalServerError)
		h.log.Error("error making subscription", "error", err)
		return
	}
}

func (h *SubscriptionHandler) Confirm(w http.ResponseWriter, r *http.Request) {
	token := r.PathValue("token")

	if err := h.subscriptionService.Confirm(r.Context(), token); err != nil {
		if errors.Is(err, commonerrors.InvalidToken) {
			http.Error(w, "invalid token", http.StatusBadRequest)
			h.log.Error("invalid token", "error", err)
			return
		} else if errors.Is(err, commonerrors.TokenNotFound) {
			http.Error(w, "token not found", http.StatusNotFound)
			h.log.Error("token not found", "error", err)
			return
		}

		http.Error(w, "", http.StatusInternalServerError)
		h.log.Error("error confirming subscription", "error", err)
		return
	}
}

func (h *SubscriptionHandler) Unsubscribe(w http.ResponseWriter, r *http.Request) {
	token := r.PathValue("token")

	if err := h.subscriptionService.Unsubscribe(r.Context(), token); err != nil {
		if errors.Is(err, commonerrors.InvalidToken) {
			http.Error(w, "invalid token", http.StatusBadRequest)
			h.log.Error("invalid token", "error", err)
			return
		} else if errors.Is(err, commonerrors.TokenNotFound) {
			http.Error(w, "token not found", http.StatusNotFound)
			h.log.Error("token not found", "error", err)
			return
		}

		http.Error(w, "", http.StatusInternalServerError)
		h.log.Error("error confirming subscription", "error", err)
		return
	}
}
