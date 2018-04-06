package main

import (
	"net/http"
	"database/sql"
	"encoding/json"
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

		w.Write(json)
		w.WriteHeader(http.StatusConflict)
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

		w.Write(json)
		w.WriteHeader(http.StatusNotModified)
		return
	}

	// Get the new user ID
	id, err := res.LastInsertId()
	if err != nil {
		error := ErrorMessage{"internal_error"}
		json, _ := json.Marshal(error)

		w.Write(json)
		w.WriteHeader(http.StatusInternalServerError)
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
	var firstName, lastName, birthDate, city, phoneNumber string

	mail := r.PostFormValue("mail")
	password := r.PostFormValue("password")

	println("Mail from call: ", mail)

	row := db.QueryRow("SELECT id, type, mail, firstName, lastName, birthDate, city, phoneNumber FROM users WHERE mail = ? AND password = ?", mail, password).Scan(&id, &userType, &mail, &firstName, &lastName, &birthDate, &city, &phoneNumber)

	print("Row value: ", row)

	// If user does not exist
	if row == sql.ErrNoRows {
		print("sqlErrNoRows !!!")

		error := ErrorMessage{"user_not_exist"}
		json, _ := json.Marshal(error)

		w.WriteHeader(http.StatusNotFound)
		w.Write(json)
		return
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
	w.Write(json)
	w.WriteHeader(http.StatusOK)
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

		w.Write(json)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Format the response
	newPassword := NewPassword{randomPasswd(12)}
	json, _ := json.Marshal(newPassword)

	// Send the response back
	w.Write(json)
	w.WriteHeader(http.StatusOK)
}