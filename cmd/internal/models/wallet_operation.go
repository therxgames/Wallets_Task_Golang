package models

import "github.com/google/uuid"

type OperationType string

const (
	Deposit  OperationType = "DEPOSIT"
	Withdraw OperationType = "WITHDRAW"
)

type WalletOperation struct {
	ID            uuid.UUID     `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	WalletID      uuid.UUID     `gorm:"type:uuid;not null;index"`
	Wallet        *Wallet       `gorm:"foreignKey:WalletID" json:"material"`
	OperationType OperationType `gorm:"type:varchar(10);not null"`
	Amount        int64         `gorm:"not null"`
}

func (ot OperationType) IsValid() bool {
	return ot == Deposit || ot == Withdraw
}
