// dbstorage пакет для размещения методов обработки запросов к базе данных.
package dbstorage

import (
	"context"
	"database/sql"
	"strconv"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/developerc/reductorUrl/internal/general"
)

// CreateMapUser создает Map пользователей читая таблицу при запуске приложения
func CreateMapUser(ctx context.Context, db *sql.DB) (map[string]bool, error) {
	rows, err := db.QueryContext(ctx, "SELECT DISTINCT usr FROM url WHERE usr IS NOT NULL")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	mapUser := make(map[string]bool)
	for rows.Next() {
		var cookie string
		err = rows.Scan(&cookie)
		if err != nil {
			return nil, err
		}
		mapUser[cookie] = true
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return mapUser, nil
}

// InsertBatch2 вставляет несколько длинных URL в таблицу
func InsertBatch2(ctx context.Context, arrLongURL []general.ArrLongURL, db *sql.DB, usr string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err = tx.Rollback(); err != nil {
			return
		}
	}()

	stmt, err := tx.PrepareContext(ctx,
		"insert into url( original_url, usr) values ($1, $2)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	for _, longURL := range arrLongURL {
		_, err := stmt.ExecContext(ctx, longURL.OriginalURL, usr)
		if err != nil {
			return err
		}
	}
	return tx.Commit()
}

// SetDelBatch2 в таблице делает отметку об удалении для нескольких коротких URL
func SetDelBatch(arrShortURL []string, db *sql.DB, usr string) error {
	_, err := db.Exec("UPDATE url SET is_deleted = true WHERE uuid = ANY($1) AND usr = $2", arrShortURL, usr)

	return err
}

// CreateTable создает таблицу url если она не существовала
func CreateTable(ctx context.Context, db *sql.DB) error {
	const duration uint = 20
	const cr string = "CREATE TABLE IF NOT EXISTS url( uuid serial primary key, " +
		"original_url TEXT CONSTRAINT must_be_different UNIQUE, usr TEXT, is_deleted BOOLEAN NOT NULL DEFAULT FALSE)"
	_, err := db.ExecContext(ctx, cr)
	if err != nil {
		return err
	}
	return nil
}

// GetLongByUUID из таблицы получает длинный URL по UUID
func GetLongByUUID(ctx context.Context, db *sql.DB, uuid int) (longURL string, isDeleted bool, err error) {
	row := db.QueryRowContext(ctx, "SELECT original_url, is_deleted FROM url WHERE uuid=$1", uuid)
	err = row.Scan(&longURL, &isDeleted)
	if err != nil {
		return "", false, err
	}
	return
}

// GetShortByOriginalURL из таблицы получает короткий URL по длинному
func GetShortByOriginalURL(ctx context.Context, db *sql.DB, originalURL string) (string, error) {
	row := db.QueryRowContext(ctx, "SELECT uuid FROM url WHERE original_url=$1", originalURL)
	var shURL int
	err := row.Scan(&shURL)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(shURL), err
}

// InsertRecord вставляет длинный URL в таблицу
func InsertRecord(ctx context.Context, db *sql.DB, originalURL, usr string) (string, error) {
	_, err := db.ExecContext(ctx, "insert into url( original_url, usr) values ($1, $2)", originalURL, usr)

	if err != nil {
		return "", err
	}

	shURL, err := GetShortByOriginalURL(ctx, db, originalURL)
	if err != nil {
		return "", err
	}
	return shURL, nil
}

// ListRepoURLs из таблицы получает список длинных URL для определенного пользователя
func ListRepoURLs(ctx context.Context, db *sql.DB, addresBase, usr string) ([]general.ArrRepoURL, error) {
	rows, err := db.QueryContext(ctx, "SELECT uuid, original_url FROM url WHERE usr = $1", usr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	arrRepoURL := make([]general.ArrRepoURL, 0)

	for rows.Next() {
		var repoURL general.ArrRepoURL
		err = rows.Scan(&repoURL.ShortURL, &repoURL.OriginalURL)
		if err != nil {
			return nil, err
		}
		repoURL.ShortURL = addresBase + "/" + repoURL.ShortURL
		arrRepoURL = append(arrRepoURL, repoURL)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return arrRepoURL, nil
}
