package main

import (
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type ShortURLAttr struct {
	BeginURL string
	Cntr     int
	MapURL   map[int]string
}

var shortURLAttr ShortURLAttr

func PostHandler(res http.ResponseWriter, req *http.Request) {
	if shortURLAttr.MapURL == nil {
		shortURLAttr = ShortURLAttr{}
		shortURLAttr.MapURL = make(map[int]string)
	}
	shortURLAttr.Cntr++
	body, err := io.ReadAll(req.Body)
	if err != nil {
		log.Println(err)
		return
	}
	shortURLAttr.MapURL[shortURLAttr.Cntr] = string(body)
	shortURLAttr.BeginURL = req.Host
	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte("http://" + shortURLAttr.BeginURL + "/" + strconv.Itoa(shortURLAttr.Cntr)))
}

func GetHandler(res http.ResponseWriter, req *http.Request) {
	if shortURLAttr.MapURL == nil {
		shortURLAttr = ShortURLAttr{}
		shortURLAttr.MapURL = make(map[int]string)
	}
	i, err := strconv.Atoi(chi.URLParam(req, "id"))
	if err != nil {
		log.Println(err)
	}
	longURL := shortURLAttr.MapURL[i]
	res.Header().Set("Location", longURL)
	res.WriteHeader(http.StatusTemporaryRedirect)
}

func BadHandler(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusBadRequest)
	res.Write([]byte("400 StatusBadRequest"))
}

func main() {
	r := chi.NewRouter()
	r.Post("/*", PostHandler)
	r.Get("/{id}", GetHandler)
	r.Put("/*", BadHandler)
	r.Delete("/*", BadHandler)
	r.Options("/*", BadHandler)
	r.Head("/*", BadHandler)
	r.Trace("/*", BadHandler)
	r.Connect("/*", BadHandler)
	r.Patch("/*", BadHandler)

	log.Fatal(http.ListenAndServe(":8080", r))
}
