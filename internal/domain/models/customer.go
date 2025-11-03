package models

import (
	"time"

	"github.com/google/uuid"
)

type Customer struct {
	ID        uuid.UUID `db:"id" json:"id"`
	FirstName string    `db:"first_name" json:"first_name"`
	LastName  string    `db:"last_name" json:"last_name"`
	Gender    string    `db:"gender" json:"gender"`
	Timezone  string    `db:"timezone" json:"timezone"`
	Birthday  time.Time `db:"birthday" json:"birthday"`
	UserID    uuid.UUID `db:"user_id" json:"user_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}

type CustomerAddress struct {
	ID         uuid.UUID `db:"id" json:"id"`
	Address    string    `db:"address" json:"address"`
	Apartment  string    `db:"apartment" json:"apartment"`
	Floor      int       `db:"floor" json:"floor"`
	Comments   string    `db:"comments" json:"comments"`
	CustomerID uuid.UUID `db:"customer_id" json:"customer_id"`
	CreatedAt  time.Time `db:"created_at" json:"created_at"`
}
