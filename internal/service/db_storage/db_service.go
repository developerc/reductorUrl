package dbstorage

import (
	"context"
	"database/sql"
	"time"

	"github.com/developerc/reductorUrl/internal/general"
	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func InsertBatch(arrLongURL []general.ArrLongURL, dbStorage string) error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*2)
	defer cancelFunc()
	conn, err := pgx.Connect(ctx, dbStorage)
	if err != nil {
		return err
	}
	defer conn.Close(ctx)
	batch := &pgx.Batch{}
	for _, longURL := range arrLongURL {
		batch.Queue("insert into url( original_url) values ($1)", longURL.OriginalURL)
	}
	br := conn.SendBatch(ctx, batch)
	_, err = br.Exec()
	if err != nil {
		return err
	}

	return nil
}

func CreateTable(db *sql.DB) error {
	const duration uint = 20
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(duration)*time.Second)
	defer cancel()
	const cr string = "CREATE TABLE IF NOT EXISTS url( uuid serial primary key, original_url TEXT CONSTRAINT must_be_different UNIQUE)"
	_, err := db.ExecContext(ctx, cr)
	if err != nil {
		return err
	}
	return nil
}

func GetLongByUUID(db *sql.DB, uuid int) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	row := db.QueryRowContext(ctx, "SELECT original_url FROM url WHERE uuid=$1", uuid)
	var longURL string
	err := row.Scan(&longURL)
	if err != nil {
		return "", err
	}
	return longURL, nil
}
