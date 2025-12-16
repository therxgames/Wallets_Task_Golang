package main

import (
	"wallet/cmd/internal/config"
	"wallet/cmd/internal/database"
	"wallet/cmd/internal/router"
)

func main() {
	config.Init()
	database.Init()
	router.Init()
}
