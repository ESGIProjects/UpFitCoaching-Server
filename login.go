package main

import (
	"net/http"
	"database/sql"
	"encoding/json"
)

func SignUp(w http.ResponseWriter, r *http.Request) {
	// Set the response header
	w.Header().Set("Content-Type", "application/json")

	// Connecting to the database
	db := dbConn()

	// Check if the user already exists
	mail := r.PostFormValue("mail")
	row := db.QueryRow("SELECT id FROM users WHERE mail = ?", mail).Scan()

	if row != sql.ErrNoRows {
		w.WriteHeader(http.StatusConflict)
		return
	}

	// Retrieve the other fields
	userType := r.PostFormValue("type")
	password := r.PostFormValue("password")
	firstName := r.PostFormValue("firstName")
	lastName := r.PostFormValue("lastName")
	birthDate := r.PostFormValue("birthDate")
	city := r.PostFormValue("city")
	phoneNumber := r.PostFormValue("phoneNumber")

	// Insertion
	res, err := db.Exec("INSERT INTO users (type, mail, password, firstName, lastName, birthDate, city, phoneNumber) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", userType, mail, password, firstName, lastName, birthDate, city, phoneNumber)
	if err != nil {
		w.WriteHeader(http.StatusNotModified)
		return
	}

	// Get the new user ID
	id, err := res.LastInsertId()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Format the response
	UserId := Id{id}
	json,_ := json.Marshal(UserId)

	// Send the response back
	w.WriteHeader(http.StatusCreated)
	w.Write(json)

}

func SignIn(w http.ResponseWriter, r *http.Request) {
	// Set the response header
	w.Header().Set("Content-Type", "application/json")

	// Connecting to the database
	db := dbConn()

	// Get the selected user
	var id, userType, firstName, lastName, birthDate, city, phoneNumber string
	mail := r.PostFormValue("mail")
	password := r.PostFormValue("password")
	row := db.QueryRow("SELECT id, type, firstName, lastName, birthDate, city, phoneNumber FROM users WHERE mail = ? AND password = ?", mail, password).Scan(&id, &userType, &firstName, &lastName, &birthDate, &city, &phoneNumber)

	if row == sql.ErrNoRows {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Format the response
	UserConnection := Connection{}
	UserConnection.Id = id
	UserConnection.Firstname = firstName
	UserConnection.Lastname = lastName
	UserConnection.Birthdate = birthDate
	UserConnection.City = city
	UserConnection.Mail = mail
	UserConnection.Tel = phoneNumber

	json,_ := json.Marshal(UserConnection)

	// Send the response back
	w.Write(json)
	w.WriteHeader(http.StatusOK)
}

func Forgot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db := dbConn()
	row := db.QueryRow("SELECT * FROM User WHERE mail = '" + r.PostFormValue("mail") + "'").Scan()
	if row != sql.ErrNoRows {
		passwd := NewPasswd{}
		passwd.Passwd = randomPasswd(12)
		json,_ := json.Marshal(passwd)
		w.Write(json)
		w.WriteHeader(http.StatusOK)
	} else{
		w.WriteHeader(http.StatusNotFound)
	}
}