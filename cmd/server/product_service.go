package main

import (
	"Chotot-Microservice/cmd/pb"
	"Chotot-Microservice/entity"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"strconv"
	"time"
)

func entityToProductProto(product *entity.Product) *pb.Product {
	return &pb.Product{
		Name:        product.ProductName,
		UserId:      int64(product.UserId),
		CatId:       product.CatId,
		TypeId:      product.TypeId,
		Price:       product.Price,
		State:       product.State,
		CreatedTime: product.CreatedTime.Unix(),
		ExpiredTime: product.ExpiredTime.Unix(),
		Address:     product.Address,
		Content:     product.Content,
	}
}

func productProtoToEntity(product *pb.Product) *entity.Product {
	return &entity.Product{
		ProductName: product.Name,
		UserId:      int(product.UserId),
		CatId:       product.CatId,
		TypeId:      product.TypeId,
		Price:       product.Price,
		State:       product.State,
		CreatedTime: time.Unix(product.CreatedTime, 0),
		ExpiredTime: time.Unix(product.ExpiredTime, 0),
		Address:     product.Address,
		Content:     product.Content,
	}
}

//func (s *productServer) GetProduct(ctx context.Context, in *pb.ProductId) (*pb.Product, error)

func (s *productServer) GetProducts(in *emptypb.Empty, stream pb.ProductService_GetProductsServer) error {
	log.Println("GetProducts")
	cur, err := s.productRepo.GetAllProducts()
	if err != nil {
		return status.Errorf(codes.Internal, err.Error())
	}

	for i := range cur {
		product := entityToProductProto(&cur[i])
		err = stream.Send(product)
		if err != nil {
			return status.Errorf(codes.Internal, err.Error())
		}
		time.Sleep(200 * time.Millisecond)
	}
	return nil
}

func (s *productServer) GetProductsByUserId(in *pb.UserToken, stream pb.ProductService_GetProductsByUserIdServer) error {
	log.Printf("GetProductsByUserId: %v\n", in.Header)
	userId, err := s.jwtService.ParseToken(in.Header)
	if err != nil {
		return status.Errorf(codes.Internal, err.Error())
	}
	cur, err := s.productRepo.UserProducts(userId)
	for i := range cur {
		product := entityToProductProto(&cur[i])
		err = stream.Send(product)
		if err != nil {
			return status.Errorf(codes.Internal, err.Error())
		}
		time.Sleep(200 * time.Millisecond)
	}
	return nil
}

func (s *productServer) CreateProduct(ctx context.Context, in *pb.Product) (*pb.Product, error) {
	log.Printf("CreateProduct: %v\n", in)
	userId, err := s.jwtService.ParseToken(in.Key.Header)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	in.UserId = int64(userId)
	product := productProtoToEntity(in)
	res, err := s.productRepo.InsertProduct(product)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return entityToProductProto(res), nil
}

func (s *productServer) UpdateProduct(ctx context.Context, in *pb.Product) (*emptypb.Empty, error) {
	log.Printf("UpdateProduct: %v\n", in)
	userId, err := s.jwtService.ParseToken(in.Key.Header)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	in.UserId = int64(userId)
	product := productProtoToEntity(in)
	product.Id, _ = strconv.Atoi(in.Id)
	_, err = s.productRepo.UpdateProduct(product)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (s *productServer) DeleteProduct(ctx context.Context, in *pb.ProductId) (*emptypb.Empty, error) {
	log.Printf("DeleteProduct: %v\n", in)
	userId, err := s.jwtService.ParseToken(in.Key.Header)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	err = s.productRepo.DeleteProduct(in.Id, strconv.FormatInt(int64(userId), 10))
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (s *productServer) SearchProducts(in *pb.SearchRequest, stream pb.ProductService_SearchProductsServer) error {
	log.Printf("SearchProduct: %v\n", in)
	cur, err := s.productRepo.SearchProductsByName(in.KeyName, in.KeyPage)
	if err != nil {
		return status.Errorf(codes.Internal, err.Error())
	}

	for i := range cur {
		product := entityToProductProto(&cur[i])
		err = stream.Send(product)
		if err != nil {
			return status.Errorf(codes.Internal, err.Error())
		}
		time.Sleep(200 * time.Millisecond)
	}
	return nil
}
