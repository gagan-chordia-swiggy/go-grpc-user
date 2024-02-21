package main

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go-project/database"
	"go-project/mocks"
	"go-project/proto"
	"gorm.io/gorm"
	_ "gorm.io/gorm"
	"testing"
)

func setup() *gorm.DB {
	return database.DatabaseConnection()
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
	req := &proto.UserRequest{
		User: &proto.User{
			Username: "john",
			Password: "abc123",
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
	assert.Equal(t, expectedResponse.User, res.User)
	teardown()
}

func TestServer_AddNameAndAge(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	gormDb := setup()
	userService := mocks.NewMockUserServiceServer(controller)
	server := &Server{
		DB: gormDb,
	}
	req := &proto.UserRequest{
		User: &proto.User{
			Username: "john",
			Password: "abc123",
		},
	}
	expectedMessage := "User created successfully"
	expectedResponse := &proto.UserResponse{
		User:    req.User,
		Message: expectedMessage,
	}
	userService.EXPECT().Create(gomock.Any(), req).Return(expectedResponse, nil).AnyTimes()
	ctx := context.Background()
	server.Create(ctx, req)

	req = &proto.UserRequest{
		User: &proto.User{
			Id:   0,
			Name: "name",
			Age:  12,
		},
	}
	expectedMessage = "User updated"
	expectedResponse = &proto.UserResponse{
		User:    req.User,
		Message: expectedMessage,
	}
	userService.EXPECT().AddName(gomock.Any(), req).Return(expectedResponse, nil).AnyTimes()
	res, err := server.AddName(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, expectedResponse.User, res.User)
	teardown()
}

func TestServer_Get(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	gormDb := setup()
	userService := mocks.NewMockUserServiceServer(controller)
	server := &Server{
		DB: gormDb,
	}
	req := &proto.UserRequest{
		User: &proto.User{
			Username: "john",
			Password: "abc123",
		},
	}
	expectedMessage := "User created successfully"
	expectedResponse := &proto.UserResponse{
		User:    req.User,
		Message: expectedMessage,
	}
	userService.EXPECT().Create(gomock.Any(), req).Return(expectedResponse, nil).AnyTimes()
	ctx := context.Background()
	server.Create(ctx, req)

	req = &proto.UserRequest{
		User: &proto.User{
			Username: "john",
		},
	}
	expectedMessage = "User fetched"
	expectedResponse = &proto.UserResponse{
		User:    req.User,
		Message: expectedMessage,
	}
	userService.EXPECT().Get(gomock.Any(), req).Return(expectedResponse, nil).AnyTimes()
	res, err := server.Get(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, expectedResponse.User, res.User)
	teardown()
}

func TestServer_Update(t *testing.T) {
	controller := gomock.NewController(t)
	defer controller.Finish()
	gormDb := setup()
	userService := mocks.NewMockUserServiceServer(controller)
	server := &Server{
		DB: gormDb,
	}
	req := &proto.UserRequest{
		User: &proto.User{
			Username: "john",
			Password: "abc123",
		},
	}
	expectedMessage := "User created successfully"
	expectedResponse := &proto.UserResponse{
		User:    req.User,
		Message: expectedMessage,
	}
	userService.EXPECT().Create(gomock.Any(), req).Return(expectedResponse, nil).AnyTimes()
	ctx := context.Background()
	server.Create(ctx, req)

	req = &proto.UserRequest{
		User: &proto.User{
			Name: "john",
		},
	}
	expectedMessage = "User updated"
	expectedResponse = &proto.UserResponse{
		User:    req.User,
		Message: expectedMessage,
	}
	userService.EXPECT().Update(gomock.Any(), req).Return(expectedResponse, nil).AnyTimes()
	res, err := server.Get(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, expectedResponse.User, res.User)
	teardown()
}
