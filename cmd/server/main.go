package main

import (
	"Chotot-Microservice/cmd/pb"
	"Chotot-Microservice/config"
	"Chotot-Microservice/repository"
	"Chotot-Microservice/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"net"
)

const (
	addr     = "0.0.0.0:50051"
	certFile = "cmd/ssl/server.crt"
	keyFile  = "cmd/ssl/server.pem"
)

var (
	db = config.ConnectDatabase()

	userRepo    = repository.NewUserRepository(db)
	productRepo = repository.NewProductRepository(db)

	jwtService = service.NewJWTService()
)

type authServer struct {
	pb.AuthServiceServer
	userRepo   repository.UserRepository
	jwtService *service.JWTService
}

type userServer struct {
	pb.UserServiceServer
	userRepo   repository.UserRepository
	jwtService *service.JWTService
}

type productServer struct {
	pb.ProductServiceServer
	productRepo repository.ProductRepository
	jwtService  *service.JWTService
}

func runGRPCServer(enabledTLS bool, lis net.Listener) error {
	var opts []grpc.ServerOption
	if enabledTLS {
		creds, err := credentials.NewServerTLSFromFile(certFile, keyFile)
		if err != nil {
			return err
		}
		opts = append(opts, grpc.Creds(creds))
	}

	s := grpc.NewServer(opts...)
	pb.RegisterAuthServiceServer(s, &authServer{
		userRepo:   userRepo,
		jwtService: jwtService,
	})

	pb.RegisterUserServiceServer(s, &userServer{
		userRepo:   userRepo,
		jwtService: jwtService,
	})

	pb.RegisterProductServiceServer(s, &productServer{
		productRepo: productRepo,
		jwtService:  jwtService,
	})
	log.Printf("listening on %s\n", addr)
	return s.Serve(lis)
}

func main() {
	defer config.CloseDatabase(db)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	if err = runGRPCServer(true, lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
