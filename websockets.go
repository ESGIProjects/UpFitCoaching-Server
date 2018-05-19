package main

import (
	"net/http"
	"log"
	"strconv"
	"github.com/gorilla/websocket"
	"errors"
	"server/global"
	"server/message"
)

type Client struct {
	Id			int64
	Socket		* websocket.Conn
}

var clients = make(map[Client]bool) // connected clients
var broadcast = make(chan message.Info) // broadcast channel

var upgrader = websocket.Upgrader{
	CheckOrigin: func (r *http.Request) bool {
		return true
	}}

func saveMessage(messageInfo message.Info) (int64, error) {
	db := global.OpenDB()

	// Insertion
	res, err := message.Save(db, messageInfo)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

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

	log.Printf("New user Entered (id: %d)", id)

	defer ws.Close()

	// register new client
	client, err := getClient(id)
	if err != nil {
		client.Id = id
		client.Socket = ws
	}

	clients[client] = true

	for {
		var messageInfo message.Info

		// Read new message as JSON
		err := ws.ReadJSON(&messageInfo)
		log.Printf("Message: %s", messageInfo.Content)

		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, client)
			break
		}

		// Save the message in database
		id, err := saveMessage(messageInfo)
		if err != nil {
			log.Printf("error: %v", err)
			break
		}

		messageInfo.Id = id

		// Send the new message to broadcast channel
		broadcast <- messageInfo
	}
}

func handleMessages() {
	for {
		message := <-broadcast

		for client := range clients {

			if message.Receiver.Id == client.Id {
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