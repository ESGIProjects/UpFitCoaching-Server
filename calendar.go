package main

import (
	"net/http"
	"strconv"
	"encoding/json"
	"database/sql"
	"io/ioutil"
)

func GetEvents(w http.ResponseWriter, r *http.Request) {
	// Set the response header
	w.Header().Set("Content-Type", "application/json")

	// Connecting to the database
	db := dbConn()
	defer db.Close()

	// Retrieve field
	userId, err := strconv.Atoi(r.URL.Query().Get("userId"))
	if err != nil {
		error := ErrorMessage{"internal_error1"}
		json, _ := json.Marshal(error)

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(json)

		db.Close()
		return
	}

	// Retrieve every concerned users
	SQLQuery := `
	SELECT * FROM users
	NATURAL LEFT JOIN coaches
	NATURAL LEFT JOIN clients
	WHERE id IN
	(SELECT client AS id FROM events WHERE coach = ?
	UNION SELECT coach AS id FROM events WHERE client = ?
	UNION SELECT ? AS id);
	`

	rows, err := db.Query(SQLQuery, userId, userId, userId)
	if err != nil {
		print(err.Error())

		error := ErrorMessage{"internal_error2"}
		json, _ := json.Marshal(error)

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(json)

		db.Close()
		return
	}

	usersMap := make(map[int64]UserInfo)

	for rows.Next() {
		user := UserInfo{}
		var nullableAddress, nullableBirthDate sql.NullString
		rows.Scan(&user.Id, &user.Type, &user.Mail, &user.FirstName, &user.LastName, &user.City, &user.PhoneNumber, &nullableAddress, &nullableBirthDate)
		user.Address = nullableAddress.String
		user.BirthDate = nullableBirthDate.String

		usersMap[user.Id] = user
	}

	// Retrieve every events
	rows, err = db.Query("SELECT * FROM events WHERE client = ? OR coach = ? ORDER BY start", userId, userId)
	if err != nil {
		error := ErrorMessage{"internal_error3"}
		json, _ := json.Marshal(error)

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(json)

		db.Close()
		return
	}

	defer rows.Close()

	eventsList := []Event{}

	for rows.Next() {
		event := Event{}
		var client, coach, createdBy, updatedBy int64
		rows.Scan(&event.Id, &event.Name, &client, &coach, &event.Start, &event.End, &event.Created, &createdBy, &event.Updated, &updatedBy)
		event.Client = usersMap[client]
		event.Coach = usersMap[coach]
		event.CreatedBy = usersMap[createdBy]
		event.UpdatedBy = usersMap[updatedBy]

		eventsList = append(eventsList, event)
	}

	// Format the response
	json, _ := json.Marshal(eventsList)

	// Send the response back
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

func AddEvent(w http.ResponseWriter, r *http.Request) {
	// Set the response header
	w.Header().Set("Content-Type", "application/json")

	// Connecting to the database
	db := dbConn()

	// Retrieve the body of the request
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		println(err.Error())
		errorMessage := ErrorMessage{"internal_error"}
		json, _ := json.Marshal(errorMessage)

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(json)

		db.Close()
		return
	}

	// Creating Event struct from JSON
	var event Event
	err = json.Unmarshal(body, &event)
	if err != nil {
		println(err.Error())
		errorMessage := ErrorMessage{"internal_error"}
		json, _ := json.Marshal(errorMessage)

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(json)

		db.Close()
		return
	}

	// Inserting into DB
	res, err := db.Exec("INSERT INTO events (name, client, coach, start, end, created, createdBy, updated, updatedBy) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)", event.Name, event.Client.Id, event.Coach.Id, event.Start, event.End, event.Created, event.CreatedBy.Id, event.Updated, event.UpdatedBy.Id)
	if err != nil {
		println(err.Error())
		errorMessage := ErrorMessage{"internal_error"}
		json, _ := json.Marshal(errorMessage)

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(json)

		db.Close()
		return
	}

	// Get the new event ID
	id, err := res.LastInsertId()
	if err != nil {
		println(err.Error())
		errorMessage := ErrorMessage{"internal_error"}
		json, _ := json.Marshal(errorMessage)

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(json)

		db.Close()
		return
	}

	// Format the response
	event.Id = id
	json, _ := json.Marshal(event)

	// Send the response back
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}
