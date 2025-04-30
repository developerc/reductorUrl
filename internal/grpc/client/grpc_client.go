// main - пакет для приложения клиента gRPC
// После запуска сервера, когда распечатается лог об успешном старте gRPC сервера, зайти из консоли в папку с файлом
// internal/grpc/client/grpc_client.go
// запустить клиента
// go run .
package main

import (
	"context"
	"fmt"

	//"net/http"

	"log"
	"os"
	"time"

	//"github.com/developerc/reductorUrl/internal/general"
	pb "github.com/developerc/reductorUrl/internal/grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure" // для упрощения не будем использовать SSL/TLS аутентификация
)

// main запускает клиента gRPC
func main() {
	//var cookieUsr0 string
	host := "localhost"
	port := "5000"
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	addr := fmt.Sprintf("%s:%s", host, port) // используем адрес сервера
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println("could not connect to grpc server: ", err)
		os.Exit(1)
	}
	// закроем соединение, когда выйдем из функции
	defer conn.Close()

	grpcClient := pb.NewReductorServiceClient(conn)
	// получим куку для user0
	strStrErrResp, err := grpcClient.HandleCookie(ctx, &pb.StrReq{CookieValue: ""})
	if err != nil {
		log.Println("error request: ", err)
	} else {
		if strStrErrResp.Err != "nil" {
			log.Println("error: ", strStrErrResp.Err)
		} else {
			//cookieUsr0 = strStrErrResp.CookieValue
			log.Println(strStrErrResp.CookieValue)
			log.Println(strStrErrResp.Usr)
		}

	}
	// добавляем длинный URL
	shortLink, err := grpcClient.AddLink(ctx, &pb.LinkUsrReq{Link: "http://blabla1.ru", Usr: "user0"})
	if err != nil {
		log.Println("error request: ", err)
	} else {
		if shortLink.Err != "nil" {
			log.Println("error: ", shortLink.Err)
		} else {
			log.Println(shortLink.ShortUrl, shortLink.Err)
		}
	}
	// добавляем длинный URL
	shortLink, err = grpcClient.AddLink(ctx, &pb.LinkUsrReq{Link: "http://blabla2.ru", Usr: "user0"})
	if err != nil {
		log.Println("error request: ", err)
	} else {
		if shortLink.Err != "nil" {
			log.Println("error: ", shortLink.Err)
		} else {
			log.Println(shortLink.ShortUrl, shortLink.Err)
		}
	}
	// проверяем живучесть БД
	errMess, err := grpcClient.Ping(ctx, &pb.StrReq{CookieValue: ""})
	if err != nil {
		log.Println("error request: ", err)
	} else {
		if errMess.Err != "nil" {
			log.Println("БД недоступна")
		} else {
			log.Println("БД доступна")
		}
	}
	// получаем длинный URL по короткому
	longLinkResp, err := grpcClient.GetLongLink(ctx, &pb.IDReq{Id: "1"})
	if err != nil {
		log.Println("error request: ", err)
	} else {
		if longLinkResp.Err != "nil" {
			log.Println("error: ", longLinkResp.Err)
		} else {
			log.Println(longLinkResp.OriginalUrl, longLinkResp.IsDeleted, longLinkResp.Err)
		}
	}
	// добавляем несколько длинных URL
	byteJSON := []byte("[{\"correlation_id\":\"ident1\",\"original_url\":\"http://blabla17.ru\"}]")
	sliceByteErrResp, err := grpcClient.HandleBatchJSON(ctx, &pb.HandleBatchJSONReq{Buf: byteJSON, Usr: "user0"})
	if err != nil {
		log.Println("error request: ", err)
	} else {
		if sliceByteErrResp.Err != "nil" {
			log.Println("error: ", longLinkResp.Err)
		} else {
			log.Println(string(sliceByteErrResp.JsonBytes), sliceByteErrResp.Err)
		}
	}
	// получает URL-ы определенного пользователя !!!!!!!
	/*var usr string
	var cookie *http.Cookie
	u := &general.User{
		Name: usr,
	}
	usr = "user0"
	u.Name = usr
	if encoded, err := s.Secure.Encode("user", u); err == nil {
		cookie = &http.Cookie{
			Name:  "user",
			Value: encoded,
		}
		return cookie, usr, nil
	} else {
		return nil, "", err
	}*/

	sliceByteErrResp, err = grpcClient.FetchURLs(ctx, &pb.StrReq{CookieValue: "MTc0NTk0ODQ4NHxJczFXbGtZbU5ObEFjTjVXZmVNMS1tQ01NOFV6VGJwMDNCbkt0SlZnOFF4dnJzT1lrYlVTWGNvYXBwRllsVDBDWm5TbEV4Y3p8PbePG0Z2PDsQi2gCULhOluhMm2DHSyth3QXt004gTdA"})
	if err != nil {
		log.Println("error request: ", err)
	} else {
		if sliceByteErrResp.Err != "nil" {
			log.Println("error: ", sliceByteErrResp.Err)
		} else {
			log.Println(string(sliceByteErrResp.JsonBytes), sliceByteErrResp.Err)
		}
	}

}
