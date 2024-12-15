package filestorage

import (
	"bufio"
	"encoding/json"
	"os"
)

type Event struct {
	UUID        uint   `json:"uuid"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

type Producer struct {
	file   *os.File
	writer *bufio.Writer
}

type Consumer struct {
	file    *os.File
	scanner *bufio.Scanner
}

var producer Producer
var consumer Consumer

func GetProducer() *Producer {
	return &producer
}

func GetConsumer() *Consumer {
	return &consumer
}

func NewProducer(filename string) error {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	producer = Producer{
		file:   file,
		writer: bufio.NewWriter(file),
	}
	return nil
}

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

func NewConsumer(filename string) error {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	consumer = Consumer{
		file:    file,
		scanner: bufio.NewScanner(file),
	}
	return nil
}

func (c *Consumer) GetEvents() ([]Event, error) {
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
