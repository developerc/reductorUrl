// server пакет сервера gRPC
package server

import (
	"bytes"
	"context"
	"fmt"

	//"fmt"

	"net"
	"net/http"

	///"os"

	pb "github.com/developerc/reductorUrl/internal/grpc/proto"
	"github.com/developerc/reductorUrl/internal/service"
	"go.uber.org/zap"
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
	shortURL, err := s.Service.AddLink(ctx, in.Link, in.Usr)
	if err != nil {
		return &pb.StrErrResp{ShortUrl: shortURL, Err: err.Error()}, err
	} else {
		return &pb.StrErrResp{ShortUrl: shortURL, Err: "nil"}, nil
	}
}

// Ping проверяет живучесть БД
func (s *Server) Ping(ctx context.Context, in *pb.StrReq) (*pb.ErrMess, error) {
	err := s.Service.Ping()
	if err != nil {
		return &pb.ErrMess{Err: err.Error()}, err
	}
	return &pb.ErrMess{Err: "nil"}, nil
}

// GetLongLink получает длинный URL по короткому
func (s *Server) GetLongLink(ctx context.Context, in *pb.IDReq) (*pb.LongLinkResp, error) {
	originalURL, isDeleted, err := s.Service.GetLongLink(ctx, in.Id)
	if err != nil {
		return nil, err
	}
	return &pb.LongLinkResp{OriginalUrl: originalURL, IsDeleted: isDeleted, Err: "nil"}, nil
}

// HandleBatchJSON добавляет в хранилище несколько длинных URL
func (s *Server) HandleBatchJSON(ctx context.Context, in *pb.HandleBatchJSONReq) (*pb.SliceByteErrResp, error) {
	var buf bytes.Buffer = *bytes.NewBuffer(in.Buf)
	jsonBytes, err := s.Service.HandleBatchJSON(ctx, buf, in.Usr)
	if err != nil {
		return nil, err
	}
	return &pb.SliceByteErrResp{JsonBytes: jsonBytes, Err: "nil"}, nil
}

// FetchURLs получает URL-ы определенного пользователя
func (s *Server) FetchURLs(ctx context.Context, in *pb.StrReq) (*pb.SliceByteErrResp, error) {
	/*cookie, usr, err := s.Service.HandleCookie("")
	if err != nil {
		return nil, err
	}
	fmt.Println(cookie.String())
	fmt.Println(usr)*/
	jsonBytes, err := s.Service.FetchURLs(ctx, in.CookieValue)
	if err != nil {
		return nil, err
	}
	return &pb.SliceByteErrResp{JsonBytes: jsonBytes, Err: "nil"}, nil
}

func (s *Server) HandleCookie(ctx context.Context, in *pb.StrReq) (*pb.StrStrErrResp, error) {
	cookie, usr, err := s.Service.HandleCookie("")
	if err != nil {
		return nil, err
	}
	fmt.Println(cookie.String())
	fmt.Println(usr)
	return &pb.StrStrErrResp{CookieValue: cookie.String(), Usr: usr, Err: "nil"}, nil
}

// NewGRPCserver создает объект структуры, которая содержит реализацию серверной части
func NewGRPCserver(service *service.Service) {
	//var host string
	//var port string
	//hostAndPort := strings.Split(address, ":")
	//hostAndPort := strings.Split(service.Shu.Settings.GRPCaddress, ":")
	//host = hostAndPort[0]
	//port = hostAndPort[1]

	//addr := fmt.Sprintf("%s:%s", host, port)
	lis, err := net.Listen("tcp", service.Shu.Settings.GRPCaddress) // будем ждать запросы по этому адресу

	if err != nil {
		service.Logger.Info("Init gRPC service", zap.String("error", err.Error()))
		//log.Println("error starting gRPC server: ", err)
		return
		//os.Exit(1)
	}

	service.Logger.Info("Init gRPC service", zap.String("start at host:port", service.Shu.Settings.GRPCaddress))
	//log.Println("tcp listener started at host and port: ", service.Shu.Settings.GRPCaddress)
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
		service.Logger.Info("Init gRPC service", zap.String("error", err.Error()))
		return
		//log.Println("error serving grpc: ", err)
		//os.Exit(1)
	}
}
