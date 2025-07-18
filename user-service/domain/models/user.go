package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uint      `gorm:"primaryKey;autoIncrement"`
	UUID        uuid.UUID `gorm:"type:uuid;not null"`
	Name        string    `gorm:"type:varhcar(100);not null"`
	Username    string    `gorm:"type:varhcar(15);not null"`
	Password    string    `gorm:"type:varhcar(255);not null"`
	PhoneNumber string    `gorm:"type:varhcar(15);not null"`
	Email       string    `gorm:"type:varhcar(100);not null"`
	RoleID      uint      `gorm:"type:uint;not null"`
	CreatedAt   *time.Time
	UpdatdeAt   *time.Time
	Role        Role `gorm:"foreignKey:role_id;references:id;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
