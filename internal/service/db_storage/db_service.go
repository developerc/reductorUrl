package dbstorage

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/developerc/reductorUrl/internal/general"
	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func CreateMapUser(db *sql.DB) (map[string]bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	rows, err := db.QueryContext(ctx, "SELECT DISTINCT usr FROM url WHERE usr IS NOT NULL")
	if err != nil {
		//fmt.Println(err)
		return nil, err
	}
	defer rows.Close()
	mapUser := make(map[string]bool)
	for rows.Next() {
		var cookie string
		err = rows.Scan(&cookie)
		if err != nil {
			//fmt.Println(err)
			return nil, err
		}
		mapUser[cookie] = true
	}
	// проверяем на ошибки
	err = rows.Err()
	if err != nil {
		//fmt.Println(err)
		return nil, err
	}
	return mapUser, nil
}

func InsertBatch(arrLongURL []general.ArrLongURL, dbStorage string, usr string) error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*2)
	defer cancelFunc()
	conn, err := pgx.Connect(ctx, dbStorage)
	if err != nil {
		return err
	}
	defer conn.Close(ctx)
	batch := &pgx.Batch{}
	for _, longURL := range arrLongURL {
		batch.Queue("insert into url( original_url, usr) values ($1, $2)", longURL.OriginalURL, usr)
	}
	br := conn.SendBatch(ctx, batch)
	_, err = br.Exec()
	if err != nil {
		return err
	}

	return nil
}

func SetDelBatch(arrShortURL []string, dbStorage string, usr string) error {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*2)
	defer cancelFunc()
	conn, err := pgx.Connect(ctx, dbStorage)
	if err != nil {
		return err
	}
	defer conn.Close(ctx)
	batch := &pgx.Batch{}
	for _, shortURL := range arrShortURL {
		batch.Queue("UPDATE url SET is_deleted = true WHERE uuid = $1 AND usr = $2", shortURL, usr)
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
	const cr string = "CREATE TABLE IF NOT EXISTS url( uuid serial primary key, original_url TEXT CONSTRAINT must_be_different UNIQUE, usr TEXT, is_deleted BOOLEAN NOT NULL DEFAULT FALSE)"
	_, err := db.ExecContext(ctx, cr)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func GetLongByUUID(db *sql.DB, uuid int) (string, bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	row := db.QueryRowContext(ctx, "SELECT original_url, is_deleted FROM url WHERE uuid=$1", uuid)
	var longURL string
	var isDeleted bool
	err := row.Scan(&longURL, &isDeleted)
	if err != nil {
		return "", false, err
	}
	return longURL, isDeleted, nil
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

func InsertRecord(db *sql.DB, originalURL string, usr string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err := db.ExecContext(ctx, "insert into url( original_url, usr) values ($1, $2)", originalURL, usr)

	if err != nil {
		return "", err
	}

	shURL, err := GetShortByOriginalURL(db, originalURL)
	if err != nil {
		return "", err
	}
	return shURL, nil
}

func ListRepoURLs(db *sql.DB, addresBase string, usr string) ([]general.ArrRepoURL, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	rows, err := db.QueryContext(ctx, "SELECT uuid, original_url FROM url WHERE usr = $1", usr)
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
		//fmt.Println(repoURL)
		if err != nil {
			//fmt.Println(err)
			return nil, err
		}
		repoURL.ShortURL = addresBase + "/" + repoURL.ShortURL
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
