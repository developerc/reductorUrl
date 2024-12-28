package memory

import (
	"context"
	"database/sql"
	"log"
	"math"
	"time"

	filestorage "github.com/developerc/reductorUrl/internal/service/file_storage"
)

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
	/*zapLogger, err := logger.Initialize(memory.NewInMemoryService().GetLogLevel())
	if err != nil {
		return err
	}*/
	//service := memory.NewInMemoryService()
	//dsn, err := service.GetDSN()
	//dsn := s.repo.(*ShortURLAttr).Settings.DBStorage
	dsn := shu.Settings.DBStorage
	/*if err != nil {
		zapLogger.Info("CreateTable", zap.String("error", err.Error()))
		return err
	}*/
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		//zapLogger.Info("CreateTable", zap.String("error", err.Error()))
		log.Println(err)
		return err
	}
	defer db.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err = db.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS url_table( uuid serial primary key, short_url INT, original_url TEXT)")
	if err != nil {
		//zapLogger.Info("CreateTable", zap.String("error", err.Error()))
		log.Println(err)
		return err
	}
	log.Println("Table created")
	return nil
}
