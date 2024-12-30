package memory

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"math"
	"time"

	filestorage "github.com/developerc/reductorUrl/internal/service/file_storage"
)

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
		//zapLogger, err := logger.Initialize(memory.NewInMemoryService().GetLogLevel())
		//zapLogger.Info("HandleBatchJSON", zap.String("error", "demarshalling"))

		return nil, err
	}
	//fmt.Println(arrLongURL)
	return arrLongURL, nil
}

func (s Service) handleArrLongURL(arrLongURL []ArrLongURL) ([]byte, error) {
	arrShortURL := make([]ArrShortURL, 0)
	for _, longURL := range arrLongURL {
		URL, err := s.AddLink(longURL.OriginalURL)
		if err != nil {
			return nil, err
		}
		shortURL := ArrShortURL{CorellationID: longURL.CorellationID, ShortURL: URL}
		arrShortURL = append(arrShortURL, shortURL)
	}
	//fmt.Println(arrShortURL)
	jsonBytes, err := json.Marshal(arrShortURL)
	if err != nil {
		return nil, err
	}
	//fmt.Println(string(jsonBytes))
	return jsonBytes, nil
}

func getFileSettings(shu *ShortURLAttr) {
	if _, err := filestorage.NewConsumer(shu.Settings.FileStorage); err != nil {
		log.Println(err)
	}
	consumer, err := filestorage.NewConsumer(shu.Settings.FileStorage)
	if err != nil {
		log.Println(err)
	}
	events, err := consumer.ListEvents()
	if err != nil {
		log.Println(err)
	}
	for _, event := range events {
		if event.UUID > math.MaxInt32 {
			event.UUID = math.MaxInt32
		}
		shu.MapURL[int(event.UUID)] = event.OriginalURL
	}
	shu.Cntr = len(events)

	if _, err := filestorage.NewProducer(shu.Settings.FileStorage); err != nil {
		log.Println(err)
	}
}

func createTable(shu *ShortURLAttr) error {
	dsn := shu.Settings.DBStorage
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Println(err)
		return err
	}
	defer db.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	_, err = db.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS url_table( uuid serial primary key, short_url INT, original_url TEXT)")
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("Table created")

	var count int
	db.QueryRowContext(ctx, "SELECT COUNT(*) FROM url_table").Scan(&count)
	shu.Cntr = count
	//log.Println("Cntr: ", shu.Cntr)
	var rows *sql.Rows
	rows, err = db.QueryContext(ctx, "SELECT short_url, original_url FROM url_table")
	if err != nil {
		log.Println(err)
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var key int
		var val string
		err = rows.Scan(&key, &val)
		if err != nil {
			log.Println(err)
			return err
		}
		//log.Println(key, val)
		shu.MapURL[key] = val
	}
	//log.Println(shu.MapURL)
	err = rows.Err()
	if err != nil {
		return err
	}

	return nil
}

func insertRecord(shu *ShortURLAttr, originalURL string) error {
	//shu.Cntr++
	shu.MapURL[shu.Cntr] = originalURL
	dsn := shu.Settings.DBStorage
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Println(err)
		return err
	}
	defer db.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err = db.ExecContext(ctx, "insert into url_table(short_url, original_url) values ($1, $2)", shu.Cntr, originalURL)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("Record inserted")
	return nil
}
