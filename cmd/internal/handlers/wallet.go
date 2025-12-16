package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"wallet/cmd/internal/database"
	"wallet/cmd/internal/models"
)

func GetWallet(c *gin.Context) {
	var wallet models.Wallet

	id := c.Param("id")

	if err := database.DB.Where("id = ?", id).First(&wallet).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "There is no wallet with this id."})
		fmt.Println(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"wallet": wallet,
	})
}
