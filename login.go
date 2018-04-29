package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
)

func CheckMail(w http.ResponseWriter, r *http.Request) {
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

	w.WriteHeader(http.StatusOK)
}

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
	city := r.PostFormValue("city")
	phoneNumber := r.PostFormValue("phoneNumber")

	// Potentially nil fields
	birthDate := r.PostFormValue("birthDate")
	address := r.PostFormValue("address")

	// Start transaction
	tx, err := db.Begin()
	if err != nil {
		db.Close()

		error := ErrorMessage{"internal_error"}
		json, _ := json.Marshal(error)

		w.WriteHeader(http.StatusNotModified)
		w.Write(json)
		return
	}

	// User insertion
	res, err := tx.Exec("INSERT INTO users (type, mail, password, firstName, lastName, city, phoneNumber) VALUES (?, ?, ?, ?, ?, ?, ?)", userType, mail, password, firstName, lastName, city, phoneNumber)
	if err != nil {
		tx.Rollback()
		db.Close()

		error := ErrorMessage{"user_insert_failed"}
		json, _ := json.Marshal(error)

		w.WriteHeader(http.StatusNotModified)
		w.Write(json)
		return
	}

	// Get the new user ID
	id, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		db.Close()

		error := ErrorMessage{"internal_error"}
		json, _ := json.Marshal(error)

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(json)
		return
	}

	// Insert information depending on user type
	if userType == 2 {
		_, err = tx.Exec("INSERT INTO coaches (id, address) VALUES(?, ?)", id, address)
	} else {
		_, err = tx.Exec("INSERT INTO clients (id, birthDate) VALUES(?, ?)", id, birthDate)
	}

	if err != nil {
		tx.Rollback()
		db.Close()

		error := ErrorMessage{"user_insert_failed"}
		json, _ := json.Marshal(error)

		w.WriteHeader(http.StatusNotModified)
		w.Write(json)
		return
	}

	// Commit to DB
	if tx.Commit() != nil {
		db.Close()

		error := ErrorMessage{"user_insert_failed"}
		json, _ := json.Marshal(error)

		w.WriteHeader(http.StatusNotModified)
		w.Write(json)
		return
	}

	db.Close()

	// Format the response
	user := UserInfo{}
	user.Id = id

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
	user := UserInfo{}
	user.Mail = r.PostFormValue("mail")
	password := r.PostFormValue("password")
	var dbPassword string
	var nullableBirthDate, nullableAddress sql.NullString

	row := db.QueryRow("SELECT * FROM users NATURAL LEFT JOIN coaches NATURAL LEFT JOIN clients WHERE mail = ?", user.Mail).Scan(&user.Id, &user.Type, &user.Mail, &dbPassword, &user.FirstName, &user.LastName, &user.City, &user.PhoneNumber, &nullableAddress, &nullableBirthDate)

	// If user does not exist
	if row == sql.ErrNoRows {
		error := ErrorMessage{"user_not_exist"}
		json, _ := json.Marshal(error)
		w.WriteHeader(http.StatusNotFound)
		w.Write(json)
		return
	}

	if password != dbPassword {
		error := ErrorMessage{"user_wrong_password"}
		json, _ := json.Marshal(error)
		w.WriteHeader(http.StatusNotFound)
		w.Write(json)
		return
	}

	// Adds optional value
	user.Address = nullableAddress.String
	user.BirthDate = nullableBirthDate.String

	// Format the response
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
