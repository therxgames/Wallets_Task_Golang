package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
	"wallet/cmd/internal/database"
	"wallet/cmd/internal/models"
)

var router *gin.Engine
var walletID uuid.UUID

func TestMain(m *testing.M) {
	os.Setenv("DB_NAME", "wallet_test")
	os.Setenv("DB_HOST", "postgres")
	os.Setenv("DB_USER", "postgres")
	os.Setenv("DB_PASSWORD", "postgres")
	os.Setenv("DB_PORT", "5432")

	database.InitTestDB()

	wallet := models.Wallet{
		ID:      uuid.New(),
		Balance: 0,
	}

	database.DB.Create(&wallet)
	walletID = wallet.ID

	gin.SetMode(gin.TestMode)
	router = gin.Default()
	router.POST("/api/v1/wallet", CreateWalletOperation)
	router.GET("/api/v1/wallets/:id", GetWallet)

	code := m.Run()
	os.Exit(code)
}

func TestDeposit(t *testing.T) {
	reqBody := map[string]interface{}{
		"wallet_id":      walletID.String(),
		"operation_type": "DEPOSIT",
		"amount":         1000,
	}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/wallet", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d, body: %s", resp.Code, resp.Body.String())
	}

	var wallet models.Wallet
	database.DB.First(&wallet, "id = ?", walletID)
	if wallet.Balance != 1000 {
		t.Fatalf("Expected balance 1000, got %d", wallet.Balance)
	}
}
func TestWithdraw(t *testing.T) {
	database.DB.Model(&models.Wallet{}).Where("id = ?", walletID).Update("balance", 2000)

	reqBody := WalletOperationInput{
		WalletID:      walletID.String(),
		OperationType: models.Withdraw,
		Amount:        500,
	}
	body, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("Failed to marshal request: %v", err)
	}

	req, _ := http.NewRequest(http.MethodPost, "/api/v1/wallet", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d, body: %s", resp.Code, resp.Body.String())
	}

	var wallet models.Wallet
	database.DB.First(&wallet, "id = ?", walletID)
	if wallet.Balance != 1500 {
		t.Fatalf("Expected balance 1500, got %d", wallet.Balance)
	}
}

func TestGetBalance(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "/api/v1/wallets/"+walletID.String(), nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", resp.Code)
	}
}

func TestConcurrentDeposits(t *testing.T) {
	database.DB.Model(&models.Wallet{}).Where("id = ?", walletID).Update("balance", 0)

	const goroutines = 100
	const depositAmount = 10
	wg := sync.WaitGroup{}
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			reqBody := WalletOperationInput{
				WalletID:      walletID.String(),
				OperationType: models.Deposit,
				Amount:        depositAmount,
			}
			body, err := json.Marshal(reqBody)
			if err != nil {
				t.Errorf("Failed to marshal request: %v", err)
				return
			}

			req, _ := http.NewRequest(http.MethodPost, "/api/v1/wallet", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			if resp.Code != http.StatusOK {
				t.Errorf("Expected status 200, got %d, body: %s", resp.Code, resp.Body.String())
			}
		}()
	}

	wg.Wait()

	var wallet models.Wallet

	database.DB.First(&wallet, "id = ?", walletID)
	expected := int64(goroutines * depositAmount)

	if wallet.Balance != expected {
		t.Fatalf("Expected balance %d, got %d", expected, wallet.Balance)
	}
}
