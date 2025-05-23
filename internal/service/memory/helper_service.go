package memory

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"math"
	"net/http"
	"strconv"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/developerc/reductorUrl/internal/config"
	"github.com/developerc/reductorUrl/internal/general"
	dbstorage "github.com/developerc/reductorUrl/internal/service/db_storage"
	filestorage "github.com/developerc/reductorUrl/internal/service/file_storage"
)

// MapURLVal структура для значения map MapURL
type MapURLVal struct {
	OriginalURL string
	Usr         string
	IsDeleted   string
}

// ShortURLAttr структура аттрибутов коротких URL
type ShortURLAttr struct {
	MapURL   map[int]MapURLVal
	MapUser  map[string]bool
	DB       *sql.DB
	Settings config.ServerSettings
	Cntr     int
}

// ArrShortURL структура массива коротких URL
type ArrShortURL struct {
	CorellationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

// User структура пользователя
type User struct {
	Name string
}

// HandleCookie метод для работы с куками
func (s *Service) HandleCookie(cookieValue string) (*http.Cookie, string, error) {
	var usr string
	var cookie *http.Cookie
	u := &User{
		Name: usr,
	}

	if cookieValue == "" {
		usr = "user" + strconv.Itoa(s.GetCounter())
		u.Name = usr
		if encoded, err := s.secure.Encode("user", u); err == nil {
			cookie = &http.Cookie{
				Name:  "user",
				Value: encoded,
			}
			return cookie, usr, nil
		} else {
			return nil, "", err
		}
	}
	if err := s.secure.Decode("user", cookieValue, u); err != nil {
		return nil, "", err
	}
	s.mu.RLock()
	_, ok := s.shu.MapUser[u.Name]
	s.mu.RUnlock()
	if ok {
		return nil, u.Name, nil
	} else {
		usr = "user" + strconv.Itoa(s.GetCounter())
		u.Name = usr
		if encoded, err := s.secure.Encode("user", u); err == nil {
			cookie = &http.Cookie{
				Name:  "user",
				Value: encoded,
			}
			s.mu.Lock()
			s.shu.MapUser[usr] = true
			s.mu.Unlock()
			return cookie, usr, nil
		} else {
			return nil, "", err
		}
	}
}

// CreateMapUser создает Map пользователей
func CreateMapUser(ctx context.Context, shu *ShortURLAttr) (map[string]bool, error) {
	mapUser, err := dbstorage.CreateMapUser(ctx, shu.DB)
	if err != nil {
		return nil, err
	}
	shu.Cntr = len(mapUser)
	return mapUser, nil
}

func (s *Service) setDelMemory(arrShortURL []string, usr string) error {
	var err error

	for _, shortURL := range arrShortURL {
		intShortURL, err := strconv.Atoi(shortURL)
		if err != nil {
			return err
		}
		if val, ok := s.shu.MapURL[intShortURL]; ok {
			if val.Usr == usr {
				s.mu.Lock()
				mapURLVal := s.shu.MapURL[intShortURL]
				mapURLVal.IsDeleted = "true"
				s.shu.MapURL[intShortURL] = mapURLVal
				s.mu.Unlock()
			}
		}
	}
	if s.shu.Settings.TypeStorage == config.FileStorage {
		s.changeFileStorage()
	}
	return err
}

// DelURLs делает отметку об удалении коротких URL-ы определенного пользователя
func (s *Service) DelURLs(cookieValue string, buf bytes.Buffer) error {
	general.CntrAtomVar.IncrCntr()
	u := &User{}
	if err := s.secure.Decode("user", cookieValue, u); err != nil {
		return err
	}

	s.mu.RLock()
	_, ok := s.shu.MapUser[u.Name]
	s.mu.RUnlock()
	if !ok {
		return http.ErrNoCookie
	}

	arrShortURL := make([]string, 0)
	if err := json.Unmarshal(buf.Bytes(), &arrShortURL); err != nil {
		return err
	}

	if s.shu.Settings.TypeStorage != config.DBStorage {
		if err := s.setDelMemory(arrShortURL, u.Name); err != nil {
			return err
		}
	} else {
		if err := dbstorage.SetDelBatch(arrShortURL, s.shu.DB, u.Name); err != nil {
			return err
		}
	}

	general.CntrAtomVar.DecrCntr()
	general.CntrAtomVar.SentNotif()
	return nil
}

// listURLsMemory для определенного пользователя получает список пар короткий URL, длинный URL
func (s *Service) listURLsMemory(usr string) ([]general.ArrRepoURL, error) {
	arrRepoURL := make([]general.ArrRepoURL, 0)
	s.mu.RLock()
	defer s.mu.RUnlock()
	for uuid, val := range s.shu.MapURL {
		if val.Usr == usr {
			repoURL := general.ArrRepoURL{}
			repoURL.ShortURL = s.shu.Settings.AdresBase + "/" + strconv.Itoa(uuid)
			repoURL.OriginalURL = val.OriginalURL
			arrRepoURL = append(arrRepoURL, repoURL)
		}
	}
	return arrRepoURL, nil
}

// FetchURLs получает URL-ы определенного пользователя
func (s *Service) FetchURLs(ctx context.Context, cookieValue string) ([]byte, error) {
	u := &User{}
	if err := s.secure.Decode("user", cookieValue, u); err != nil {
		return nil, err
	}

	s.mu.RLock()
	_, ok := s.shu.MapUser[u.Name]
	s.mu.RUnlock()
	if !ok {
		return nil, http.ErrNoCookie
	}
	var jsonBytes []byte
	var arrRepoURL []general.ArrRepoURL
	var err error
	if s.shu.Settings.TypeStorage != config.DBStorage {
		arrRepoURL, err = s.listURLsMemory(u.Name)
		if err != nil {
			return nil, err
		}
	} else {
		arrRepoURL, err = dbstorage.ListRepoURLs(ctx, s.shu.DB, s.GetAdresBase(), u.Name)
		if err != nil {
			return nil, err
		}
	}

	jsonBytes, err = json.Marshal(arrRepoURL)
	if err != nil {
		return nil, err
	}
	return jsonBytes, nil
}

func listLongURL(buf bytes.Buffer) ([]general.ArrLongURL, error) {
	arrLongURL := make([]general.ArrLongURL, 0)
	if err := json.Unmarshal(buf.Bytes(), &arrLongURL); err != nil {
		return nil, err
	}
	return arrLongURL, nil
}

func (s *Service) handleArrLongURL(ctx context.Context, arrLongURL []general.ArrLongURL, usr string) ([]byte, error) {
	shu := s.shu
	if shu.Settings.TypeStorage != config.DBStorage {
		arrShortURL := make([]ArrShortURL, 0)
		for _, longURL := range arrLongURL {
			URL, err := s.AddLink(ctx, longURL.OriginalURL, usr)
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

	if err := dbstorage.InsertBatch2(ctx, arrLongURL, shu.DB, usr); err != nil {
		return nil, err
	}

	arrShortURL := make([]ArrShortURL, 0)
	for _, longURL := range arrLongURL {
		short, err := dbstorage.GetShortByOriginalURL(ctx, shu.DB, longURL.OriginalURL)
		if err != nil {
			return nil, err
		}
		shortURL := ArrShortURL{CorellationID: longURL.CorellationID, ShortURL: s.GetAdresBase() + "/" + short}
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

		shu.MapURL[int(event.UUID)] = MapURLVal{OriginalURL: event.OriginalURL, Usr: event.Usr, IsDeleted: event.IsDeleted}
		shu.MapUser[shu.MapURL[int(event.UUID)].Usr] = true
	}
	shu.Cntr = len(events)

	if _, err := filestorage.NewProducer(shu.Settings.FileStorage); err != nil {
		return err
	}
	return nil
}

// Ping делает проверку живучести БД
func (s *Service) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := s.shu.DB.PingContext(ctx); err != nil {
		return err
	}
	return nil
}
