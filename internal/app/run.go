package app

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/developerc/reductorUrl/internal/config"
)

type ShortURLAttr struct {
	Settings config.ServerSettings
	Cntr     int
	MapURL   map[int]string
}

var shu ShortURLAttr

func GetShortURLAttr() *ShortURLAttr {
	return &shu
}

func NewShortURLAttr(settings config.ServerSettings) *ShortURLAttr {
	shortURLAttr := ShortURLAttr{}
	shortURLAttr.Settings = settings
	shortURLAttr.MapURL = make(map[int]string)
	return &shortURLAttr
}

func PostHandler(shortURLAttr ShortURLAttr) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		shortURLAttr.Cntr++
		body, err := io.ReadAll(req.Body)
		if err != nil {
			log.Println(err)
			return
		}
		if len(body) == 0 {
			res.WriteHeader(http.StatusBadRequest)
			res.Write([]byte("Base URL not correct"))
			return
		}
		shortURLAttr.MapURL[shortURLAttr.Cntr] = string(body)
		res.Header().Set("Content-Type", "text/plain")
		res.WriteHeader(http.StatusCreated)
		res.Write([]byte(shortURLAttr.Settings.AdresBase + "/" + strconv.Itoa(shortURLAttr.Cntr)))
	}
}

func GetHandler(shortURLAttr ShortURLAttr) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		fmt.Println("id: ", chi.URLParam(req, "id"))
		i, err := strconv.Atoi(chi.URLParam(req, "id"))
		if err != nil {
			log.Println(err)
			return
		}
		longURL := shortURLAttr.MapURL[i]
		res.Header().Set("Location", longURL)
		res.WriteHeader(http.StatusTemporaryRedirect)
	}
}

func BadHandler(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusBadRequest)
	res.Write([]byte("400 StatusBadRequest"))
}

func Run() {
	settings := config.NewServerSettings()
	shu = *NewShortURLAttr(*settings)
	fmt.Println(shu)
	r := chi.NewRouter()
	r.Post("/*", PostHandler(shu))
	r.Get("/{id}", GetHandler(shu))
	r.Put("/*", BadHandler)
	r.Delete("/*", BadHandler)
	r.Options("/*", BadHandler)
	r.Head("/*", BadHandler)
	r.Trace("/*", BadHandler)
	r.Connect("/*", BadHandler)
	r.Patch("/*", BadHandler)

	log.Fatal(http.ListenAndServe(settings.AdresRun, r))
}
