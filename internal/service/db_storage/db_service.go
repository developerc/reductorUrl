package dbstorage

import (
	"context"
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
