package routes

import (
	"net/http"
	"server/global"
	"strconv"
	"database/sql"
	"server/message"
	"server/user"
)

func GetMessages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db := global.OpenDB()
	defer db.Close()

	// Get field from request
	userId, _ := strconv.Atoi(r.URL.Query().Get("userId"))

	rows, err := message.GetConcernedUsers(db, userId)
	if err != nil {
		db.Close()

		println(err.Error())
		global.SendError(w, "internal_error", http.StatusInternalServerError)
		return
	}

	users := make(map[int64]user.Info)

	for rows.Next() {
		userInfo := user.Info{}
		var address, birthDate sql.NullString
		rows.Scan(&userInfo.Id, &userInfo.Type, &userInfo.Mail, &userInfo.FirstName, &userInfo.LastName, &userInfo.City, &userInfo.PhoneNumber, &address, &birthDate)
		userInfo.Address = address.String
		userInfo.BirthDate = birthDate.String

		users[userInfo.Id] = userInfo
	}
	rows.Close()

	// Retrieve every messages
	rows, err = message.GetFromUserId(db, userId)
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