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
	w.Header().Set("Content-Type", "application/json")
	db := dbConn()
	var id,idSender,idReceiver,body,timestamp string
	test,_ := r.URL.Query()["id"]
	rows,_ := db.Query("SELECT id,idSender,idReceiver,body,timestamp FROM Message WHERE (idSender='" + test[0] + "' OR idReceiver ='" +
		test[0] + "') AND deletedSender='0'")
	defer rows.Close()
	Messages := []Message{}
	for rows.Next(){
		Message := Message{}
		rows.Scan(&id,&idSender,&idReceiver,&body,&timestamp)
		/*Message.Id = id
		Message.IdSender = idSender
		Message.IdReceiver = idReceiver
		Message.Body = body
		Message.Timestamp = timestamp*/
		Messages = append(Messages, Message)
	}
	json,_ := json.Marshal(Messages)
	w.Write(json)
	w.WriteHeader(http.StatusOK)
}

func SendMessage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db := dbConn()
	_, err := db.Query("INSERT INTO Message (idSender,idReceiver,body,timestamp,deletedSender,deletedReceiver) VALUES ('" + r.PostFormValue("idSender") +
		"','" +r.PostFormValue("idReceiver") + "','" + r.PostFormValue("body") + "','" +r.PostFormValue("timestamp") + "','0','0')")
	if err == nil{
		w.WriteHeader(http.StatusCreated)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

// Websocket

type Message struct {
	Email string `json:"email"`
	Username string `json:"username"`
	Message string `json:"message"`
}

type Client struct {
	Id int
	UserType int
	Socket *websocket.Conn
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

	log.Printf("NEW CONNECTION OPENED !!")

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	userType, err := strconv.Atoi(r.URL.Query().Get("type"))

	if err != nil {
		ws.Close()
		return
	}

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
		var msg Message

		// Read new message as JSON
		err := ws.ReadJSON(&msg)

		log.Printf("Message after: %v", msg)

		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, client)
			break
		}

		// Send the new message to broadcast channel
		broadcast <- msg
	}
}

func handleMessages() {
	for {
		msg := <-broadcast

		for client := range clients {

			err := client.Socket.WriteJSON(msg)
			if err != nil {
				log.Printf("error %v", err)
				client.Socket.Close()
				delete(clients, client)
			}
		}
	}
}