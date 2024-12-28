package memory

import (
	"log"
	"math"

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

func getDBSettings(shu *ShortURLAttr) {
	// здесь будут настройки для DBStorage
}
