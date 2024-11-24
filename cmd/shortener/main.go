package main

import (
	"io"
	"log"
	"net/http"
	"strconv"
)

type ShortURLAttr struct {
	BeginURL string
	Cntr     int
	MapURL   map[int]string
}

var shortURLAttr ShortURLAttr

func mainPage2(res http.ResponseWriter, req *http.Request) {

	switch req.Method {
	case http.MethodPost:
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
		shortURLAttr.Cntr++
		return
	case http.MethodGet:
		i, err := strconv.Atoi(req.URL.String()[1:])
		if err != nil {
			log.Println(err)
		}
		longURL := shortURLAttr.MapURL[i]
		res.Header().Set("Location", longURL)
		res.WriteHeader(http.StatusTemporaryRedirect)
		return

	default:
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("400 StatusBadRequest"))
	}
}

func main() {
	shortURLAttr = ShortURLAttr{}
	shortURLAttr.MapURL = make(map[int]string)
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, mainPage2)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
