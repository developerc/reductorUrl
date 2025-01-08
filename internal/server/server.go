package server

import (
	"bytes"
	//"errors"
	"io"
	//"log"
	"net/http"
	"time"

	"github.com/developerc/reductorUrl/internal/api"
	"github.com/developerc/reductorUrl/internal/logger"
	"github.com/developerc/reductorUrl/internal/middleware"

	"github.com/developerc/reductorUrl/internal/service/memory"
	"github.com/go-chi/chi/v5"
	m "github.com/go-chi/chi/v5/middleware"

	//"github.com/jackc/pgerrcode"
	//"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
)

type svc interface {
	AddLink(link string) (string, error)
	Ping() error
	GetLongLink(id string) (string, error)
	HandleBatchJSON(buf bytes.Buffer) ([]byte, error)
	HandleBatchJSONPgx(buf bytes.Buffer) ([]byte, error)
}

type Server struct {
	service svc
	logger  *zap.Logger
	//svc     *memory.Service
}

func NewServer(service svc) (*Server, error) {
	var err error
	srv := new(Server)
	srv.service = service
	srv.logger, err = logger.Initialize("Info")

	if err != nil {
		return srv, err
	}
	return srv, nil
}

func (s *Server) addLink(w http.ResponseWriter, r *http.Request) {
	var shortURL string
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if shortURL, err = s.service.AddLink(string(body)); err != nil {
		//var errURLExits *memory.ErrorURLExists
		if _, ok := err.(*memory.ErrorURLExists); ok {
			s.logger.Info("Add link", zap.String("error", err.Error()))
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusConflict)
			if _, err := w.Write([]byte(shortURL)); err != nil {
				return
			}
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		//if errors.As(err, &memory.ErrorURLExists){}
		/*var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) && pgErr.ConstraintName == "must_be_different" {
			log.Println("Такой оригинальный URL уже существует")
			shortURL, err := s.getService().GetShortByOriginalURL(string(body))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusConflict)
			if _, err := w.Write([]byte(shortURL)); err != nil {
				return
			}
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}*/

		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write([]byte(shortURL)); err != nil {
		return
	}
}

func (s *Server) addLinkJSON(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	var shortURL string
	var jsonBytes []byte

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	longURL, err := api.HandleAPIShorten(buf)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if shortURL, err = s.service.AddLink(longURL); err != nil {
		if _, ok := err.(*memory.ErrorURLExists); ok {
			s.logger.Info("Add link JSON", zap.String("error", err.Error()))
			jsonBytes, err := api.ShortToJSON(shortURL)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			if _, err := w.Write(jsonBytes); err != nil {
				return
			}
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		/*var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) && pgErr.ConstraintName == "must_be_different" {
			log.Println("Такой оригинальный URL уже существует")
			shortURL, err := s.getService().GetShortByOriginalURL(longURL)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			jsonBytes, err := api.ShortToJSON(shortURL)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			if _, err := w.Write(jsonBytes); err != nil {
				return
			}
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}*/
		return
	}
	if jsonBytes, err = api.ShortToJSON(shortURL); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write(jsonBytes); err != nil {
		return
	}
}

func (s *Server) addBatchJSON(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer
	var jsonBytes []byte
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if jsonBytes, err = s.service.HandleBatchJSON(buf); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write(jsonBytes); err != nil {
		return
	}
}

func (s *Server) GetLongLink(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	//longURL, err := s.getService().GetLongLink(id)
	//longURL, err := s.svc.GetLongLink(id)
	longURL, err := s.service.GetLongLink(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", longURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (s *Server) CheckPing(w http.ResponseWriter, r *http.Request) {
	//if s.getService().CheckPing() != nil {
	if s.service.Ping() != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

/*func (s *Server) getService() *memory.Service {
	val := reflect.ValueOf(s.service)
	return (*memory.Service)(val.UnsafePointer())
}*/

func (s *Server) SetupRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Middleware)
	r.Use(middleware.GzipHandle)
	r.Use(m.Timeout(3 * time.Second))
	r.Post("/", s.addLink)
	r.Post("/api/shorten", s.addLinkJSON)
	r.Get("/{id}", s.GetLongLink)
	r.Get("/ping", s.CheckPing)
	r.Post("/api/shorten/batch", s.addBatchJSON)
	return r
}
