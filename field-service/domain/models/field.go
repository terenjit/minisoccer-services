package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Field struct {
	ID            uint           `gorm:"primaryKey;autoIncrement"`
	UUID          uuid.UUID      `gorm:"type:uuid;not null"`
	Code          string         `gorm:"type:varchar(15);not null"`
	Name          string         `gorm:"type:varchar(255);not null"`
	PricePerHour  int            `gorm:"type:int;not null"`
	Images        pq.StringArray `gorm:"type:text[];not null"`
	CreatedAt     *time.Time
	UpdatdeAt     *time.Time
	DeletedAt     *time.Time
	FieldSchedule []FieldSchedule `gorm:"foreignKey:field_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
