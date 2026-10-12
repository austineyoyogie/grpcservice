package server

import (
	"gocgrpcarchitecture/generate/proto/userProto"
	"gocgrpcarchitecture/internal/app/resources/user/resource"
	"gocgrpcarchitecture/internal/app/resources/user/service/implementation"
	"gocgrpcarchitecture/internal/pkg/configs/database"
	"log"
	"net"

	"google.golang.org/grpc"
)

var (
	CD = database.CDriver()
)

type HttpServer struct {
	addr string
}

func NewHttpServer(addr string) *HttpServer {
	return &HttpServer{addr: addr}
}

func (s *HttpServer) Start() error {
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		log.Fatal("GRPC [SERVER] failed to listen.")
	}

	grpcServer := grpc.NewServer()
	userService := implementation.NewUserServiceImpl(CD)
	userProto.RegisterUserServiceServer(grpcServer, &resource.UserServerResource{UserService: userService})

	//documentService := implementation.NewDocumentServiceImpl(CD)
	//documentProto.RegisterDocumentServiceServer(grpcServer, &resource.DocumentServerResource{DocumentService: documentService})

	log.Println("GRPC [SERVER] listening on port", s.addr)
	if err := grpcServer.Serve(lis); err != nil {
		grpcServer.GracefulStop()
		log.Printf("GRPC [SERVER] failed to start: %v \n", err)
	}
	return nil
}
