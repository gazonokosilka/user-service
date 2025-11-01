package psql

import (
	"fmt"
	"log"
	"runtime/debug"

	"user-service/internal/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Storage struct {
	db *sqlx.DB
}

// Init connection to database
func Init(cfg config.PostgresConfig) *Storage {
	const op = "storage.psql.Init"

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DbName, cfg.SslMode)

	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		panic(fmt.Sprintf("%s: failed to open db: %v", op, err))
	}

	if err = db.Ping(); err != nil {
		panic(fmt.Sprintf("%s: failed to ping db: %v", op, err))
	}

	return &Storage{db: db}
}

func (s *Storage) GetDB() *sqlx.DB {
	return s.db
}

// Close connection
func (s *Storage) Close() {
	if s.db != nil {
		log.Printf("Closing DB (caller):\n%s", debug.Stack())
		log.Printf("DB stats before close: InUse=some Idle=some")
		s.db.Close()
	}
}
