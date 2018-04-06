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

	//var id string
	mail := r.PostFormValue("mail")

	// Check if the user already exists
	row := db.QueryRow("SELECT id FROM users WHERE mail = '" + mail + "'").Scan()

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

	id, err := res.LastInsertId()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	UserId := Id{id}
	json,_ := json.Marshal(UserId)

	w.WriteHeader(http.StatusCreated)
	w.Write(json)

}

func SignIn(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db := dbConn()
	var id, firstname, lastname, birthdate, city, mail, tel string
	row := db.QueryRow("SELECT id,firstname,lastname,birthdate,city,mail,tel FROM User WHERE mail ='" + r.PostFormValue("mail") + "' AND passwd ='" +
		r.PostFormValue("passwd") + "'").Scan(&id,&firstname,&lastname,&birthdate,&city,&mail,&tel)
	if row != sql.ErrNoRows {
		UserConnection := Connection{}
		UserConnection.Id = id
		UserConnection.Firstname = firstname
		UserConnection.Lastname = lastname
		UserConnection.Birthdate = birthdate
		UserConnection.City = city
		UserConnection.Mail = mail
		UserConnection.Tel = tel
		json,_ := json.Marshal(UserConnection)
		w.Write(json)
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
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