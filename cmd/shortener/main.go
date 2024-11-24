package main

import (
	//"fmt"
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

// var strUrl string
//var beginUrl string

// var cutUrl string
// var cntr int
//var mapUrl map[int]string

func mainPage2(res http.ResponseWriter, req *http.Request) {

	switch req.Method {
	case http.MethodPost:
		body, err := io.ReadAll(req.Body)
		if err != nil {
			//panic(err)
			log.Println(err)
			return
		}
		//fmt.Println(string(body))
		//mapUrl[cntr] = string(body)
		//mapUrl[shortUrlAttr.Cntr] = string(body)
		shortUrlAttr.MapUrl[shortUrlAttr.Cntr] = string(body)

		//fmt.Println(mapUrl)
		//beginUrl = req.Host
		shortUrlAttr.BeginUrl = req.Host
		//fmt.Println("beginUrl: ", beginUrl)
		//strUrl = req.URL.String()
		//strUrl = strUrl[1:]
		//fmt.Println("strUrl: " + strUrl)
		res.Header().Set("Content-Type", "text/plain")
		res.WriteHeader(http.StatusCreated)
		//_, _ = res.Write([]byte(`http://localhost:8080/EwHXdJfB `))
		//res.Write([]byte("http://" + beginUrl + "/" + strconv.Itoa(cntr)))
		//res.Write([]byte("http://" + shortUrlAttr.BeginUrl + "/" + strconv.Itoa(cntr)))
		res.Write([]byte("http://" + shortUrlAttr.BeginUrl + "/" + strconv.Itoa(shortUrlAttr.Cntr)))
		//cntr++
		shortUrlAttr.Cntr++
		return
	case http.MethodGet:
		//tmpStr := req.URL.String()[1:]
		//fmt.Println(tmpStr)
		i, err := strconv.Atoi(req.URL.String()[1:])
		if err != nil {
			// ... handle error
			//panic(err)
			log.Println(err)
		}
		//longUrl := mapUrl[i]
		longUrl := shortUrlAttr.MapUrl[i]
		//res.Header().Set("Location", "http://"+beginUrl+"/"+strUrl)
		res.Header().Set("Location", longUrl)
		res.WriteHeader(http.StatusTemporaryRedirect)
		//delete(mapUrl, i)
		return
		/*tmpStr := req.URL.String()[1:]
		//fmt.Println(tmpStr + "   " + cutUrl)
		if tmpStr == cutUrl {

			res.Header().Set("Location", "http://"+beginUrl+"/"+strUrl)
			//res.Header().Set("Allow", http.MethodPost)
			res.WriteHeader(http.StatusTemporaryRedirect)
			//res.Write([]byte("redirect"))
			return
		} else {
			res.WriteHeader(http.StatusBadRequest)
			res.Write([]byte("400 StatusBadRequest"))
			return
		}*/

	default:
		res.WriteHeader(http.StatusBadRequest)
		res.Write([]byte("400 StatusBadRequest"))
	}
}

func main() {
	shortUrlAttr = ShortUrlAttr{}
	shortUrlAttr.MapUrl = make(map[int]string)
	//cntr = 0
	//mapUrl = make(map[int]string)
	//cutUrl = "EwHXdJfB"
	mux := http.NewServeMux()
	//mux.HandleFunc(`/api/`, apiPage)
	mux.HandleFunc(`/`, mainPage2)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
