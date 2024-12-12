package server

import (
	"io"
	"log"
	"net/http"

	//"github.com/go-chi/chi/middleware"
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

// хандлер для addLink
/*func middleware(next http.Handler) http.Handler {
	start := time.Now()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uri := r.RequestURI
		method := r.Method
		duration := time.Since(start)
		logger.Log.Sugar().Infoln(
			"uri", uri,
			"method", method,
			"duration", duration,
		)
		next.ServeHTTP(w, r)
	})
}*/

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

/*func (s *Server) WithLogging(h http.Handler) http.Handler {
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
}*/

/*func (s *Server) myHandler() http.Handler {

	return http.HandlerFunc(s.addLink)
}

func rootHandle(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Привет"))
}

func (s *Server) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	data := []byte("Привет!")
	res.Write(data)
}*/

func (s *Server) SetupRoutes() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Middleware)
	r.Post("/*", s.addLink)
	r.Get("/{id}", s.GetLongLink)
	return r
	/*r.Post("/*", s.addLink)
	r.Get("/{id}", s.GetLongLink)*/

	/*r.Use(middleware.New(s.WithLogging(func ((w http.ResponseWriter, r *http.Request))  {

	})))*/
	/*r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)*/

}
