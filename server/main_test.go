package main

import (
	"context"
	"encoding/base64"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"go-project/mocks"
	"go-project/proto"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	_ "gorm.io/gorm"
	"reflect"
	"testing"
)

func TestServer_CreateUser(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error occurred while creating mock db: %v", err)
	}
	defer db.Close()
	gormDb, err := gorm.Open(postgres.New(postgres.Config{
		Conn:       db,
		DriverName: "postgres",
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("An error occured while creating mock gorm instance %v", err)
	}
	userService := mocks.NewMockUserServiceServer(controller)
	server := &Server{
		DB: gormDb,
	}
	req := &proto.CreateUserRequest{
		User: &proto.User{
			Name: "John",
			Age:  12,
		},
	}
	expectedToken := base64.StdEncoding.EncodeToString([]byte("John:12"))
	expectedQuery := "^INSERT INTO users (.+)$"
	expectedRowsAffected := int64(1)
	expectedMessage := "User created successfully"

	mock.ExpectExec(expectedQuery).WillReturnResult(sqlmock.NewResult(1, expectedRowsAffected))
	userService.EXPECT().Create(gomock.Any(), req).Return(&proto.UserResponse{
		User:    req.User,
		Message: expectedMessage,
	}, nil)
	ctx := context.Background()
	res, err := server.Create(ctx, req)

	if err != nil {
		t.Errorf("Error occured while executing CreateUser method: %v", err)
	}
	if !reflect.DeepEqual(res.User.Name, req.User.Name) {
		t.Errorf("Expected %s, Got %s", req.User.Name, res.User.Name)
	}
	if !reflect.DeepEqual(res.User.Token, expectedToken) {
		t.Errorf("Expected %s, Got %s", expectedToken, res.User.Token)
	}
	if !reflect.DeepEqual(expectedMessage, res.Message) {
		t.Errorf("Expected %s, Got %s", expectedMessage, res.Message)
	}
}
