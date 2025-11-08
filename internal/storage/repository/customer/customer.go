package customer

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"user-service/internal/domain/models"
	"user-service/internal/storage"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, customer *models.Customer) error {
	const op = "repository.customer.Create"

	query := `
        INSERT INTO customers (id, first_name, last_name, gender, timezone, birthday, user_id)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `

	_, err := r.db.ExecContext(ctx, query,
		customer.ID,
		customer.FirstName,
		customer.LastName,
		customer.Gender,
		customer.Timezone,
		customer.Birthday,
		customer.UserID,
	)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*models.Customer, error) {
	const op = "repository.customer.GetByID"

	query := `
        SELECT id, first_name, last_name, gender, timezone, birthday, user_id, created_at
        FROM customers
        WHERE id = $1
    `

	var customer models.Customer
	err := r.db.GetContext(ctx, &customer, query, id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrUserNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &customer, nil
}

func (r *Repository) GetByUserID(ctx context.Context, userID uuid.UUID) (*models.Customer, error) {
	const op = "repository.customer.GetByUserID"

	query := `
        SELECT id, first_name, last_name, gender, timezone, birthday, user_id, created_at
        FROM customers
        WHERE user_id = $1
    `

	var customer models.Customer
	err := r.db.GetContext(ctx, &customer, query, userID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, storage.ErrUserNotFound
		}
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &customer, nil
}

func (r *Repository) Update(ctx context.Context, id uuid.UUID, customer *models.Customer) error {
	const op = "repository.customer.Update"

	query := `
        UPDATE customers
        SET first_name = $1, last_name = $2, gender = $3, timezone = $4, birthday = $5
        WHERE id = $6
    `

	result, err := r.db.ExecContext(ctx, query,
		customer.FirstName,
		customer.LastName,
		customer.Gender,
		customer.Timezone,
		customer.Birthday,
		id,
	)

	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if rowsAffected == 0 {
		return storage.ErrUserNotFound
	}

	return nil
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	const op = "repository.customer.Delete"

	query := `DELETE FROM customers WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if rowsAffected == 0 {
		return storage.ErrUserNotFound
	}

	return nil
}
