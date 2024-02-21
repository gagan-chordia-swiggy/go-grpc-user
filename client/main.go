package main

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"go-project/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
)

type User struct {
	Id       uint64 `json:"id"`
	Name     string `json:"name"`
	Age      uint32 `json:"age"`
	Username string `json:"username"`
	Password string `json:"password"`
}

var address = ":8544"
var baseUrl = "/api/v1/users"

func main() {
	connection, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("Failed to connect %v", err)
	}

	defer connection.Close()

	client := proto.NewUserServiceClient(connection)
	router := mux.NewRouter()

	router.HandleFunc(baseUrl, func(w http.ResponseWriter, r *http.Request) {
		CreateUser(client, w, r)
	})

	http.ListenAndServe(":8989", router)
}

func CreateUser(client proto.UserServiceClient, w http.ResponseWriter, r *http.Request) {
	var usr User
	json.NewDecoder(r.Body).Decode(&usr)

	user := &proto.User{
		Username: usr.Username,
		Password: usr.Password,
	}

	res, err := client.Create(context.Background(), &proto.CreateUserRequest{
		User: user,
	})

	if err != nil {
		json.NewEncoder(w).Encode(err)
		return
	}

	json.NewEncoder(w).Encode(res)
}
