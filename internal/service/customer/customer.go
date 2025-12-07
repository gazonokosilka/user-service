package customer

import (
	"context"
	"fmt"
	"log/slog"

	"user-service/internal/domain/dto"
	"user-service/internal/domain/models"

	"github.com/google/uuid"
)

type CustomerRepository interface {
	Create(ctx context.Context, customer *models.Customer) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Customer, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*models.Customer, error)
	GetAll(ctx context.Context) ([]models.Customer, error)
	Update(ctx context.Context, id uuid.UUID, customer *models.Customer) error
}

type Service struct {
	log  *slog.Logger
	repo CustomerRepository
}

func New(log *slog.Logger, repo CustomerRepository) *Service {
	return &Service{
		log:  log,
		repo: repo,
	}
}

func (s *Service) CreateCustomer(ctx context.Context, req *dto.CreateCustomerRequest) (*models.Customer, error) {
	const op = "service.customer.CreateCustomer"

	log := s.log.With(slog.String("op", op))

	if err := req.Validate(); err != nil {
		log.Warn("validation failed", slog.String("error", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		log.Warn("invalid user_id", slog.String("error", err.Error()))
		return nil, fmt.Errorf("%s: invalid user_id: %w", op, err)
	}

	birthday, err := req.ParseBirthday()
	if err != nil {
		log.Warn("invalid birthday", slog.String("error", err.Error()))
		return nil, fmt.Errorf("%s: invalid birthday: %w", op, err)
	}

	customer := &models.Customer{
		ID:        uuid.New(),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Gender:    req.Gender,
		Timezone:  req.Timezone,
		Birthday:  birthday,
		UserID:    userID,
	}

	if err := s.repo.Create(ctx, customer); err != nil {
		log.Error("failed to create customer", slog.String("error", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("customer created", slog.String("customer_id", customer.ID.String()))

	return customer, nil
}

func (s *Service) GetCustomer(ctx context.Context, id uuid.UUID) (*models.Customer, error) {
	const op = "service.customer.GetCustomer"

	log := s.log.With(slog.String("op", op))

	customer, err := s.repo.GetByID(ctx, id)
	if err != nil {
		log.Error("failed to get customer", slog.String("error", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return customer, nil
}

func (s *Service) GetAllCustomers(ctx context.Context) ([]models.Customer, error) {
	const op = "service.customer.GetAllCustomers"

	log := s.log.With(slog.String("op", op))

	customers, err := s.repo.GetAll(ctx)
	if err != nil {
		log.Error("failed to get customers", slog.String("error", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("customers retrieved", slog.Int("count", len(customers)))

	return customers, nil
}

func (s *Service) UpdateCustomer(ctx context.Context, id uuid.UUID, req *dto.UpdateCustomerRequest) (*models.Customer, error) {
	const op = "service.customer.UpdateCustomer"

	log := s.log.With(slog.String("op", op))

	if err := req.Validate(); err != nil {
		log.Warn("validation failed", slog.String("error", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	existingCustomer, err := s.repo.GetByID(ctx, id)
	if err != nil {
		log.Error("failed to get customer", slog.String("error", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if req.FirstName != nil {
		existingCustomer.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		existingCustomer.LastName = *req.LastName
	}
	if req.Gender != nil {
		existingCustomer.Gender = *req.Gender
	}
	if req.Timezone != nil {
		existingCustomer.Timezone = *req.Timezone
	}
	if req.Birthday != nil {
		birthday, err := req.ParseBirthday()
		if err != nil {
			log.Warn("invalid birthday", slog.String("error", err.Error()))
			return nil, fmt.Errorf("%s: invalid birthday: %w", op, err)
		}
		existingCustomer.Birthday = *birthday
	}

	if err := s.repo.Update(ctx, id, existingCustomer); err != nil {
		log.Error("failed to update customer", slog.String("error", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("customer updated", slog.String("customer_id", id.String()))

	return existingCustomer, nil
}
