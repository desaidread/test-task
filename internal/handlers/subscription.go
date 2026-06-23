package handlers

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"testTask/internal/models"
	"testTask/internal/repository"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Handler struct {
	repo *repository.SubscriptionRepository
}

func New(repo *repository.SubscriptionRepository) *Handler {
	return &Handler{repo: repo}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var input models.CreateSubscriptionInput

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {

		http.Error(w, "Unable to decode json", http.StatusBadRequest)
		return
	}
	id, err := h.repo.Create(r.Context(), input)
	if err != nil {
		slog.Error("Unable to create subscription: ", "err", err)
		http.Error(w, "Unable to create subscription", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"id": id})
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Cannot read id", http.StatusBadRequest)
		return
	}
	sub, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "subcscription not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Cannot get id", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sub)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	var updateSubscription models.CreateSubscriptionInput

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Cannot read id", http.StatusBadRequest)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&updateSubscription); err != nil {
		http.Error(w, "Unable to decode json", http.StatusBadRequest)
		return
	}

	if err := h.repo.Update(r.Context(), id, updateSubscription); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "subscription with that id not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Unable to update subscription", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Cannot read id", http.StatusBadRequest)
		return
	}

	if err := h.repo.Delete(r.Context(), id); err != nil {
		http.Error(w, "unable to delete subscription", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {

	list, err := h.repo.List(r.Context())
	if err != nil {
		http.Error(w, "Cannot list subscriptions", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)
}

func (h *Handler) TotalCost(w http.ResponseWriter, r *http.Request) {
	userIdStr := r.URL.Query().Get("user_id")

	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		http.Error(w, "unable to parse user_id", http.StatusBadRequest)
		return
	}

	fromStr := r.URL.Query().Get("from")
	toStr := r.URL.Query().Get("to")

	if fromStr == "" || toStr == "" {
		http.Error(w, "from and to are required", http.StatusBadRequest)
		return
	}

	from, err := time.Parse("01-2006", fromStr)
	if err != nil {
		http.Error(w, "Unable to parse date", http.StatusBadRequest)
		return
	}

	to, err := time.Parse("01-2006", toStr)
	if err != nil {
		http.Error(w, "Unable to parse date", http.StatusBadRequest)
		return
	}

	serviceNameStr := r.URL.Query().Get("service_name")

	totalCost, err := h.repo.TotalCost(r.Context(), userId, serviceNameStr, from, to)
	if err != nil {
		http.Error(w, "Cannot count total cost", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]int{"total": totalCost})

}
