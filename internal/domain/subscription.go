package domain

import (
	"github.com/google/uuid"
)

type Subscription struct {
	ID          string    `json:"id,omitempty"`
	ServiceName string    `json:"service_name,omitempty"`
	Price       int       `json:"price,omitempty"`
	UserID      uuid.UUID `json:"user_id,omitempty"`
	StartDate   string    `json:"start_date,omitempty"`
	EndDate     string    `json:"end_date,omitempty"`
}
