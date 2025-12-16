package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"net/http"
	"wallet/cmd/internal/database"
	"wallet/cmd/internal/models"
)

type WalletOperationInput struct {
	WalletID      string               `json:"wallet_id" binding:"required"`
	OperationType models.OperationType `json:"operation_type" binding:"required"`
	Amount        int64                `json:"amount" binding:"required,gt=0"`
}

func CreateWalletOperation(c *gin.Context) {
	var input WalletOperationInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	walletID, err := uuid.Parse(input.WalletID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid wallet_id"})
		return
	}

	if !input.OperationType.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid operation type"})
		return
	}

	err = database.DB.Transaction(func(tx *gorm.DB) error {
		tx.Exec("SET LOCAL lock_timeout = '5s'")

		var wallet models.Wallet

		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&wallet, "id = ?", walletID).Error; err != nil {
			return err
		}

		switch input.OperationType {
		case models.Deposit:
			wallet.Balance += input.Amount
		case models.Withdraw:
			if wallet.Balance < input.Amount {
				return errors.New("insufficient funds")
			}
			wallet.Balance -= input.Amount
		}

		if err := tx.Save(&wallet).Error; err != nil {
			return err
		}

		walletOperation := models.WalletOperation{
			WalletID:      wallet.ID,
			OperationType: input.OperationType,
			Amount:        input.Amount,
		}

		return tx.Create(&walletOperation).Error
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
