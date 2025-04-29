// main - пакет для приложения клиента gRPC
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	pb "github.com/developerc/reductorUrl/internal/grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure" // для упрощения не будем использовать SSL/TLS аутентификация
)

// main запускает клиента gRPC
func main() {
	host := "localhost"
	port := "5000"

	addr := fmt.Sprintf("%s:%s", host, port) // используем адрес сервера
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println("could not connect to grpc server: ", err)
		os.Exit(1)
	}
	// закроем соединение, когда выйдем из функции
	defer conn.Close()

	grpcClient := pb.NewReductorServiceClient(conn)
	shortLink, err := grpcClient.AddLink(context.TODO(), &pb.LinkUsrReq{Link: "http://long_URL", Usr: "user0"})
	if err != nil {
		log.Println("failed invoking short link: ", err)
	}
	fmt.Println(shortLink.ShortUrl, shortLink.Err)
	//--
	shortLink, err = grpcClient.AddLink(context.TODO(), &pb.LinkUsrReq{Link: "http://long_URL", Usr: "user0"})
	if err != nil {
		log.Println("failed invoking short link: ", err)
	}
	fmt.Println(shortLink.ShortUrl, shortLink.Err)
}
