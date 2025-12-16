package router

import (
	"github.com/gin-gonic/gin"
	"log"
	"wallet/cmd/internal/handlers"
)

func Init() {
	router := gin.Default()

	{
		v1 := router.Group("/api/v1")
		v1.POST("/wallet", handlers.CreateWalletOperation)
		v1.GET("/wallets/:id", handlers.GetWallet)
	}

	if err := router.Run(); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
