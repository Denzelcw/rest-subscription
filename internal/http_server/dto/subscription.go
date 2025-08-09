package dto

import (
	"github.com/google/uuid"
)

type CreateSubDTO struct {
	ServiceName string    `json:"service_name" validate:"min=3,max=255"`
	Price       int       `json:"price" validate:"min=0"`
	UserID      uuid.UUID `json:"user_id" validate:"required,uuid4"`
	StartDate   string    `json:"start_date" validate:"required"`
	EndDate     string    `json:"end_date,omitempty"`
}

type UpdateSubDTO struct {
	ID          int       `json:"id,omitempty"`
	ServiceName string    `json:"service_name" validate:"min=3,max=255"`
	Price       int       `json:"price" validate:"min=0"`
	UserID      uuid.UUID `json:"user_id" validate:"required,uuid4"`
	StartDate   string    `json:"start_date" validate:"required"`
	EndDate     string    `json:"end_date,omitempty"`
}

type TotalCost struct {
	ServiceName string    `json:"service_name" validate:"min=3,max=255"`
	UserID      uuid.UUID `json:"user_id" validate:"required,uuid4"`
	StartDate   string    `json:"start_date" validate:"required"`
	EndDate     string    `json:"end_date,omitempty"`
}
