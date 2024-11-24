package main

import (
	"io"
	"log"
	"net/http"
	"strconv"
)

type ShortUrlAttr struct {
	BeginUrl string
	Cntr     int
	MapUrl   map[int]string
}

var shortUrlAttr ShortUrlAttr

func mainPage2(res http.ResponseWriter, req *http.Request) {

	switch req.Method {
	case http.MethodPost:
		body, err := io.ReadAll(req.Body)
		if err != nil {
			log.Println(err)
			return
		}
		shortUrlAttr.MapUrl[shortUrlAttr.Cntr] = string(body)
		shortUrlAttr.BeginUrl = req.Host
		res.Header().Set("Content-Type", "text/plain")
		res.WriteHeader(http.StatusCreated)
		res.Write([]byte("http://" + shortUrlAttr.BeginUrl + "/" + strconv.Itoa(shortUrlAttr.Cntr)))
		shortUrlAttr.Cntr++
		return
	case http.MethodGet:
		i, err := strconv.Atoi(req.URL.String()[1:])
		if err != nil {
			log.Println(err)
		}
		longUrl := shortUrlAttr.MapUrl[i]
		res.Header().Set("Location", longUrl)
		res.WriteHeader(http.StatusTemporaryRedirect)
		return

	default:
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("400 StatusBadRequest"))
	}
}

func main() {
	shortUrlAttr = ShortUrlAttr{}
	shortUrlAttr.MapUrl = make(map[int]string)
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, mainPage2)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
