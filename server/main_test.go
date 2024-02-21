package main

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go-project/mocks"
	"go-project/proto"
	"gorm.io/gorm"
	_ "gorm.io/gorm"
	"testing"
)

func setup() *gorm.DB {
	return DatabaseConnection()
}

func teardown() {
	setup().Exec("TRUNCATE TABLE USERS")
}

func TestServer_CreateUser(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	gormDb := setup()
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
	expectedMessage := "User created successfully"
	expectedResponse := &proto.UserResponse{
		User:    req.User,
		Message: expectedMessage,
	}

	userService.EXPECT().Create(gomock.Any(), req).Return(expectedResponse, nil).AnyTimes()
	ctx := context.Background()
	res, err := server.Create(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, res)
	teardown()
}
