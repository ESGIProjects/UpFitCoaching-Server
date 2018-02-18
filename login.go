package main

import (
	"net/http"
	"database/sql"
	"encoding/json"
)

func SignUp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db := dbConn()
	var id string
	checkId := db.QueryRow("SELECT id FROM User WHERE mail ='" + r.PostFormValue("mail") + "'").Scan()
	if checkId == sql.ErrNoRows {
		db.QueryRow("INSERT INTO User (type, firstName, lastName, passwd, birthdate, city, mail, tel) VALUES" +
			"('" + r.PostFormValue("type") + "','"  + r.PostFormValue("firstName") + "','" + r.PostFormValue("lastName") +
			"','" + r.PostFormValue("passwd") + "','" + r.PostFormValue("birthdate") + "','" + r.PostFormValue("city") +
			"','" + r.PostFormValue("mail") + "','" + r.PostFormValue("tel") + "')")
		db.QueryRow("SELECT id FROM User WHERE mail = '" + r.PostFormValue("mail") + "'").Scan(&id)
		UserId := Id{}
		UserId.Id = id
		json,_ := json.Marshal(UserId)
		w.Write(json)
		w.WriteHeader(http.StatusCreated)
	} else {
		w.WriteHeader(http.StatusConflict)
	}
}

func SignIn(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db := dbConn()
	var id, firstname, lastname, birthdate, city, mail, tel string
	check := db.QueryRow("SELECT id,firstname,lastname,birthdate,city,mail,tel FROM User WHERE mail ='" + r.PostFormValue("mail") + "' AND passwd ='" +
		r.PostFormValue("passwd") + "'").Scan(&id,&firstname,&lastname,&birthdate,&city,&mail,&tel)
	if check != sql.ErrNoRows {
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
	check := db.QueryRow("SELECT * FROM User WHERE mail = '" + r.PostFormValue("mail") + "'").Scan()
	if check != sql.ErrNoRows {
		passwd := NewPasswd{}
		passwd.Passwd = randomPasswd(12)
		json,_ := json.Marshal(passwd)
		w.Write(json)
		w.WriteHeader(http.StatusOK)
	} else{
		w.WriteHeader(http.StatusNotFound)
	}
}
