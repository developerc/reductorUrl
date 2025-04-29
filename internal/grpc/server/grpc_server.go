// server пакет сервера gRPC
package server

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	pb "github.com/developerc/reductorUrl/internal/grpc/proto"
	"github.com/developerc/reductorUrl/internal/service"
	"google.golang.org/grpc"
)

// svc интерфейс с функциями обработки запросов клиента gRPC
type svc interface {
	AddLink(ctx context.Context, link string, usr string) (string, error)
	Ping() error
	GetLongLink(ctx context.Context, id string) (string, bool, error)
	HandleBatchJSON(ctx context.Context, buf bytes.Buffer, usr string) ([]byte, error)
	AsURLExists(err error) bool
	FetchURLs(ctx context.Context, cookieValue string) ([]byte, error)
	HandleCookie(cookieValue string) (*http.Cookie, string, error)
	DelURLs(cookieValue string, buf bytes.Buffer) error
	GetStatsSvc(ctx context.Context, ip net.IP) ([]byte, error)
}

// Server структура сервера gRPC
type Server struct {
	pb.ReductorServiceServer
	Service svc
}

// NewServer конструктор сервера gRPC
func NewServer() *Server {
	return &Server{}
}

// AddLink метод сервера gRPC, добавляет в хранилище длинный URL, возвращает короткий
func (s *Server) AddLink(ctx context.Context, in *pb.LinkUsrReq) (*pb.StrErrResp, error) {
	shortUrl, err := s.Service.AddLink(ctx, in.Link, in.Usr)
	if err != nil {
		return &pb.StrErrResp{ShortUrl: shortUrl, Err: err.Error()}, nil
	} else {
		return &pb.StrErrResp{ShortUrl: shortUrl, Err: "nil"}, nil
	}

}

// NewGRPCserver создает объект структуры, которая содержит реализацию серверной части
func NewGRPCserver(service *service.Service) {
	var host string
	var port string
	//hostAndPort := strings.Split(address, ":")
	hostAndPort := strings.Split(service.Shu.Settings.GRPCaddress, ":")
	host = hostAndPort[0]
	port = hostAndPort[1]

	addr := fmt.Sprintf("%s:%s", host, port)
	lis, err := net.Listen("tcp", addr) // будем ждать запросы по этому адресу

	if err != nil {
		log.Println("error starting tcp listener: ", err)
		os.Exit(1)
	}

	log.Println("tcp listener started at port: ", port)
	// создадим сервер grpc
	grpcServer := grpc.NewServer()

	// объект структуры, которая содержит реализацию
	// серверной части GeometryService
	reductorServiceServer := NewServer()
	reductorServiceServer.Service = service
	// зарегистрируем нашу реализацию сервера
	pb.RegisterReductorServiceServer(grpcServer, reductorServiceServer)
	// запустим grpc сервер
	if err := grpcServer.Serve(lis); err != nil {
		log.Println("error serving grpc: ", err)
		os.Exit(1)
	}
}
