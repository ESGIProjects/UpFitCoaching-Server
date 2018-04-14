package main

import (
    "log"
    "net/http"
    _ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/signup/", SignUp).Methods("POST")
	router.HandleFunc("/signin/", SignIn).Methods("POST")
	router.HandleFunc("/forgot/", Forgot).Methods("POST")

	router.HandleFunc("/messenger/", GetMessages).Methods("GET")
	router.HandleFunc("/messenger/", SendMessage).Methods("POST")

	router.HandleFunc("/ws", handleConnections)

	go handleMessages()

	err := http.ListenAndServe(":80", router)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}