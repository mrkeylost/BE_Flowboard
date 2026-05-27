package model

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	InternalID     int64     `json:"internal_id" db:"internal_id" gorm:"primaryKey;autoIncrement"`
	PublicID       uuid.UUID `json:"public_id" db:"public_id"`
	CardInternalID int64     `json:"card_internal_id" db:"card_internal_id"`
	CardID         uuid.UUID `json:"card_id" db:"card_id"`
	UserInternalID int64     `json:"user_internal_id" db:"user_internal_id"`
	UserID         uuid.UUID `json:"user_id" db:"user_id"`
	Message        string    `json:"message" db:"message"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}
