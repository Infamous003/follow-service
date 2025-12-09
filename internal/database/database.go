package database

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/Infamous003/follow-service/internal/config"
	_ "github.com/lib/pq"
)

func NewDB(cfg config.DBConfig) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.DSN)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxIdleTime(cfg.MaxIdleTime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for i := 0; i < 5; i++ {
		err := db.PingContext(ctx)
		if err == nil {
			break
		}
		log.Println("Database ping failed, retrying...")
		time.Sleep(2 * time.Second)
		if i == 4 {
			return nil, err
		}
	}

	log.Println("Database connection established")
	return db, nil
}
