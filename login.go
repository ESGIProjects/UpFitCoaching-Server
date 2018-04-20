package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
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
		error := ErrorMessage{"user_already_exists"}
		json, _ := json.Marshal(error)

		w.WriteHeader(http.StatusConflict)
		w.Write(json)
		return
	}

	// Retrieve the other fields
	userType, _ := strconv.Atoi(r.PostFormValue("type"))
	password := r.PostFormValue("password")
	firstName := r.PostFormValue("firstName")
	lastName := r.PostFormValue("lastName")
	birthDate := r.PostFormValue("birthDate")
	city := r.PostFormValue("city")
	phoneNumber := r.PostFormValue("phoneNumber")

	// Insertion
	res, err := db.Exec("INSERT INTO users (type, mail, password, firstName, lastName, birthDate, city, phoneNumber) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", userType, mail, password, firstName, lastName, birthDate, city, phoneNumber)
	if err != nil {
		error := ErrorMessage{"user_insert_failed"}
		json, _ := json.Marshal(error)

		w.WriteHeader(http.StatusNotModified)
		w.Write(json)
		return
	}

	// Get the new user ID
	id, err := res.LastInsertId()
	if err != nil {
		error := ErrorMessage{"internal_error"}
		json, _ := json.Marshal(error)

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(json)
		return
	}

	// Format the response
	user := UserInfo{}
	user.Id = id
	user.UserType = userType

	json, _ := json.Marshal(user)

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
	var id int64
	var userType int
	var firstName, lastName, address, birthDate, city, phoneNumber string
	var UserRow, CoachRow error

	mail := r.PostFormValue("mail")
	password := r.PostFormValue("password")

	UserRow = db.QueryRow("SELECT id, type, mail, firstName, lastName, birthDate, city, phoneNumber FROM users WHERE mail = ? AND password = ?", mail, password).Scan(&id, &userType, &mail, &firstName, &lastName, &birthDate, &city, &phoneNumber)

	// If user does not exist, another request is made on coaches table
	if UserRow == sql.ErrNoRows {
		CoachRow = db.QueryRow("SELECT id, mail, firstName, lastName, address, city, phoneNumber FROM coaches WHERE mail = ? AND password = ?", mail, password).Scan(&id, &mail, &firstName, &lastName, &address, &city, &phoneNumber)
		// If coach does not exist
		if CoachRow == sql.ErrNoRows {
			error := ErrorMessage{"user_not_exist"}
			json, _ := json.Marshal(error)
			w.WriteHeader(http.StatusNotFound)
			w.Write(json)
			return
		}
	}

	// Format the response
	user := UserInfo{}
	user.Id = id
	user.UserType = userType
	user.Mail = mail
	user.FirstName = firstName
	user.LastName = lastName
	user.BirthDate = birthDate
	user.City = city
	user.PhoneNumber = phoneNumber

	json, _ := json.Marshal(user)

	// Send the response back
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

func Forgot(w http.ResponseWriter, r *http.Request) {
	// Set the response header
	w.Header().Set("Content-Type", "application/json")

	// Connecting to the database
	db := dbConn()

	// Check if the user exists
	mail := r.PostFormValue("mail")
	row := db.QueryRow("SELECT * FROM users WHERE mail = ?", mail).Scan()

	if row == sql.ErrNoRows {
		error := ErrorMessage{"user_not_exist"}
		json, _ := json.Marshal(error)

		w.WriteHeader(http.StatusNotFound)
		w.Write(json)
		return
	}

	// Format the response
	newPassword := NewPassword{randomPasswd(12)}
	json, _ := json.Marshal(newPassword)

	// Send the response back
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}
