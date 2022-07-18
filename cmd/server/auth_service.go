package main

import (
	"Chotot-Microservice/cmd/pb"
	"Chotot-Microservice/entity"
	"context"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"strconv"
)

func (s *authServer) Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginResponse, error) {
	log.Printf("Login request: %v\n", in)

	res, err := s.userRepo.VerifyCredential(in.Phone)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal error: %v", err)
	}

	if res == nil {
		return nil, status.Errorf(codes.NotFound, "User not found")
	}

	if err = bcrypt.CompareHashAndPassword([]byte(res.Passwd), []byte(in.Password)); err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Invalid password")
	}

	token := s.jwtService.GenerateToken(strconv.FormatInt(int64(res.Id), 10))
	return &pb.LoginResponse{Token: token}, nil
}

func (s *authServer) Register(ctx context.Context, in *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	log.Printf("Register request: %v\n", in)

	res, _ := s.userRepo.IsDuplicatePhone(in.Phone)
	if res {
		return nil, status.Errorf(codes.AlreadyExists, "Phone already exists")
	}

	userReq := &entity.User{
		Phone:    in.Phone,
		Username: in.Username,
		Passwd:   in.Password,
		Address:  in.Address,
		Email:    in.Email,
		IsAdmin:  in.IsAdmin,
	}

	newUser, _ := s.userRepo.InsertUser(userReq)
	if newUser == nil {
		return nil, status.Errorf(codes.Internal, "Internal error")
	}
	token := s.jwtService.GenerateToken(strconv.FormatInt(int64(newUser.Id), 10))
	return &pb.RegisterResponse{
		Phone:    newUser.Phone,
		Username: newUser.Username,
		Password: newUser.Passwd,
		Token:    token,
	}, nil
}
