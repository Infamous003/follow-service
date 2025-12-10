package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Infamous003/follow-service/internal/config"
	"github.com/Infamous003/follow-service/internal/database"
	"github.com/Infamous003/follow-service/internal/handler"
	"github.com/Infamous003/follow-service/internal/repository"
	"github.com/Infamous003/follow-service/internal/service"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	// LOading configs
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	// Connecting to the DB
	db, err := database.NewDB(cfg.DB)
	if err != nil {
		log.Fatal("failed to connect to the database:", err)
	}
	defer db.Close()

	if err := runMigrations(cfg.DB.DSN); err != nil {
		log.Fatal("failed to run migrations:", err)
	}

	// Repos and services
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	followRepo := repository.NewFollowRepository(db)
	followService := service.NewFollowService(followRepo, userRepo)
	followHandler := handler.NewFollowHandler(followService)

	// Chi router and middlewares
	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Use(middleware.Logger)

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"}, // or specific domains like "https://example.com"
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // 5 minutes
	}))

	// User routes
	router.Post("/users", userHandler.CreateUser)
	router.Get("/users/{id}", userHandler.GetUserByID)
	router.Get("/users", userHandler.ListUsers)

	// Follow routes
	router.Post("/follow", followHandler.FollowUser)
	router.Post("/unfollow", followHandler.UnfollowUser)

	router.Get("/users/{id}/followers", followHandler.ListFollowers)
	router.Get("/users/{id}/following", followHandler.ListFollowing)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// a channel to store shutdown errors
	shutdownErr := make(chan error)

	go func() {
		// chan to store OS signal
		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		sig := <-quit
		log.Printf("caught sig: %+v, shutting down server", sig)

		ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
		defer cancel()

		err := server.Shutdown(ctx)
		if err != nil {
			shutdownErr <- err
		}

		log.Println("server stopped")
		shutdownErr <- nil
	}()

	log.Printf("Starting server on %s\n", server.Addr)

	err = server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		log.Fatal("Server failed to start:", err)
	}
	err = <-shutdownErr
	if err != nil {
		log.Fatal("failed to shutdown server gracefully:", err)
	}
	log.Println("server exited properly")
}

func runMigrations(dsn string) error {
	m, err := migrate.New(
		"file://migrations",
		dsn,
	)
	if err != nil {
		return err
	}
	err = m.Up()
	if errors.Is(err, migrate.ErrNoChange) {
		return nil
	}

	return err
}
