package dto

import (
	"mime/multipart"
	"time"

	"github.com/google/uuid"
)

type FieldRequest struct {
	Code         string                 `json:"code" validate:"required"`
	Name         string                 `json:"name" validate:"required"`
	PricePerHour int                    `json:"pricePerHour" validate:"required"`
	Images       []multipart.FileHeader `json:"images" validate:"required"`
}

type UpdateFieldRequest struct {
	Code         string                 `json:"code" validate:"required"`
	Name         string                 `json:"name" validate:"required"`
	PricePerHour int                    `json:"pricePerHour" validate:"required"`
	Images       []multipart.FileHeader `json:"images"`
}

type FieldResponse struct {
	UUID         uuid.UUID  `json:"uuid"`
	Code         string     `json:"code"`
	Name         string     `json:"name"`
	PricePerHour int        `json:"pricePerHour"`
	Images       []string   `json:"images"`
	CreatedAt    *time.Time `json:"createAt"`
	UpdateAt     *time.Time `json:"updateAt"`
}

type FieldDetailResponse struct {
	Code         string     `json:"code"`
	Name         string     `json:"name"`
	PricePerHour int        `json:"pricePerHour"`
	Images       []string   `json:"images"`
	CreatedAt    *time.Time `json:"createAt"`
	UpdateAt     *time.Time `json:"updateAt"`
}

type FieldRequestParam struct {
	Page       int     `form:"page" validate:"required"`
	Limit      int     `form:"limit" validate:"required"`
	SortColumn *string `form:"sortColumn"`
	SortOrder  int     `form:"sortOrder"`
}
