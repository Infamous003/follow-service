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

	followRepo := repository.NewFollowRepository(db)
	followService := service.NewFollowService(followRepo, userRepo)
	followHandler := handler.NewFollowHandler(followService)

	router := chi.NewRouter()
	router.Use(middleware.Logger)

	router.Post("/users", userHandler.CreateUser)
	router.Get("/users/{id}", userHandler.GetUserByID)
	router.Get("/users", userHandler.ListUsers)
	router.Get("/users/{id}/followers", followHandler.ListFollowers)
	router.Get("/users/{id}/following", followHandler.ListFollowing)

	router.Post("/follow", followHandler.FollowUser)
	router.Post("/unfollow", followHandler.UnfollowUser)

	s := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: router,
	}

	log.Printf("Starting server on %s\n", s.Addr)

	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
