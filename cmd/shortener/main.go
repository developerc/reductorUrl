package main

import (
	//"fmt"
	"net/http"
)

var strUrl string
var beginUrl string
var cutUrl string

func mainPage2(res http.ResponseWriter, req *http.Request) {

	switch req.Method {
	case http.MethodPost:
		beginUrl = req.Host
		//fmt.Println("beginUrl: ", beginUrl)
		strUrl = req.URL.String()
		strUrl = strUrl[1:]
		//fmt.Println("strUrl: " + strUrl)
		res.Header().Set("Content-Type", "text/plain")
		res.WriteHeader(http.StatusCreated)
		//_, _ = res.Write([]byte(`http://localhost:8080/EwHXdJfB `))
		res.Write([]byte("http://" + beginUrl + "/" + cutUrl))
		return
	case http.MethodGet:
		res.Header().Set("Location", "http://"+beginUrl+"/"+strUrl)
		res.WriteHeader(http.StatusTemporaryRedirect)
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

/*func mainPage(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		if req.Method != http.MethodGet {
			res.WriteHeader(http.StatusBadRequest)
			res.Write([]byte("400 StatusBadRequest"))
			return
		}
		//fmt.Println("handle GET request")
		fmt.Println(req.URL.String())
		if req.URL.String()[1:] == cutUrl {

			res.Header().Set("Location", "http://"+beginUrl+"/"+strUrl)
			res.Header().Set("Allow", http.MethodPost)
			res.WriteHeader(http.StatusTemporaryRedirect)
			//res.Write([]byte("redirect"))
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

	res.Header().Set("Content-Type", "text/plain")
	res.WriteHeader(http.StatusCreated)
	//_, _ = res.Write([]byte(`http://localhost:8080/EwHXdJfB `))
	res.Write([]byte("http://" + beginUrl + "/" + cutUrl))
}*/

/*func apiPage(res http.ResponseWriter, req *http.Request) {
	res.Write([]byte("Это страница /api."))
}*/

func main() {
	cutUrl = "EwHXdJfB"
	mux := http.NewServeMux()
	//mux.HandleFunc(`/api/`, apiPage)
	mux.HandleFunc(`/`, mainPage2)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
