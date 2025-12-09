package main

import (
	"fmt"
	"log"

	"github.com/Infamous003/follow-service/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Configuration loaded: %+v\n", cfg)
}
