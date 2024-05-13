package db

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"

	_ "github.com/jackc/pgx/v4/stdlib" // load pgx driver for PostgreSQL
)

func CreateConnection() (*sqlx.DB, error) {
	maxConn, _ := strconv.Atoi(os.Getenv("DB_MAX_CONNECTIONS"))
	maxIdleConn, _ := strconv.Atoi(os.Getenv("DB_MAX_IDLE_CONNECTIONS"))
	maxLifetimeConn, _ := strconv.Atoi(os.Getenv("DB_MAX_LIFETIME_CONNECTIONS"))

	url := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSL_MODE"),
	)

	db, err := sqlx.Connect("pgx", url)
	if err != nil {
		return nil, fmt.Errorf("error connected to database, %w", err)
	}

	db.SetMaxOpenConns(maxConn)
	db.SetMaxIdleConns(maxIdleConn)
	db.SetConnMaxIdleTime(time.Duration(maxLifetimeConn))

	if err := db.Ping(); err != nil {
		defer db.Close()
		return nil, fmt.Errorf("error sent ping to database, %w ", err)
	}

	return db, nil
}
