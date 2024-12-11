package server

import (
	"io"
	"log"
	"net/http"

	//"reductorUrl/internal/app"

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
	//link := r.FormValue("link")
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

func (s *Server) GetLongLink(w http.ResponseWriter, r *http.Request) {
	log.Println("id: ", chi.URLParam(r, "id"))
	id := chi.URLParam(r, "id")
	longURL, err := GetService().GetLongLink(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Location", longURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (s *Server) SetupRoutes() http.Handler {
	/*rt := http.NewServeMux()
	rt.HandleFunc("/links", s.addLink)*/
	r := chi.NewRouter()
	r.Post("/*", s.addLink)
	r.Get("/{id}", s.GetLongLink)
	return r
}
