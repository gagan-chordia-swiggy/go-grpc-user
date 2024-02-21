package main

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"go-project/database"
	"go-project/models"
	"go-project/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/gorm"
	"log"
	"net/http"
	"reflect"
)

var address = ":8545"
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
	}).Methods("POST")
	router.HandleFunc(baseUrl+"/add", func(w http.ResponseWriter, r *http.Request) {
		AddNameAndAge(client, w, r)
	}).Methods("POST")

	http.ListenAndServe(":8989", router)
}

func CreateUser(client proto.UserServiceClient, w http.ResponseWriter, r *http.Request) {
	var usr models.User
	json.NewDecoder(r.Body).Decode(&usr)

	user := &proto.User{
		Username: usr.Username,
		Password: usr.Password,
	}

	res, err := client.Create(context.Background(), &proto.CreateUserRequest{
		User: user,
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

func AddNameAndAge(client proto.UserServiceClient, w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	if username == "" || password == "" || !ok {
		w.WriteHeader(http.StatusUnprocessableEntity)
	}

	db := database.DatabaseConnection()

	user, err := getUserByUsername(db, username)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("Wrong credentials")
		return
	}

	if !reflect.DeepEqual(user.Password, password) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode("Wrong credentials")
		return
	}
	var usr models.User
	json.NewDecoder(r.Body).Decode(&usr)

	u := &proto.User{
		Id:       user.Id,
		Name:     usr.Name,
		Username: user.Username,
		Password: user.Password,
		Age:      usr.Age,
	}

	res, err := client.AddName(context.Background(), &proto.CreateUserRequest{
		User: u,
	})

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

func getUserByUsername(db *gorm.DB, username string) (*models.User, error) {
	var user models.User
	err := db.Where("username = ?", username).Find(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}
