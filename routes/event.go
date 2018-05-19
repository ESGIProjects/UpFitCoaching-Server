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

	events := make([]event.Info, 0)

	for rows.Next() {
		eventInfo := event.Info{}
		var client, coach, createdBy, updatedBy int64
		rows.Scan(&eventInfo.Id, &eventInfo.Name, &eventInfo.Type, &client, &coach, &eventInfo.Start, &eventInfo.End, &eventInfo.Created, &createdBy, &eventInfo.Updated, &updatedBy)
		eventInfo.Client = users[client]
		eventInfo.Coach = users[coach]
		eventInfo.CreatedBy = users[createdBy]
		eventInfo.UpdatedBy = users[updatedBy]

		events = append(events, eventInfo)
	}
	rows.Close()

	global.SendJSON(w, events, http.StatusOK)
}

func AddEvent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db := global.OpenDB()
	defer db.Close()

	eventInfo := event.Info{}

	// Get fields from request
	eventInfo.Name = r.PostFormValue("name")
	eventType, _ := strconv.Atoi(r.PostFormValue("type"))
	client, _ := strconv.Atoi(r.PostFormValue("client"))
	coach, _ := strconv.Atoi(r.PostFormValue("coach"))
	eventInfo.Start = r.PostFormValue("start")
	eventInfo.End = r.PostFormValue("end")
	eventInfo.Created = r.PostFormValue("created")
	createdBy, _ := strconv.Atoi(r.PostFormValue("createdBy"))

	// Inserting into DB
	res, err := db.Exec("INSERT INTO events (name, type, client, coach, start, end, created, createdBy, updated, updatedBy) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", eventInfo.Name, eventType, client, coach, eventInfo.Start, eventInfo.End, eventInfo.Created, createdBy, eventInfo.Created, createdBy)
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

	eventInfo.Id = id
	eventInfo.Type = &eventType

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
	eventInfo.Updated = eventInfo.Created

	if createdBy == client {
		eventInfo.CreatedBy = eventInfo.Client
	} else {
		eventInfo.CreatedBy = eventInfo.Coach
	}
	eventInfo.UpdatedBy = eventInfo.CreatedBy

	global.SendJSON(w, eventInfo, http.StatusCreated)
}
