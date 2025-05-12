// server пакет сервера gRPC
package server

import (
	"bytes"
	"context"
	"log"
	"net"
	"net/http"

	"github.com/developerc/reductorUrl/internal/general"
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
	buf := *bytes.NewBuffer(in.Buf)
	jsonBytes, err := s.Service.HandleBatchJSON(ctx, buf, in.Usr)
	if err != nil {
		return nil, err
	}
	return &pb.SliceByteErrResp{JsonBytes: jsonBytes, Err: "nil"}, nil
}

// FetchURLs получает URL-ы определенного пользователя
func (s *Server) FetchURLs(ctx context.Context, in *pb.StrReq) (*pb.SliceByteErrResp, error) {
	jsonBytes, err := s.Service.FetchURLs(ctx, in.CookieValue)
	if err != nil {
		return nil, err
	}
	return &pb.SliceByteErrResp{JsonBytes: jsonBytes, Err: "nil"}, nil
}

// HandleCookie обрабатывает куки
func (s *Server) HandleCookie(ctx context.Context, in *pb.StrReq) (*pb.StrStrErrResp, error) {
	cookie, usr, err := s.Service.HandleCookie(in.CookieValue)
	if err != nil {
		return nil, err
	}
	return &pb.StrStrErrResp{CookieValue: cookie.String(), Usr: usr, Err: "nil"}, nil
}

// GetStatsSvc обрабатывает статистику
func (s *Server) GetStatsSvc(ctx context.Context, in *pb.StrReq) (*pb.SliceByteErrResp, error) {
	ipNet := net.ParseIP(in.CookieValue) //здесь в CookieValue находится IP адрес
	jsonBytes, err := s.Service.GetStatsSvc(ctx, ipNet)
	if err != nil {
		return nil, err
	}
	return &pb.SliceByteErrResp{JsonBytes: jsonBytes, Err: "nil"}, nil
}

// DelURLs делает отметку об удалении коротких URL-ы определенного пользователя
func (s *Server) DelURLs(ctx context.Context, in *pb.StrByteReq) (*pb.ErrMess, error) {
	buf := bytes.NewBuffer(in.JsonBytes)
	err := s.Service.DelURLs(in.CookieValue, *buf)
	if err != nil {
		return nil, err
	}
	return &pb.ErrMess{Err: "nil"}, nil
}

// NewGRPCserver создает объект структуры, которая содержит реализацию серверной части
func NewGRPCserver(ctx context.Context, service *service.Service) {
	general.CntrAtomVar.IncrCntr()
	lis, err := net.Listen("tcp", service.Shu.Settings.GRPCaddress) // будем ждать запросы по этому адресу

	if err != nil {
		service.Logger.Info("Init gRPC service", zap.String("error", err.Error()))
		return
	}

	service.Logger.Info("Init gRPC service", zap.String("start at host:port", service.Shu.Settings.GRPCaddress))
	// создадим сервер grpc с перехватчиком
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(serverInterceptor))

	go func() {
		<-ctx.Done()
		service.Logger.Info("Server gRPC", zap.String("shutdown", "begin"))
		grpcServer.GracefulStop()
		service.Logger.Info("Server gRPC", zap.String("shutdown", "end"))
		general.CntrAtomVar.DecrCntr()
		general.CntrAtomVar.SentNotif()
	}()

	reductorServiceServer := NewServer()
	reductorServiceServer.Service = service

	pb.RegisterReductorServiceServer(grpcServer, reductorServiceServer)

	if err := grpcServer.Serve(lis); err != nil {
		service.Logger.Info("Init gRPC service", zap.String("error", err.Error()))
		return
	}
}

// serverInterceptor серверный перехватчик
func serverInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Printf("[INFO]  %v", info.FullMethod)
	return handler(ctx, req)
}
