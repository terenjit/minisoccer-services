package dto

import (
	"mime/multipart"
	"time"

	"github.com/google/uuid"
)

type FieldRequest struct {
	Code         string                 `form:"code" validate:"required"`
	Name         string                 `form:"name" validate:"required"`
	PricePerHour int                    `form:"pricePerHour" validate:"required"`
	Images       []multipart.FileHeader `form:"images" validate:"required"`
}

type UpdateFieldRequest struct {
	Code         string                 `form:"code" validate:"required"`
	Name         string                 `form:"name" validate:"required"`
	PricePerHour int                    `form:"pricePerHour" validate:"required"`
	Images       []multipart.FileHeader `form:"images"`
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
	SortOrder  *string `form:"sortOrder"`
}
