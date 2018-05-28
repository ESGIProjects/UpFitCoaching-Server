package main

import (
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"server/routes"
)

func main() {
	// Creating router
	router := mux.NewRouter()

	// User-handling routes
	router.HandleFunc("/checkmail/", routes.ExistingMail).Methods("POST")
	router.HandleFunc("/signup/", routes.SignUp).Methods("POST")
	router.HandleFunc("/signin/", routes.SignIn).Methods("POST")
	router.HandleFunc("/forgot/", routes.Forgot).Methods("POST")
	router.HandleFunc("/users/", routes.UpdateProfile).Methods("PUT")

	// Message-handling routes
	router.HandleFunc("/messages/", routes.GetMessages).Methods("GET")
	router.HandleFunc("/ws", handleConnections)
	go handleMessages()

	// Event-handling routes
	router.HandleFunc("/events/", routes.GetEvents).Methods("GET")
	router.HandleFunc("/events/", routes.AddEvent).Methods("POST")
	router.HandleFunc("/events/", routes.UpdateEvent).Methods("PUT")

	// Listen on port 80
	err := http.ListenAndServe(":80", router)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
