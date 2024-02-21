package main

import (
	"context"
	"errors"
	"go-project/database"
	"go-project/models"
	"go-project/proto"
	"google.golang.org/grpc"
	"gorm.io/gorm"
	"log"
	"net"
)

type Server struct {
	DB *gorm.DB
	proto.UserServiceServer
}

var address = ":8545"

func main() {
	lis, err := net.Listen("tcp", address)

	if err != nil {
		log.Fatalf("Failed to listed: %v\n", err)
	}

	log.Printf("Listening to: %v\n", address)

	server := grpc.NewServer()
	db := database.DatabaseConnection()
	proto.RegisterUserServiceServer(server, &Server{DB: db})

	err = server.Serve(lis)

	if err != nil {
		log.Fatalf("Failed to serve: %v\n", err)
	}
}

func (s *Server) Create(_ context.Context, req *proto.CreateUserRequest) (*proto.UserResponse, error) {
	usr := req.User
	user := models.User{
		Username: usr.Username,
		Password: usr.Password,
	}

	res := s.DB.Create(&user)

	if res.RowsAffected == 0 {
		return nil, errors.New("user creation unsuccessful")
	}

	response := &proto.UserResponse{
		User:    usr,
		Message: "User created successfully",
	}

	return response, nil
}

func (s *Server) AddName(_ context.Context, req *proto.CreateUserRequest) (*proto.UserResponse, error) {
	usr := req.User
	user := models.User{
		Id:       usr.Id,
		Username: usr.Username,
		Password: usr.Password,
		Name:     usr.Name,
		Age:      usr.Age,
	}

	res := s.DB.Save(&user)

	if res.RowsAffected == 0 {
		return nil, errors.New("user update unsuccessful")
	}

	response := &proto.UserResponse{
		User:    usr,
		Message: "User updated",
	}

	return response, nil
}
