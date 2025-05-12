// main - пакет для приложения клиента gRPC
// После запуска сервера, когда распечатается лог об успешном старте gRPC сервера, зайти из консоли в папку с файлом
// internal/grpc/client/grpc_client.go
// запустить клиента
// go run .
package main

import (
	"context"
	"log"
	"os"
	"time"

	pb "github.com/developerc/reductorUrl/internal/grpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// main запускает клиента gRPC
func main() {
	var cookieUsr0 string
	addr := "localhost:5000"
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// создадим клиент grpc с перехватчиком
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithUnaryInterceptor(clientInterceptor))
	if err != nil {
		log.Println("could not connect to grpc server: ", err)
		os.Exit(1)
	}
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
			cookieUsr0 = strStrErrResp.CookieValue[5:]
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
			log.Println("error: ", sliceByteErrResp.Err)
		} else {
			log.Println(string(sliceByteErrResp.JsonBytes), sliceByteErrResp.Err)
		}
	}
	// обрабатываем статистику
	sliceByteErrResp, err = grpcClient.GetStatsSvc(ctx, &pb.StrReq{CookieValue: "192.168.0.1"}) //здесь в CookieValue находится IP адрес
	if err != nil {
		log.Println("error request: ", err)
	} else {
		if sliceByteErrResp.Err != "nil" {
			log.Println("error: ", sliceByteErrResp.Err)
		} else {
			log.Println(string(sliceByteErrResp.JsonBytes), sliceByteErrResp.Err)
		}
	}
	// получает URL-ы определенного пользователя
	sliceByteErrResp, err = grpcClient.FetchURLs(ctx, &pb.StrReq{CookieValue: cookieUsr0})
	if err != nil {
		log.Println("error request: ", err)
	} else {
		if sliceByteErrResp.Err != "nil" {
			log.Println("error: ", sliceByteErrResp.Err)
		} else {
			log.Println(string(sliceByteErrResp.JsonBytes), sliceByteErrResp.Err)
		}
	}
	// делает отметку об удалении коротких URL-ы определенного пользователя
	byteJSON = []byte("[\"1\",\"2\"]")
	errMess, err = grpcClient.DelURLs(ctx, &pb.StrByteReq{CookieValue: cookieUsr0, JsonBytes: byteJSON})
	if err != nil {
		log.Println("error request: ", err)
	} else {
		if errMess.Err != "nil" {
			log.Println("error: ", errMess.Err)
		} else {
			log.Println("Список коротких URL принят для отметки о удалении")
		}
	}
}

// clientInterceptor клиентский перехватчик
func clientInterceptor(ctx context.Context, method string, req interface{},
	reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	start := time.Now()

	err := invoker(ctx, method, req, reply, cc, opts...)

	if err != nil {
		log.Printf("[ERROR] %s, %v", method, err)
	} else {
		log.Printf("[INFO] %s, %v", method, time.Since(start))
	}
	return err
}
