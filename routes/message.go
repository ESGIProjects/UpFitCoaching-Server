package routes

import (
	"net/http"
	"server/global"
	"strconv"
	"server/message"
)

func GetMessages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	db := global.OpenDB()
	defer db.Close()

	// Get field from request
	userId, err := strconv.Atoi(r.URL.Query().Get("userId"))

	if err != nil {
		db.Close()

		print(err.Error())
		global.SendError(w, "internal_error", http.StatusInternalServerError)
		return
	}

	// Get concerned users
	users, err := message.GetUsersList(db, userId)
	if err != nil {
		db.Close()

		print(err.Error())
		global.SendError(w, "internal_error", http.StatusInternalServerError)
		return
	}

	// Retrieve every messages
	rows, err := message.GetFromUserId(db, userId)
	if err != nil {
		db.Close()

		print(err.Error())
		global.SendError(w, "internal_error", http.StatusInternalServerError)
		return
	}

	messages := make([]message.Info, 0)

	for rows.Next() {
		message := message.Info{}
		var sender, receiver int64
		rows.Scan(&message.Id, &sender, &receiver, &message.Date, &message.Content)

		message.Sender = users[sender]
		message.Receiver = users[receiver]

		messages = append(messages, message)
	}
	rows.Close()

	global.SendJSON(w, messages, http.StatusOK)
}