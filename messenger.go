package main

import (
	"net/http"
	"encoding/json"
	"log"
	"strconv"
	"github.com/gorilla/websocket"
	"errors"
)

func GetMessages(w http.ResponseWriter, r *http.Request) {
	// Set the response header
	w.Header().Set("Content-Type", "application/json")

	// Connecting to the database
	db := dbConn()
	defer db.Close()

	// Retrieve field
	userId, err := strconv.Atoi(r.URL.Query().Get("userId"))
	if err != nil {
		error := ErrorMessage{"internal_error"}
		json, _ := json.Marshal(error)

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(json)

		db.Close()
		return
	}

	rows, err := db.Query("SELECT * FROM messages WHERE sender = ? OR receiver = ? ORDER BY date DESC", userId, userId)
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
	var id, sender, receiver int64
	var date, content string

	for rows.Next() {
		rows.Scan(&id, &sender, &receiver, &date, &content)
		message := Message{}
		message.Id = id
		message.Sender = sender
		message.Receiver = receiver
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
	res, err := db.Exec("INSERT INTO messages (sender, receiver, date, content) VALUES(?, ?, ?, ?)", message.Sender, message.Receiver, message.Date, message.Content)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

// Websocket

type Client struct {
	Id			int64
	Socket		* websocket.Conn
}

var clients = make(map[Client]bool) // connected clients
var broadcast = make(chan Message) // broadcast channel

var upgrader = websocket.Upgrader{
	CheckOrigin: func (r *http.Request) bool {
		return true
}}

func getClient(id int64) (Client, error) {
	for client := range clients {
		if client.Id == id {
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

	clientId, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		ws.Close()
		return
	}
	id := int64(clientId)

	log.Printf("New User Entered (id: %d)", id)

	defer ws.Close()

	// register new client
	client, err := getClient(id)
	if err != nil {
		client.Id = id
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

			if message.Receiver == client.Id {
				log.Printf("found correct client")
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