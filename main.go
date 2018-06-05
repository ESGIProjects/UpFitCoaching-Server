package main

import (
	"log"
	"net/http"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"server/routes"

	"google.golang.org/api/option"
	"context"
	"firebase.google.com/go/messaging"
	"firebase.google.com/go"
	"fmt"
)

func main() {
	// Creating router
	router := mux.NewRouter()

	// User-handling routes
	router.HandleFunc("/checkmail/", routes.ExistingMail).Methods("POST")
	router.HandleFunc("/signup/", routes.SignUp).Methods("POST")
	router.HandleFunc("/signin/", routes.SignIn).Methods("POST")
	router.HandleFunc("/forgot/", routes.Forgot).Methods("POST")

	// Message-handling routes
	router.HandleFunc("/messages/", routes.GetMessages).Methods("GET")
	router.HandleFunc("/ws", handleConnections)
	go handleMessages()

	// Event-handling routes
	router.HandleFunc("/events/", routes.GetEvents).Methods("GET")
	router.HandleFunc("/events/", routes.AddEvent).Methods("POST")

	// Notification test routes
	router.HandleFunc("/notification", NotificationTest).Methods("GET")

	// Listen on port 80
	err := http.ListenAndServe(":80", router)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func NotificationTest(w http.ResponseWriter, r *http.Request) {
	// Initializing Firebase app
	opt := option.WithCredentialsFile("upfit-serviceAccountKey.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	// Obtain a messaging client from the Firebase app
	ctx := context.Background()
	client, err := app.Messaging(ctx)

	// This registration token comes from the client FCM SDKs
	registrationToken := "f4SDWAFsSGk:APA91bFT-yb9NzMbF4idW3cS7ckW8oAZCYQ77fJ9HVhCLQ5zUjfm_5jN7n_1EiEM50aDfOA1bsDGemv121OiDJQL8se5MrMkgQiIdKNdKMIjeKF7bfVchFVymEZR7Zf_my78K9CPRX7x"

	message := &messaging.Message{
		Notification: &messaging.Notification{
			Title: "Titre test",
			Body: "Body test",
		},
		Token: registrationToken,
	}

	response, err := client.Send(ctx, message)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Successfully sent message:", response)
}