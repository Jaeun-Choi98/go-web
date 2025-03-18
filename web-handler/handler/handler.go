package handler

import (
	"net/http"

	"github.com/gorilla/mux"
)

func NewHandler() http.Handler {
	mux := mux.NewRouter()

	// practice json, file handler
	mux.HandleFunc("/bar", BarHandler)
	mux.Handle("/foo", FooHandler{})
	mux.HandleFunc("/json", JsonHandler)
	mux.HandleFunc("/file", FileHandler)

	// practice rest
	mux.HandleFunc("/user", CreateUser).Methods("POST")
	mux.HandleFunc("/users/{id:[0-9]+}", GetUserById).Methods("GET")
	mux.HandleFunc("/users", UpdateUserByEmail).Methods("POST")
	mux.HandleFunc("/users", PutUserByEmail).Methods("PUT")
	mux.HandleFunc("/users", DeleteUser).Methods("DELETE")
	return mux
}
