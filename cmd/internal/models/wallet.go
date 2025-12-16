package models

import "github.com/google/uuid"

type Wallet struct {
	ID      uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	Balance int64     `gorm:"not null" json:"balance"`
}
