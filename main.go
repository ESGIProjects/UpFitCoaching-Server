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
	"server/auth"
)

func main() {
	// Creating router
	router := mux.NewRouter()

	loginRouter := router.PathPrefix("/login").Subrouter()
	protectedRouter := router.PathPrefix("/").Subrouter()

	protectedRouter.Use(auth.VerifyToken)

	// Notification token
	protectedRouter.HandleFunc("/token/", AddOrUpdateToken).Methods("PUT")

	// User-handling routes
	loginRouter.HandleFunc("/checkmail/", routes.ExistingMail).Methods("POST")
	loginRouter.HandleFunc("/signup/", routes.SignUp).Methods("POST")
	loginRouter.HandleFunc("/signin/", routes.SignIn).Methods("POST")
	loginRouter.HandleFunc("/forgot/", routes.Forgot).Methods("POST")
	protectedRouter.HandleFunc("/users/", routes.UpdateProfile).Methods("PUT")

	// Message-handling routes
	protectedRouter.HandleFunc("/messages/", routes.GetMessages).Methods("GET")
	protectedRouter.HandleFunc("/ws", handleConnections)
	go handleMessages()

	// Event-handling routes
	protectedRouter.HandleFunc("/events/", routes.GetEvents).Methods("GET")
	protectedRouter.HandleFunc("/events/", routes.AddEvent).Methods("POST")
	protectedRouter.HandleFunc("/events/", routes.UpdateEvent).Methods("PUT")
	protectedRouter.HandleFunc("/events/", routes.CancelEvent).Methods("DELETE")

	// Forum-handling routes
	protectedRouter.HandleFunc("/forums/", routes.GetForums).Methods("GET")
	protectedRouter.HandleFunc("/threads/", routes.GetThreads).Methods("GET")
	protectedRouter.HandleFunc("/thread/", routes.GetThread).Methods("GET")
	protectedRouter.HandleFunc("/thread/", routes.CreateThread).Methods("POST")
	protectedRouter.HandleFunc("/post/", routes.AddPost).Methods("POST")

	// Follow Up-handling routes
	protectedRouter.HandleFunc("/appraisals/", routes.GetLastAppraisal).Methods("GET")
	protectedRouter.HandleFunc("/appraisals/", routes.CreateAppraisal).Methods("POST")
	protectedRouter.HandleFunc("/measurements/", routes.GetMeasurements).Methods("GET")
	protectedRouter.HandleFunc("/measurements/", routes.CreateMeasurements).Methods("POST")
	protectedRouter.HandleFunc("/tests/", routes.GetTests).Methods("GET")
	protectedRouter.HandleFunc("/tests/", routes.CreateTest).Methods("POST")

	// Prescriptions-handling routes
	protectedRouter.HandleFunc("/prescriptions/", routes.GetPrescriptions).Methods("GET")
	protectedRouter.HandleFunc("/prescriptions/", routes.CreatePrescription).Methods("POST")

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