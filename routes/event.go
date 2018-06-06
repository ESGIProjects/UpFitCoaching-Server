package routes

import (
	"net/http"
	"server/global"
	"strconv"
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

	// Get concerned users
	users, err := event.GetUsersList(db, userId)
	if err != nil {
		db.Close()

		println(err.Error())
		global.SendError(w, "internal_error", http.StatusInternalServerError)
		return
	}

	// Retrieve every events
	rows, err := event.GetFromUserId(db, userId)
	if err != nil {
		db.Close()

		print(err.Error())
		global.SendError(w, "internal_error", http.StatusInternalServerError)
		return
	}

	events := make([]event.Info, 0)

	for rows.Next() {
		eventInfo := event.Info{}
		var firstUser, secondUser, createdBy, updatedBy int64
		rows.Scan(&eventInfo.Id, &eventInfo.Name, &eventInfo.Type, &eventInfo.Status, &firstUser, &secondUser, &eventInfo.Start, &eventInfo.End, &eventInfo.Created, &createdBy, &eventInfo.Updated, &updatedBy)

		eventInfo.FirstUser = users[firstUser]
		eventInfo.SecondUser = users[secondUser]
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
	firstUser, _ := strconv.Atoi(r.PostFormValue("firstUser"))
	secondUser, _ := strconv.Atoi(r.PostFormValue("secondUser"))
	eventInfo.Start = r.PostFormValue("start")
	eventInfo.End = r.PostFormValue("end")
	eventInfo.Created = r.PostFormValue("created")
	createdBy, _ := strconv.Atoi(r.PostFormValue("createdBy"))

	status := 0

	// Inserting into DB
	res, err := db.Exec("INSERT INTO events (name, type, status, firstUser, secondUser, start, end, created, createdBy, updated, updatedBy) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", eventInfo.Name, eventType, status, firstUser, secondUser, eventInfo.Start, eventInfo.End, eventInfo.Created, createdBy, eventInfo.Created, createdBy)
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
	eventInfo.Type = eventType
	eventInfo.Status = status

	_, err = user.GetFromId(db, &eventInfo.FirstUser, int64(firstUser))
	if err != nil {
		db.Close()

		print(err.Error())
		global.SendError(w, "internal_error", http.StatusInternalServerError)
		return
	}

	_, err = user.GetFromId(db, &eventInfo.SecondUser, int64(secondUser))
	if err != nil {
		db.Close()

		print(err.Error())
		global.SendError(w, "internal_error", http.StatusInternalServerError)
		return
	}
	eventInfo.Updated = eventInfo.Created

	if createdBy == firstUser {
		eventInfo.CreatedBy = eventInfo.FirstUser
	} else {
		eventInfo.CreatedBy = eventInfo.SecondUser
	}
	eventInfo.UpdatedBy = eventInfo.CreatedBy

	global.SendJSON(w, eventInfo, http.StatusCreated)
}

func UpdateEvent(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db := global.OpenDB()
	defer db.Close()

	// Get event ID from request
	eventId, _ := strconv.Atoi(r.PostFormValue("eventId"))

	// Retrieve event from DB
	eventInfo := event.Info{}
	err := event.GetFromId(db, &eventInfo, int64(eventId))
	if err != nil {
		db.Close()

		print(err.Error())
		global.SendError(w, "internal_error", http.StatusInternalServerError)
		return
	}

	// Update fields from request
	eventInfo.Name = r.PostFormValue("name")
	eventInfo.Status, _ = strconv.Atoi(r.PostFormValue("status"))
	eventInfo.Start = r.PostFormValue("start")
	eventInfo.End = r.PostFormValue("end")
	eventInfo.Updated = r.PostFormValue("updated")
	updatedBy, _ := strconv.Atoi(r.PostFormValue("updatedBy"))

	if int64(updatedBy) == eventInfo.FirstUser.Id {
		eventInfo.UpdatedBy = eventInfo.FirstUser
	} else {
		eventInfo.UpdatedBy = eventInfo.SecondUser
	}

	// TODO - Change status if needed

	// Update the event
	event.Update(db, eventInfo)

	// TODO - Notification with status change

	global.SendJSON(w, eventInfo, http.StatusOK)
}