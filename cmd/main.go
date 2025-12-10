package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Infamous003/follow-service/internal/config"
	"github.com/Infamous003/follow-service/internal/database"
	"github.com/Infamous003/follow-service/internal/handler"
	"github.com/Infamous003/follow-service/internal/repository"
	"github.com/Infamous003/follow-service/internal/service"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
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

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Post("/users", userHandler.CreateUser)

	s := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: router,
	}

	log.Printf("Starting server on %s\n", s.Addr)

	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
