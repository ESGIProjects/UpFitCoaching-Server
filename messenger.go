package main

import (
	"net/http"
	"encoding/json"
	"log"
	"strconv"
	"github.com/gorilla/websocket"
	"errors"
)

func GetConversations(w http.ResponseWriter, r *http.Request) {

}

func GetConversation(w http.ResponseWriter, r *http.Request) {
	// Set the response header
	w.Header().Set("Content-Type", "application/json")

	// Connecting to the database
	db := dbConn()
	defer db.Close()

	// Retrieve fields
	coachId, err := strconv.Atoi(r.URL.Query().Get("coachId"))
	if err != nil {
		error := ErrorMessage{"internal_error"}
		json, _ := json.Marshal(error)

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(json)

		db.Close()
		return
	}

	clientId, err := strconv.Atoi(r.URL.Query().Get("clientId"))
	if err != nil {
		error := ErrorMessage{"internal_error"}
		json, _ := json.Marshal(error)

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(json)

		db.Close()
		return
	}

	var numberByPage = 20
	var page int
	pageString := r.URL.Query().Get("page")
	if pageString == "" {
		page = 0
	} else {
		page, err = strconv.Atoi(pageString)
		if err != nil {
			error := ErrorMessage{"internal_error"}
			json, _ := json.Marshal(error)

			w.WriteHeader(http.StatusInternalServerError)
			w.Write(json)

			db.Close()
			return
		}
	}

	rows, err := db.Query("SELECT * FROM messages WHERE (fromId = ? AND toId = ?) OR (fromId = ? AND toId = ?) ORDER BY date DESC LIMIT ?,?", coachId, clientId, clientId, coachId, page * numberByPage, numberByPage)
	if err != nil {
		error := ErrorMessage{"internal_error"}
		json, _ := json.Marshal(error)

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(json)

		db.Close()
		return
	}

	defer rows.Close()

	messagesList := []Message{}
	var id int64
	var fromId, fromType, toId, toType int
	var date, content string

	for rows.Next() {
		rows.Scan(&id, &fromId, &fromType, &toId, &toType, &date, &content)
		message := Message{}
		message.Id = id
		message.FromUserId = fromId
		message.FromUserType = fromType
		message.ToUserId = toId
		message.ToUserType = toType
		message.Date = date
		message.Content = content

		messagesList = append(messagesList, message)
	}

	// Format the response
	json, _ := json.Marshal(messagesList)

	// Send the response back
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

func SaveMessage(message Message) (int64, error) {
	// Connecting to the database
	db := dbConn()

	// Insertion
	res, err := db.Exec("INSERT INTO messages (fromId, fromType, toId, toType, date, content) VALUES(?, ?, ?, ?, ?, ?)", message.FromUserId, message.FromUserType, message.ToUserId, message.ToUserType, message.Date, message.Content)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

// Websocket

type Client struct {
	Id			int
	UserType	int
	Socket		* websocket.Conn
}

var clients = make(map[Client]bool) // connected clients
var broadcast = make(chan Message) // broadcast channel

var upgrader = websocket.Upgrader{
	CheckOrigin: func (r *http.Request) bool {
		return true
}}

func getClient(id int, userType int) (Client, error) {
	for client := range clients {
		if client.Id == id && client.UserType == userType {
			return client, nil
		}
	}

	return Client{}, errors.New("not connected")
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// upgrade to ws
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	userType, err := strconv.Atoi(r.URL.Query().Get("type"))

	if err != nil {
		ws.Close()
		return
	}

	log.Printf("New User Entered (id: %d, type: %d)", id, userType)

	defer ws.Close()

	// register new client
	client, err := getClient(id, userType)
	if err != nil {
		client.Id = id
		client.UserType = userType
		client.Socket = ws
	}

	clients[client] = true

	for {
		var message Message

		// Read new message as JSON
		err := ws.ReadJSON(&message)
		log.Printf("Message: %s", message.Content)

		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, client)
			break
		}

		// Save the message in database
		id, err := SaveMessage(message)
		if err != nil {
			log.Printf("error: %v", err)
			break
		}

		message.Id = id

		// Send the new message to broadcast channel
		broadcast <- message
	}
}

func handleMessages() {
	for {
		message := <-broadcast

		for client := range clients {

			if message.ToUserId == client.Id && message.ToUserType == client.UserType {
				err := client.Socket.WriteJSON(message)
				if err != nil {
					log.Printf("error %v", err)
					client.Socket.Close()
					delete(clients, client)
				}
			}
		}

		// PUSH NOTIFICATION HERE
	}
}