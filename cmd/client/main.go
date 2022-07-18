package main

import (
	"Chotot-Microservice/cmd/pb"
	"Chotot-Microservice/middleware"
	"Chotot-Microservice/service"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
)

const (
	addr     = "localhost:50051"
	certFile = "cmd/ssl/ca.crt"
)

func loadTLSCredentials() (credentials.TransportCredentials, error) {
	creds, err := credentials.NewClientTLSFromFile(certFile, "")
	if err != nil {
		return nil, err
	}
	return creds, nil
}

func main() {
	var opts []grpc.DialOption
	creds, err := loadTLSCredentials()

	if err != nil {
		log.Fatalf("Failed to load credentials: %v", err)
	}

	opts = append(opts, grpc.WithTransportCredentials(creds))

	conn, err := grpc.Dial(addr, opts...)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	authClient := pb.NewAuthServiceClient(conn)
	userClient := pb.NewUserServiceClient(conn)
	productClient := pb.NewProductServiceClient(conn)

	r := gin.Default()
	authRoutes := r.Group("/cho-tot/auth")
	{
		authRoutes.POST("/login", login(authClient))
		authRoutes.POST("/register", register(authClient))
	}

	userRoutes := r.Group("/cho-tot/user", middleware.AuthorizeJWT(service.NewJWTService()))
	{
		userRoutes.GET("/profile", getUser(userClient))
		userRoutes.PUT("/profile/update", updateUser(userClient))
		userRoutes.GET("/products", getUserProducts(productClient))
		userRoutes.POST("/products/create", createUserProduct(productClient))
		userRoutes.PUT("/products/update/:id", updateUserProduct(productClient))
		userRoutes.DELETE("/products/delete/:id", deleteUserProduct(productClient))
	}

	productRoutes := r.Group("/cho-tot/product")
	{
		productRoutes.GET("/", getProducts(productClient))
		productRoutes.GET("/search", searchProducts(productClient))
	}

	if err = r.Run(":8080"); err != nil {
		log.Fatalf("failed to run server: %v", err)
	}
}
