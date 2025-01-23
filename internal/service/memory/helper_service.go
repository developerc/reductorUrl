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

	"github.com/developerc/reductorUrl/internal/config"
	"github.com/developerc/reductorUrl/internal/general"
	dbstorage "github.com/developerc/reductorUrl/internal/service/db_storage"
	filestorage "github.com/developerc/reductorUrl/internal/service/file_storage"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type ShortURLAttr struct {
	Settings config.ServerSettings
	Cntr     int
	MapURL   map[int]string
	MapUser  map[string]bool
	DB       *sql.DB
}

type ArrShortURL struct {
	CorellationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

func (s *Service) HandleCookie(r *http.Request) (*http.Cookie, string, error) {
	var usr string
	var gc *http.Cookie
	var err error
	if s.GetShortURLAttr().Settings.TypeStorage == config.DBStorage {
		_, err = r.Cookie("user")
		if err != nil {
			usr = "user" + strconv.Itoa(s.GetCounter())
			gc, err = s.SetCookie(usr)
			if err != nil {
				return nil, "", err
			}
			s.GetShortURLAttr().MapUser[usr] = true
			return gc, usr, nil
		}

		usr, err = s.ReadCookie(r)
		if err != nil {
			return nil, "", err
		}

		if _, ok := s.GetShortURLAttr().MapUser[usr]; ok {
			return nil, usr, nil
		} else {
			usr = "user" + strconv.Itoa(s.GetCounter())
			gc, err = s.SetCookie(usr)
			if err != nil {
				return nil, "", err
			}
			s.GetShortURLAttr().MapUser[usr] = true
			return gc, usr, nil
		}
	} else {
		return nil, "", nil
	}
}

func CreateMapUser(shu *ShortURLAttr) (map[string]bool, error) {
	mapUser, err := dbstorage.CreateMapUser(shu.DB)
	if err != nil {
		return nil, err
	}
	shu.Cntr = len(mapUser)
	return mapUser, nil
}

func (s *Service) DelURLs(r *http.Request) (bool, error) {
	_, err := r.Cookie("user")
	if err != nil {
		return false, err
	}

	usr, err := s.ReadCookie(r)
	if err != nil {
		return false, http.ErrNoCookie
	}
	if _, ok := s.GetShortURLAttr().MapUser[usr]; !ok {
		return false, http.ErrNoCookie
	}

	var buf bytes.Buffer
	_, err = buf.ReadFrom(r.Body)
	if err != nil {
		return false, err
	}

	arrShortURL := make([]string, 0)
	if err := json.Unmarshal(buf.Bytes(), &arrShortURL); err != nil {
		return false, err
	}

	if err := dbstorage.SetDelBatch(arrShortURL, s.GetShortURLAttr().Settings.DBStorage, usr); err != nil {
		return false, err
	}

	return true, nil
}

func (s *Service) FetchURLs( /*r *http.Request*/ cookieValue string) ([]byte, error) {
	/*_, err := r.Cookie("user")
	if err != nil {
		return nil, err
	}*/

	//usr, err := s.ReadCookie(r)
	/*if err != nil {
		return nil, http.ErrNoCookie
	}*/
	//usr, err := s.ReadCookie2(cookieValue)
	var err error
	usr := &User{}
	if err = secure.Decode("user", cookieValue, usr); err != nil {
		return nil, err
	}

	if _, ok := s.GetShortURLAttr().MapUser[usr.Name]; !ok {
		return nil, http.ErrNoCookie
	}

	arrRepoURL, err := dbstorage.ListRepoURLs(s.GetShortURLAttr().DB, s.GetAdresBase(), usr.Name)
	if err != nil {
		return nil, err
	}

	jsonBytes, err := json.Marshal(arrRepoURL)
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

func (s *Service) handleArrLongURL(arrLongURL []general.ArrLongURL, usr string) ([]byte, error) {
	shu := s.GetShortURLAttr()
	if shu.Settings.TypeStorage != config.DBStorage {
		arrShortURL := make([]ArrShortURL, 0)
		for _, longURL := range arrLongURL {
			URL, err := s.AddLink(longURL.OriginalURL, usr)
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

	if err := dbstorage.InsertBatch(arrLongURL, shu.Settings.DBStorage, usr); err != nil {
		return nil, err
	}

	arrShortURL := make([]ArrShortURL, 0)
	for _, longURL := range arrLongURL {
		short, err := dbstorage.GetShortByOriginalURL(shu.DB, longURL.OriginalURL)
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
		shu.MapURL[int(event.UUID)] = event.OriginalURL
	}
	shu.Cntr = len(events)

	if _, err := filestorage.NewProducer(shu.Settings.FileStorage); err != nil {
		return err
	}
	return nil
}

func (s *Service) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := s.GetShortURLAttr().DB.PingContext(ctx); err != nil {
		return err
	}
	return nil
}
