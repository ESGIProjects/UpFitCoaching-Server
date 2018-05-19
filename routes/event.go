package routes

import (
	"net/http"
	"server/global"
	"strconv"
	"database/sql"
	"server/user"
	"server/event"
)

func GetEvents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db := global.OpenDB()
	defer db.Close()

	// Get fields from request
	userId, err := strconv.Atoi(r.URL.Query().Get("userId"))
	if err != nil {
		db.Close()

		print(err.Error())
		global.SendError(w, "internal_error", http.StatusInternalServerError)
		return
	}

	rows, err := event.GetConcernedUsers(db, userId)
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

	// Retrieve every events
	rows, err = event.GetFromUserId(db, userId)
	if err != nil {
		db.Close()

		print(err.Error())
		global.SendError(w, "internal_error", http.StatusInternalServerError)
		return
	}

	var events []event.Info

	for rows.Next() {
		event := event.Info{}
		var client, coach, createdBy, updatedBy int64
		rows.Scan(&event.Id, &event.Name, &client, &coach, &event.Start, &event.End, &event.Created, &createdBy, &event.Updated, &updatedBy)
		event.Client = users[client]
		event.Coach = users[coach]
		event.CreatedBy = users[createdBy]
		event.UpdatedBy = users[updatedBy]

		events = append(events, event)
	}
	rows.Close()

	global.SendJSON(w, events, http.StatusOK)
}

func AddEvent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db := global.OpenDB()
	defer db.Close()

	// Get fields from request
	name := r.PostFormValue("name")
	client, _ := strconv.Atoi(r.PostFormValue("client"))
	coach, _ := strconv.Atoi(r.PostFormValue("coach"))
	start := r.PostFormValue("start")
	end := r.PostFormValue("end")
	created := r.PostFormValue("created")
	createdBy, _ := strconv.Atoi(r.PostFormValue("createdBy"))

	// Inserting into DB
	res, err := db.Exec("INSERT INTO events (name, client, coach, start, end, created, createdBy, updated, updatedBy) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)", name, client, coach, start, end, created, createdBy, created, createdBy)
	if err != nil {
		db.Close()

		print(err.Error())
		global.SendError(w, "internal_error", http.StatusInternalServerError)
		return
	}

	// Get the new event ID
	id, err := res.LastInsertId()
	if err != nil {
		db.Close()

		print(err.Error())
		global.SendError(w, "internal_error", http.StatusInternalServerError)
		return
	}

	eventInfo := event.Info{}
	eventInfo.Id = id
	_, err = user.GetFromId(db, &eventInfo.Client, int64(client))
	if err != nil {
		db.Close()

		print(err.Error())
		global.SendError(w, "internal_error", http.StatusInternalServerError)
		return
	}

	_, err = user.GetFromId(db, &eventInfo.Coach, int64(coach))
	if err != nil {
		db.Close()

		print(err.Error())
		global.SendError(w, "internal_error", http.StatusInternalServerError)
		return
	}

	eventInfo.Start = start
	eventInfo.Created = created
	eventInfo.Updated = created

	if createdBy == client {
		eventInfo.CreatedBy = eventInfo.Client
	} else {
		eventInfo.CreatedBy = eventInfo.Coach
	}
	eventInfo.UpdatedBy = eventInfo.CreatedBy

	global.SendJSON(w, eventInfo, http.StatusOK)
}
