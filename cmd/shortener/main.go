package main

import (
	"fmt"
	"net/http"
)

var strUrl string
var beginUrl string
var cutUrl string

func mainPage(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		if req.Method != http.MethodGet {
			res.WriteHeader(http.StatusBadRequest)
			res.Write([]byte("400 StatusBadRequest"))
			return
		}
		//fmt.Println("handle GET request")
		fmt.Println(req.URL.String())
		if req.URL.String()[1:] == cutUrl {
			res.WriteHeader(http.StatusTemporaryRedirect)
			res.Header().Set("Location", "http://"+beginUrl+"/"+strUrl)
			res.Header().Set("Allow", http.MethodPost)
			res.Write([]byte("redirect"))
		} else {
			res.WriteHeader(http.StatusBadRequest)
			res.Write([]byte("400 StatusBadRequest"))
		}
		return
	}
	beginUrl = req.Host
	fmt.Println("beginUrl: ", beginUrl)
	strUrl = req.URL.String()
	strUrl = strUrl[1:]
	fmt.Println(strUrl)
	/*if url != "/" {
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("400 StatusBadRequest"))
		return
	}*/
	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	//_, _ = res.Write([]byte(`http://localhost:8080/EwHXdJfB `))
	res.Write([]byte("http://" + beginUrl + "/" + cutUrl))
}

/*func apiPage(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("Это страница /api."))
}*/

func main() {
	cutUrl = "EwHXdJfB"
	mux := http.NewServeMux()
	//mux.HandleFunc(`/api/`, apiPage)
	mux.HandleFunc(`/`, mainPage)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
