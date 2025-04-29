// main - пакет для приложения клиента gRPC
// После запуска сервера, когда распечатается лог об успешном старте gRPC сервера, зайти из консоли в папку с файлом
// internal/grpc/client/grpc_client.go
// запустить клиента
// go run .
package main

import (
	"context"
	"fmt"

	"log"
	"os"
	"time"

	pb "github.com/developerc/reductorUrl/internal/grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure" // для упрощения не будем использовать SSL/TLS аутентификация
)

// main запускает клиента gRPC
func main() {
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
	//--
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
	//--
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
	//--
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
	//--
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
}
