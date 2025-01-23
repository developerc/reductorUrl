package server

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/developerc/reductorUrl/internal/api"
	"github.com/developerc/reductorUrl/internal/logger"
	"github.com/developerc/reductorUrl/internal/middleware"

	"github.com/developerc/reductorUrl/internal/service/memory"
	"github.com/go-chi/chi/v5"
	m "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

type svc interface {
	AddLink(link string, usr string) (string, error)
	Ping() error
	GetLongLink(id string) (string, bool, error)
	HandleBatchJSON(buf bytes.Buffer, usr string) ([]byte, error)
	AsURLExists(err error) bool
	FetchURLs(usr string) ([]byte, error)
	HandleCookie(r *http.Request) (*http.Cookie, string, error)
	DelURLs(r *http.Request) (bool, error)
}

type Server struct {
	service svc
	logger  *zap.Logger
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

	gc, usr, err := s.service.HandleCookie(r)
	if err == nil && gc != nil {
		http.SetCookie(w, gc)
	}

	if shortURL, err = s.service.AddLink(string(body), usr); err != nil {
		if s.service.AsURLExists(err) {
			s.logger.Info("Add link", zap.String("error", err.Error()))
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusConflict)
			if _, err := w.Write([]byte(shortURL)); err != nil {
				return
			}
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
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

	gc, usr, err := s.service.HandleCookie(r)
	if err == nil && gc != nil {
		http.SetCookie(w, gc)
	}

	longURL, err := api.HandleAPIShorten(buf, s.logger)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if shortURL, err = s.service.AddLink(longURL, usr); err != nil {
		if _, ok := err.(*memory.ErrorURLExists); ok {
			s.logger.Info("Add link JSON", zap.String("error", err.Error()))
			jsonBytes, err := api.ShortToJSON(shortURL, s.logger)
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
		return
	}
	if jsonBytes, err = api.ShortToJSON(shortURL, s.logger); err != nil {
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

	gc, usr, err := s.service.HandleCookie(r)
	if err == nil && gc != nil {
		http.SetCookie(w, gc)
	}

	if jsonBytes, err = s.service.HandleBatchJSON(buf, usr); err != nil {
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
	longURL, isDeleted, err := s.service.GetLongLink(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if isDeleted {
		w.WriteHeader(http.StatusGone)
		return
	}
	w.Header().Set("Location", longURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (s *Server) CheckPing(w http.ResponseWriter, r *http.Request) {
	if s.service.Ping() != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func (s *Server) UserURLs(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("user")
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	//if cookie, err := r.Cookie("user"); err == nil {
	//userStruct := &memory.User{}

	//usr, err := s.ReadCookie(r)

	jsonBytes, err := s.service.FetchURLs(cookie.Value)
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	if len(jsonBytes) < 3 {
		w.WriteHeader(http.StatusNoContent)

		return
	} else {
		w.WriteHeader(http.StatusOK)
	}

	if _, err := w.Write(jsonBytes); err != nil {
		return
	}
}

func (s *Server) DelUserURLs(w http.ResponseWriter, r *http.Request) {
	ok, err := s.service.DelURLs(r)
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	if ok {
		w.WriteHeader(http.StatusAccepted)
		if _, err := w.Write([]byte("Accepted")); err != nil {
			return
		}
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		if _, err := w.Write([]byte("Not accepted!")); err != nil {
			return
		}
	}
}

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
	r.Get("/api/user/urls", s.UserURLs)
	r.Delete("/api/user/urls", s.DelUserURLs)
	return r
}
