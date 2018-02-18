package main

import (
    "log"
    "net/http"
    "github.com/gorilla/mux"
    _ "github.com/go-sql-driver/mysql"
)

func main() {
    router := mux.NewRouter()

	router.HandleFunc("/signup/", SignUp).Methods("POST")
	router.HandleFunc("/signin/", SignIn).Methods("POST")
	router.HandleFunc("/forgot/", Forgot).Methods("POST")

	router.HandleFunc("/messenger/", GetMessages).Methods("GET")
	router.HandleFunc("/messenger/", SendMessage).Methods("POST")
	router.HandleFunc("/messenger/", DeleteMessages).Methods("DELETE")

    log.Fatal(http.ListenAndServe(":8000", router))
}