package main

import (
	"log"

	"github.com/Infamous003/follow-service/internal/config"
	"github.com/Infamous003/follow-service/internal/database"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.NewDB(cfg.DB)
	if err != nil {
		log.Fatal("failed to connect to the database:", err)
	}
	defer db.Close()
}
