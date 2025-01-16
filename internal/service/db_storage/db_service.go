package dbstorage

import (
	"context"
	"database/sql"
	"strconv"
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

func GetShortByOriginalURL(db *sql.DB, originalURL string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	row := db.QueryRowContext(ctx, "SELECT uuid FROM url WHERE original_url=$1", originalURL)
	var shURL int
	err := row.Scan(&shURL)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(shURL), err
}

func InsertRecord(db *sql.DB, originalURL string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err := db.ExecContext(ctx, "insert into url( original_url) values ($1)", originalURL)

	if err != nil {
		return "", err
	}

	shURL, err := GetShortByOriginalURL(db, originalURL)
	if err != nil {
		return "", err
	}
	return shURL, nil
}

func ListRepoURLs(db *sql.DB, addresBase string) ([]general.ArrRepoURL, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	rows, err := db.QueryContext(ctx, "SELECT uuid, original_url FROM url ")
	if err != nil {
		//fmt.Println(err)
		return nil, err
	}
	defer rows.Close()
	arrRepoURL := make([]general.ArrRepoURL, 0)
	// пробегаем по всем записям
	for rows.Next() {
		//repoURL := general.ArrRepoURL{}

		var repoURL general.ArrRepoURL
		err = rows.Scan(&repoURL.ShortURL, &repoURL.OriginalURL)
		repoURL.ShortURL = addresBase + "/" + repoURL.ShortURL
		//fmt.Println(repoURL)
		if err != nil {
			//fmt.Println(err)
			return nil, err
		}
		arrRepoURL = append(arrRepoURL, repoURL)
	}
	// проверяем на ошибки
	err = rows.Err()
	if err != nil {
		//fmt.Println(err)
		return nil, err
	}
	return arrRepoURL, nil
}
