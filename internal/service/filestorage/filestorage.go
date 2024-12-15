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
	file *os.File
	// добавляем Writer в Producer
	writer *bufio.Writer
}

type Consumer struct {
	file *os.File
	// заменяем Reader на Scanner
	scanner *bufio.Scanner
}

var producer Producer
var consumer Consumer

/*func InitStorage(filename string) error {
	if err := NewProducer(filename); err != nil {
		return err
	}

	return nil
}*/

func GetProducer() *Producer {
	return &producer
}

func GetConsumer() *Consumer {
	return &consumer
}

// -- Producer
func NewProducer(filename string) error /*(*Producer, error)*/ {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	//defer file.Close()
	//fmt.Println(file.Name())

	producer = Producer{
		file: file,
		// создаём новый Writer
		writer: bufio.NewWriter(file),
	}
	return nil
	/*return &Producer{
		file: file,
		// создаём новый Writer
		writer: bufio.NewWriter(file),
	}, nil*/
}

func (p *Producer) WriteEvent(event *Event) error {
	/*_, err := os.OpenFile(p.file.Name(), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}*/
	//defer file.Close()
	data, err := json.Marshal(&event)
	if err != nil {
		return err
	}

	// записываем событие в буфер
	if _, err := p.writer.Write(data); err != nil {
		return err
	}

	// добавляем перенос строки
	if err := p.writer.WriteByte('\n'); err != nil {
		return err
	}

	// записываем буфер в файл
	return p.writer.Flush()
}

// --Consumer
func NewConsumer(filename string) error /*(*Consumer, error)*/ {
	file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}

	consumer = Consumer{
		file: file,
		// создаём новый scanner
		scanner: bufio.NewScanner(file),
	}
	return nil
	/*return &Consumer{
		file: file,
		// создаём новый scanner
		scanner: bufio.NewScanner(file),
	}, nil*/
}

/*func (c *Consumer) ReadEvent() (*Event, error) {
	// одиночное сканирование до следующей строки
	if !c.scanner.Scan() {
		return nil, c.scanner.Err()
	}
	// читаем данные из scanner
	data := c.scanner.Bytes()

	event := Event{}
	err := json.Unmarshal(data, &event)
	if err != nil {
		return nil, err
	}

	return &event, nil
}*/

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

/*func (c *Consumer) InitFillMap() (map[int]string, error) {
	var mapURL map[int]string = make(map[int]string)
	c.scanner.Split(bufio.ScanLines)
	for c.scanner.Scan() {
		data := c.scanner.Bytes()
		//fmt.Println(string(data))
		event := Event{}
		err := json.Unmarshal(data, &event)
		if err != nil {
			return nil, err
		}
		//mapURL[1]=""
		mapURL[int(event.Uuid)] = event.Original_url

	}
	fmt.Println(mapURL)
	c.Close()
	return mapURL, nil
}*/

func (c *Consumer) Close() error {
	return c.file.Close()
}
