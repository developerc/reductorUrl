package memory

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"math"
	"strconv"
	"time"

	"github.com/developerc/reductorUrl/internal/config"
	filestorage "github.com/developerc/reductorUrl/internal/service/file_storage"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type ShortURLAttr struct {
	Settings config.ServerSettings
	Cntr     int
	MapURL   map[int]string
}

type ArrLongURL struct {
	CorellationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type ArrShortURL struct {
	CorellationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

func listLongURL(buf bytes.Buffer) ([]ArrLongURL, error) {
	arrLongURL := make([]ArrLongURL, 0)
	if err := json.Unmarshal(buf.Bytes(), &arrLongURL); err != nil {
		return nil, err
	}
	return arrLongURL, nil
}

func (s *Service) handleArrLongURL(arrLongURL []ArrLongURL) ([]byte, error) {
	arrShortURL := make([]ArrShortURL, 0)
	for _, longURL := range arrLongURL {
		URL, err := s.AddLink(longURL.OriginalURL)
		if err != nil {
			return nil, err
		}
		shortURL := ArrShortURL{CorellationID: longURL.CorellationID, ShortURL: URL}
		arrShortURL = append(arrShortURL, shortURL)
	}
	jsonBytes, err := json.Marshal(arrShortURL)
	if err != nil {
		return nil, err
	}
	return jsonBytes, nil
}

func getFileSettings(shu *ShortURLAttr) error {
	if _, err := filestorage.NewConsumer(shu.Settings.FileStorage); err != nil {
		return err
	}
	consumer, err := filestorage.NewConsumer(shu.Settings.FileStorage)
	if err != nil {
		return err
	}
	events, err := consumer.ListEvents()
	if err != nil {
		return err
	}
	for _, event := range events {
		if event.UUID > math.MaxInt32 {
			event.UUID = math.MaxInt32
		}
		shu.MapURL[int(event.UUID)] = event.OriginalURL
	}
	shu.Cntr = len(events)

	if _, err := filestorage.NewProducer(shu.Settings.FileStorage); err != nil {
		return err
	}
	return nil
}

func createTable(shu *ShortURLAttr) error {
	dsn := shu.Settings.DBStorage
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return err
	}
	defer db.Close()
	const duration uint = 20
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(duration)*time.Second)
	defer cancel()
	const cr string = "CREATE TABLE IF NOT EXISTS url_table( uuid serial primary key, original_url TEXT CONSTRAINT must_be_different UNIQUE)"
	_, err = db.ExecContext(ctx, cr)
	if err != nil {
		return err
	}

	var count int
	if err := db.QueryRowContext(ctx, "SELECT COUNT(*) FROM url_table").Scan(&count); err != nil {
		return err
	}
	shu.Cntr = count
	var rows *sql.Rows
	rows, err = db.QueryContext(ctx, "SELECT uuid, original_url FROM url_table")
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var key int
		var val string
		err = rows.Scan(&key, &val)
		if err != nil {
			return err
		}
		shu.MapURL[key] = val
	}
	err = rows.Err()
	if err != nil {
		return err
	}

	return nil
}

func insertRecord(shu *ShortURLAttr, originalURL string) (string, error) {
	dsn := shu.Settings.DBStorage
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return "", err
	}
	defer db.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err = db.ExecContext(ctx, "insert into url_table( original_url) values ($1)", originalURL)

	if err != nil {
		return "", err
	}

	shURL, err := getShortByOriginalURL(shu, originalURL)
	if err != nil {
		return "", err
	}
	return shURL, nil
}

func getShortByOriginalURL(shu *ShortURLAttr, originalURL string) (string, error) {
	dsn := shu.Settings.DBStorage
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return "", err
	}
	defer db.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	row := db.QueryRowContext(ctx, "SELECT uuid FROM url_table WHERE original_url=$1", originalURL)
	var shURL int
	err = row.Scan(&shURL)
	if err != nil {
		return "", err
	}
	return strconv.Itoa(shURL), err
}

func getLongByUUID(shu *ShortURLAttr, uuid int) (string, error) {
	dsn := shu.Settings.DBStorage
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return "", err
	}
	defer db.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	row := db.QueryRowContext(ctx, "SELECT original_url FROM url_table WHERE uuid=$1", uuid)
	var longURL string
	err = row.Scan(&longURL)
	if err != nil {
		return "", err
	}
	return longURL, nil
}

func (s *Service) CheckPing() error {
	dsn, err := s.GetDSN()
	if err != nil {
		return err
	}
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return err
	}
	defer db.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return err
	}
	return nil
}
