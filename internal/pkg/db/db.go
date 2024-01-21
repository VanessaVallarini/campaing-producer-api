package db

import (
	"campaing-producer-service/internal/config"
	"context"
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/lockp111/go-easyzap"
)

type IDb interface {
	Ping(ctx context.Context) error
	Query(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRow(query string, args ...any) *sql.Row
	Close()
}

type Db struct {
	conn *sql.DB
}

func NewDbClient(cfg config.Config) *Db {
	db, err := sql.Open(cfg.DatabaseConfig.PostgresDriver, cfg.DatabaseConfig.DatabaseConnStr)
	if err != nil {
		easyzap.Panicf("configuration error: %v", err)
	}

	return &Db{
		conn: db,
	}
}

func (db *Db) Ping(ctx context.Context) error {
	err := db.conn.Ping()
	if err != nil {
		easyzap.Error(ctx, err, "failed to ping database")

		return err
	}

	return nil
}

func (db *Db) Query(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return db.conn.Query(query, args)
}

func (db *Db) QueryRow(query string, args ...any) *sql.Row {
	return db.conn.QueryRow(query, args)
}

func (db *Db) Close() {
	db.conn.Close()
}
