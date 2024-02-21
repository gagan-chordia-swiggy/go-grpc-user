package main

import (
	"context"
	"errors"
	"fmt"
	"go-project/proto"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net"
)

type Server struct {
	DB *gorm.DB
	proto.UserServiceServer
}

type User struct {
	gorm.Model
	Id       uint64 `json:"id"`
	Name     string `json:"name"`
	Age      uint32 `json:"age"`
	Username string `json:"username" gorm:"unique"`
	Password string `json:"password"`
}

var address = ":8544"

func DatabaseConnection() *gorm.DB {
	host := "localhost"
	port := 5432
	user := "postgres"
	password := "root1234"
	dbName := "postgres"

	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName)

	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})

	if err != nil {
		log.Fatalln("Error connecting to the database")
	}

	log.Print("Connected to the database")

	err = db.AutoMigrate(&User{})

	if err != nil {
		log.Fatalln("Error migrating the database")
	}

	return db
}

func main() {
	lis, err := net.Listen("tcp", address)

	if err != nil {
		log.Fatalf("Failed to listed: %v\n", err)
	}

	log.Printf("Listening to: %v\n", address)

	server := grpc.NewServer()
	db := DatabaseConnection()
	proto.RegisterUserServiceServer(server, &Server{DB: db})

	err = server.Serve(lis)

	if err != nil {
		log.Fatalf("Failed to serve: %v\n", err)
	}
}

func (s *Server) Create(_ context.Context, req *proto.CreateUserRequest) (*proto.UserResponse, error) {
	usr := req.User
	user := User{
		Name:     "",
		Username: usr.Username,
		Password: usr.Password,
		Age:      0,
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
