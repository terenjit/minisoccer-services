package clients

import (
	"time"

	"github.com/google/uuid"
)

type FieldResponse struct {
	Code    int       `json:"code"`
	Status  string    `json:"status"`
	Message string    `json:"message"`
	Data    FieldData `json:"data"`
}

type FieldData struct {
	UUID         uuid.UUID  `json:"uuid"`
	FieldName    string     `json:"name"`
	PricePerHour float64    `json:"pricePerHour"`
	Date         string     `json:"date"`
	StartTime    string     `json:"startTime"`
	EndTime      string     `json:"endTime"`
	Status       string     `json:"status"`
	CreatedAt    *time.Time `json:"createdAt"`
	UpdatedAt    *time.Time `json:"updatedAt"`
}
