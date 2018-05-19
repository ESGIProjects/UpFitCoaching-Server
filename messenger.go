package main

import (
	"net/http"
	"encoding/json"
	"log"
	"strconv"
	"github.com/gorilla/websocket"
	"errors"
	"database/sql"
	"server/user"
	"server/global"
)

func GetMessages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db := global.OpenDB()
	defer db.Close()

	// Retrieve field
	userId, err := strconv.Atoi(r.URL.Query().Get("userId"))
	if err != nil {
		db.Close()

		println(err.Error())
		global.SendError(w, "internal_error", http.StatusNotModified)
		return
	}

	// Retrieve every concerned users
	rows, err := db.Query("(SELECT id, type, mail, firstName, lastName, city, phoneNumber, address, birthDate FROM users NATURAL LEFT JOIN coaches NATURAL LEFT JOIN clients WHERE id = ?) UNION (SELECT id, type, mail, firstName, lastName, city, phoneNumber, address, birthDate FROM users NATURAL RIGHT JOIN (SELECT sender AS id FROM messages WHERE receiver = ? UNION SELECT receiver AS id FROM messages WHERE sender = ?) AS u NATURAL LEFT JOIN coaches NATURAL LEFT JOIN clients)", userId, userId, userId)
	if err != nil {
		db.Close()

		println(err.Error())
		global.SendError(w, "internal_error", http.StatusNotModified)
		return
	}

	usersMap := make(map[int64]user.Info)

	for rows.Next() {
		userInfo := user.Info{}
		var address, birthDate sql.NullString
		rows.Scan(&userInfo.Id, &userInfo.Type, &userInfo.Mail, &userInfo.FirstName, &userInfo.LastName, &userInfo.City, &userInfo.PhoneNumber, &address, &birthDate)
		userInfo.Address = address.String
		userInfo.BirthDate = birthDate.String

		usersMap[userInfo.Id] = userInfo
	}

	// Retrieve every messages
	rows, err = db.Query("SELECT * FROM messages WHERE sender = ? OR receiver = ? ORDER BY date DESC", userId, userId)
	if err != nil {
		db.Close()

		println(err.Error())
		global.SendError(w, "internal_error", http.StatusNotModified)
		return
	}

	defer rows.Close()

	messagesList := []Message{}

	for rows.Next() {
		message := Message{}
		var sender, receiver int64
		rows.Scan(&message.Id, &sender, &receiver, &message.Date, &message.Content)
		message.Sender = usersMap[sender]
		message.Receiver = usersMap[receiver]

		messagesList = append(messagesList, message)
	}

	// Format the response
	json, _ := json.Marshal(messagesList)

	// Send the response back
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

func SaveMessage(message Message) (int64, error) {
	db := global.OpenDB()

	// Insertion
	res, err := db.Exec("INSERT INTO messages (sender, receiver, date, content) VALUES(?, ?, ?, ?)", message.Sender.Id, message.Receiver.Id, message.Date, message.Content)
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