// filestorage пакет для размещения методов хранения данных в файловой системе.
package filestorage

import (
	"bufio"
	"encoding/json"
	"io/fs"
	"os"
)

// Event структура для хранения объекта
type Event struct {
	UUID        uint   `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// EventWriter интерфейс записи объекта в файл
type EventWriter interface {
	WriteEvent(event *Event) error
}

// EventReader интерфейс чтения списка объектов из файла
type EventReader interface {
	ListEvents() ([]Event, error)
}

// Producer структура для записи объекта в файл
type Producer struct {
	evwr   EventWriter
	file   *os.File
	writer *bufio.Writer
}

// Consumer структура для чтения из файла
type Consumer struct {
	evre    EventReader
	file    *os.File
	scanner *bufio.Scanner
}

var producer Producer
var consumer Consumer

// NewProducer конструктор продюсера
func NewProducer(filename string) (*Producer, error) {
	if producer.evwr != nil {
		return &producer, nil
	}
	const filePermission fs.FileMode = 0666
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, filePermission)
	if err != nil {
		return nil, err
	}

	producer = Producer{
		file:   file,
		writer: bufio.NewWriter(file),
	}
	return &producer, err
}

// WriteEvent метод записывает в файл
func (p *Producer) WriteEvent(event *Event) error {
	data, err := json.Marshal(&event)
	if err != nil {
		return err
	}

	if _, err := p.writer.Write(data); err != nil {
		return err
	}

	if err := p.writer.WriteByte('\n'); err != nil {
		return err
	}

	return p.writer.Flush()
}

// NewConsumer конструктор консумера
func NewConsumer(filename string) (*Consumer, error) {
	if consumer.evre != nil {
		return &consumer, nil
	}
	const filePermission fs.FileMode = 0666
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, filePermission)
	if err != nil {
		return nil, err
	}

	consumer = Consumer{
		file:    file,
		scanner: bufio.NewScanner(file),
	}
	return &consumer, nil
}

// ListEvents метод читает список объектов из файла
func (c *Consumer) ListEvents() ([]Event, error) {
	events := make([]Event, 0)
	c.scanner.Split(bufio.ScanLines)
	for c.scanner.Scan() {
		data := c.scanner.Bytes()
		event := Event{}
		err := json.Unmarshal(data, &event)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	c.file.Close()
	return events, nil
}

func (c *Consumer) Close() error {
	return c.file.Close()
}
