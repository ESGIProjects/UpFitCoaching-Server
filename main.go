package main

import (
	"log"
	"net/http"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"server/routes"

	"firebase.google.com/go/messaging"
	"server/global"
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

	// Forum-handling routes
	router.HandleFunc("/forums/", routes.GetForums).Methods("GET")
	router.HandleFunc("/threads/", routes.GetThreads).Methods("GET")
	router.HandleFunc("/thread/", routes.GetThread).Methods("GET")
	router.HandleFunc("/thread/", routes.CreateThread).Methods("POST")
	router.HandleFunc("/post/", routes.AddPost).Methods("POST")

	// Debug routes
	router.HandleFunc("/notification", DebugNotifications).Methods("GET")

	// Listen on port 80
	err := http.ListenAndServe(":80", router)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func DebugNotifications(w http.ResponseWriter, r *http.Request) {
	db := global.OpenDB()
	defer db.Close()

	tokensDictionary := global.GetTokens(db, 15)
	notifications := make([]*messaging.Message, 0)

	for _, tokens := range tokensDictionary {
		for _, token := range tokens {

			notification := global.DebugNotification(token)
			notifications = append(notifications, notification)
		}
	}

	global.SendNotifications(notifications...)
}