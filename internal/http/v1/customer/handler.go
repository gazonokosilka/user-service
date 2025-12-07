package customer

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"user-service/internal/domain/dto"
	"user-service/internal/domain/models"
	"user-service/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type CustomerService interface {
	CreateCustomer(ctx context.Context, req *dto.CreateCustomerRequest) (*models.Customer, error)
	GetCustomer(ctx context.Context, id uuid.UUID) (*models.Customer, error)
	GetAllCustomers(ctx context.Context) ([]models.Customer, error)
	UpdateCustomer(ctx context.Context, id uuid.UUID, req *dto.UpdateCustomerRequest) (*models.Customer, error)
}

type Handler struct {
	log     *slog.Logger
	service CustomerService
}

func NewHandler(log *slog.Logger, service CustomerService) *Handler {
	return &Handler{
		log:     log,
		service: service,
	}
}

func (h *Handler) CreateCustomer(w http.ResponseWriter, r *http.Request) {
	const op = "handler.customer.CreateCustomer"

	log := h.log.With(slog.String("op", op))

	var req dto.CreateCustomerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Warn("failed to decode request", slog.String("error", err.Error()))
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	customer, err := h.service.CreateCustomer(r.Context(), &req)
	if err != nil {
		log.Error("failed to create customer", slog.String("error", err.Error()))
		respondWithError(w, http.StatusInternalServerError, "failed to create customer")
		return
	}

	respondWithJSON(w, http.StatusCreated, customer)
}

func (h *Handler) GetCustomer(w http.ResponseWriter, r *http.Request) {
	const op = "handler.customer.GetCustomer"

	log := h.log.With(slog.String("op", op))

	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		log.Warn("invalid customer id", slog.String("error", err.Error()))
		respondWithError(w, http.StatusBadRequest, "invalid customer id")
		return
	}

	customer, err := h.service.GetCustomer(r.Context(), id)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			respondWithError(w, http.StatusNotFound, "customer not found")
			return
		}
		log.Error("failed to get customer", slog.String("error", err.Error()))
		respondWithError(w, http.StatusInternalServerError, "failed to get customer")
		return
	}

	respondWithJSON(w, http.StatusOK, customer)
}

func (h *Handler) GetAllCustomers(w http.ResponseWriter, r *http.Request) {
	const op = "handler.customer.GetAllCustomers"

	log := h.log.With(slog.String("op", op))

	customers, err := h.service.GetAllCustomers(r.Context())
	if err != nil {
		log.Error("failed to get customers", slog.String("error", err.Error()))
		respondWithError(w, http.StatusInternalServerError, "failed to get customers")
		return
	}

	respondWithJSON(w, http.StatusOK, customers)
}

func (h *Handler) UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	const op = "handler.customer.UpdateCustomer"

	log := h.log.With(slog.String("op", op))

	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		log.Warn("invalid customer id", slog.String("error", err.Error()))
		respondWithError(w, http.StatusBadRequest, "invalid customer id")
		return
	}

	var req dto.UpdateCustomerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Warn("failed to decode request", slog.String("error", err.Error()))
		respondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	customer, err := h.service.UpdateCustomer(r.Context(), id, &req)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			respondWithError(w, http.StatusNotFound, "customer not found")
			return
		}
		log.Error("failed to update customer", slog.String("error", err.Error()))
		respondWithError(w, http.StatusInternalServerError, "failed to update customer")
		return
	}

	respondWithJSON(w, http.StatusOK, customer)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}
