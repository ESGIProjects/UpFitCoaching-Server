package main

import (
	"log"
	"net/http"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"server/routes"

	"firebase.google.com/go/messaging"
	"server/global"
	"strconv"
)

func main() {
	// Creating router
	router := mux.NewRouter()

	// Notification token
	router.HandleFunc("/token/", AddOrUpdateToken).Methods("PUT")

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
	router.HandleFunc("/events/", routes.CancelEvent).Methods("DELETE")

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

func AddOrUpdateToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	db := global.OpenDB()
	defer db.Close()

	// Get fields from request
	userId, _ := strconv.Atoi(r.PostFormValue("userId"))
	token := r.PostFormValue("token")

	// Possible nil value
	oldToken := r.PostFormValue("oldToken")

	var query string
	args := make([]interface{}, 0)


	if oldToken == "" {
		// Add
		query = "INSERT INTO tokens (userId, token) VALUES(?, ?)"
		args = append(args, userId, token)

	} else {
		// Update
		query = "UPDATE tokens SET token = ? WHERE token = ?"
		args = append(args, token, oldToken)
	}

	// Inserting into DB
	_, err := db.Exec(query, args...)
	if err != nil {
		db.Close()

		print(err.Error())
		global.SendError(w, "internal_error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
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