package server

import (
	"io"
	"log"
	"net/http"

	//"reductorUrl/internal/app"

	"github.com/developerc/reductorUrl/internal/logger"
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

// хандлер для addLink
func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uri := r.RequestURI
		method := r.Method
		next.ServeHTTP(w, r)
		logger.Log.Sugar().Infoln(
			"uri", uri,
			"method", method,
		)
	})
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

	uri := r.RequestURI
	method := r.Method
	logger.Log.Sugar().Infoln(
		"uri", uri,
		"method", method,
	)

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

func WithLogging(h http.Handler) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		uri := r.RequestURI
		method := r.Method

		// точка, где выполняется хендлер pingHandler
		h.ServeHTTP(w, r) // обслуживание оригинального запроса

		logger.Log.Sugar().Infoln(
			"uri", uri,
			"method", method,
		)
	}
	return http.HandlerFunc(logFn)
}

func (s *Server) SetupRoutes() http.Handler {
	/*rt := http.NewServeMux()
	rt.HandleFunc("/links", s.addLink)*/
	r := chi.NewRouter()
	r.Post("/*", s.addLink)
	//r.Post("/*", middleware(s.addLink))
	//r.Post("/*", WithLogging(s.addLink))
	r.Get("/{id}", s.GetLongLink)
	return r
}
