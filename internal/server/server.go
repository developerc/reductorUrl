package server

import (
	"bytes"
	"io"
	"log"
	"net/http"

	"github.com/developerc/reductorUrl/internal/api"
	"github.com/developerc/reductorUrl/internal/middleware"
	"github.com/go-chi/chi/v5"
)

type svc interface {
	AddLink(link string) (string, error)
}

type Server struct {
	service svc
}

func NewServer(service svc) Server {
	return Server{service: service}
}

func (s *Server) addLink(w http.ResponseWriter, r *http.Request) {
	var shortURL string
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if shortURL, err = s.service.AddLink(string(body)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortURL))
}

func (s *Server) addLinkJSON(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("hello from addLinkJson")
	var buf bytes.Buffer
	var shortURL string
	var jsonBytes []byte
	// читаем тело запроса
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
	if shortURL, err = s.service.AddLink(string(longURL)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if jsonBytes, err = api.ShortToJSON(shortURL); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonBytes)
}

func (s *Server) GetLongLink(w http.ResponseWriter, r *http.Request) {
	log.Println("id: ", chi.URLParam(r, "id"))
	id := chi.URLParam(r, "id")
	longURL, err := GetService().GetLongLink(id)
	//log.Println("longURL", longURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", longURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (s *Server) SetupRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Middleware)
	r.Post("/", s.addLink)
	r.Post("/api/shorten", s.addLinkJSON)
	r.Get("/{id}", s.GetLongLink)
	return r

}
