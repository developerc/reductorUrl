package server

import (
	"bytes"
	"fmt"
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
	AddLink(link string) (string, error)
	Ping() error
	GetLongLink(id string) (string, error)
	HandleBatchJSON(buf bytes.Buffer) ([]byte, error)
	AsURLExists(err error) bool
	GetCripto() (string, error)
	FetchURLs() ([]byte, error)
	IsRegisteredCookie(cookie string) bool
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
	fmt.Println("from addLink")
	existCookie := true
	var shortURL string
	var cripto string
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//здесь проверяем пришла ли кука
	cookie, err := r.Cookie("idUser")
	if err != nil {
		fmt.Println(err)
		existCookie = false
	} else {
		if s.service.IsRegisteredCookie(cookie.Value) == false { //проверим есть ли такая кука в списке
			existCookie = false
		}
	}

	if existCookie { //если пришла, берем эту куку
		cripto = cookie.Value
	} else { //если нет, генерируем куку
		cripto, err = s.service.GetCripto()
		fmt.Println("from UserURLs cripto: ", cripto)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if shortURL, err = s.service.AddLink(string(body)); err != nil {
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

	//если кука сгенерирована, добавляем куку
	if !existCookie {
		c := http.Cookie{
			Name:  "idUser",
			Value: cripto,
		}
		http.SetCookie(w, &c)
	}
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write([]byte(shortURL)); err != nil {
		return
	}
}

func (s *Server) addLinkJSON(w http.ResponseWriter, r *http.Request) {
	fmt.Println("from addLinkJSON")
	var buf bytes.Buffer
	var shortURL string
	var jsonBytes []byte

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	longURL, err := api.HandleAPIShorten(buf, s.logger)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if shortURL, err = s.service.AddLink(longURL); err != nil {
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
	fmt.Println("from addBatchJSON")
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
	longURL, err := s.service.GetLongLink(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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
	jsonBytes, err := s.service.FetchURLs()
	fmt.Println("len(jsonBytes): ", len(jsonBytes))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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
	/*cookie, err := r.Cookie("exampleCookie")
	if err != nil {
		fmt.Println("no cookies!")
		cripto, err := s.service.GetCripto()
		fmt.Println("from UserURLs cripto: ", cripto)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		c := http.Cookie{
			Name:  "exampleCookie",
			Value: string(cripto),
		}
		http.SetCookie(w, &c)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	fmt.Println("cookie.Value: " + cookie.Value)
	w.WriteHeader(http.StatusOK)*/
	//return
	/*fmt.Println("from UserURLs")
	cookie, err := r.Cookie("exampleCookie")
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			http.Error(w, "cookie not found", http.StatusBadRequest)
			cookie := http.Cookie{
				Name:     "exampleCookie",
				Value:    "Hello world!",
				Path:     "/",
				MaxAge:   3600,
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteLaxMode,
			}
			http.SetCookie(w, &cookie)
			w.Write([]byte(cookie.Value))
			return
		default:
			log.Println(err)
			http.Error(w, "server error", http.StatusInternalServerError)
		}
		return
	}

	// Echo out the cookie value in the response body.
	w.Write([]byte(cookie.Value))*/
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
	return r
}
