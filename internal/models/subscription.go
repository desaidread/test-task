package models

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID          int        `json:"id"`
	ServiceName string     `json:"service_name"`
	UserId      uuid.UUID  `json:"user_id"`
	Price       int        `json:"price"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     *time.Time `json:"end_date,omitempty"`
}

type CreateSubscriptionInput struct {
	ServiceName string    `json:"service_name" validate:"required"`
	Price       int       `json:"price" validate:"required"`
	UserId      uuid.UUID `json:"user_id" validate:"required"`
	StartDate   string    `json:"start_date" validate:"required"`
	EndDate     *string   `json:"end_date"`
}
