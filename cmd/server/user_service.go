package main

import (
	"Chotot-Microservice/cmd/pb"
	"Chotot-Microservice/entity"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
)

func entityToUserProto(user *entity.User) *pb.User {
	return &pb.User{
		Phone:    user.Phone,
		Username: user.Username,
		Password: user.Passwd,
		Address:  user.Address,
		Email:    user.Email,
		IsAdmin:  user.IsAdmin,
	}
}

func (s *userServer) GetUser(ctx context.Context, in *pb.UserToken) (*pb.User, error) {
	log.Printf("GetUser: %v", in)
	userId, err := s.jwtService.ParseToken(in.Header)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
	}
	user, err := s.userRepo.UserProfile(userId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal error: %v", err)
	}
	return entityToUserProto(user), nil
}

func (s *userServer) UpdateUser(ctx context.Context, in *pb.UpdateRequest) (*emptypb.Empty, error) {
	log.Printf("UpdateUser: %v", in)
	userId, err := s.jwtService.ParseToken(in.Key.Header)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
	}
	updatedUser := &entity.User{
		Username: in.Username,
		Address:  in.Address,
		Email:    in.Email,
	}
	_, err = s.userRepo.UpdateUser(userId, updatedUser)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal error: %v", err)
	}
	return &emptypb.Empty{}, nil
}
