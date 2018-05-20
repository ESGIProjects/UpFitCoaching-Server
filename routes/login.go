package routes

import (
	"net/http"
	"server/user"
	"encoding/json"
	"server/global"
	"strconv"
	"database/sql"
)

func ExistingMail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db := global.OpenDB()

	// Get field from request
	mail := r.URL.Query().Get("mail")

	// Get user from DB
	userInfo := user.Info{}
	_, err := user.GetFromMail(db, &userInfo, mail)

	// If user exists

	if err == nil {
		w.WriteHeader(http.StatusFound) // 302 plutôt que 409 ?
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func SignUp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db := global.OpenDB()

	// Get fields from request
	mail := r.PostFormValue("mail")
	userType, _ := strconv.Atoi(r.PostFormValue("type"))
	password := r.PostFormValue("password")
	firstName := r.PostFormValue("firstName")
	lastName := r.PostFormValue("lastName")
	city := r.PostFormValue("city")
	phoneNumber := r.PostFormValue("phoneNumber")

	// Potentially nil fields
	birthDate := r.PostFormValue("birthDate")
	address := r.PostFormValue("address")

	// Check if this user already exists
	userInfo := user.Info{}
	_, err := user.GetFromMail(db, &userInfo, mail)
	if err == nil {
		db.Close()

		global.SendError(w, "user_already_exists", http.StatusNotFound)
		return
	}

	// Start DB transaction
	tx, err := db.Begin()
	if err != nil {
		db.Close()

		println(err.Error())
		global.SendError(w, "internal_error", http.StatusNotModified)
		return
	}

	// Insert the user
	res, err := tx.Exec("INSERT INTO users (type, mail, password, firstName, lastName, city, phoneNumber) VALUES (?, ?, ?, ?, ?, ?, ?)", userType, mail, password, firstName, lastName, city, phoneNumber)
	if err != nil {
		tx.Rollback()
		db.Close()

		println(err.Error())
		global.SendError(w, "user_insert_failed", http.StatusNotModified)
		return
	}

	// Get the last inserted ID
	id, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		db.Close()

		println(err.Error())
		global.SendError(w, "internal_error", http.StatusInternalServerError)
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

		println(err.Error())
		global.SendError(w, "user_insert_failed", http.StatusNotModified)
		return
	}

	// Commit work to database
	if tx.Commit() != nil {
		db.Close()

		println(err.Error())
		global.SendError(w, "user_insert_failed", http.StatusNotModified)
		return
	}

	db.Close()

	// Format the response
	//userInfo = user.Info{}
	userInfo.Id = id

	global.SendJSON(w, userInfo, http.StatusCreated)
}

func SignIn(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db := global.OpenDB()

	userInfo := user.Info{}

	// Get fields from request
	userInfo.Mail = r.PostFormValue("mail")
	typedPassword := r.PostFormValue("password")

	password, err := user.GetFromMail(db, &userInfo, userInfo.Mail)
	if err != nil {
		println(err.Error())
		global.SendError(w, "user_not_exist", http.StatusNotFound)
		return
	}

	if password != typedPassword {
		global.SendError(w, "user_wrong_password", http.StatusNotFound)
		return
	}

	global.SendJSON(w, userInfo, http.StatusOK)
}

func Forgot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db := global.OpenDB()

	// Check if the user exists
	mail := r.PostFormValue("mail")
	row := db.QueryRow("SELECT * FROM users WHERE mail = ?", mail).Scan()

	if row == sql.ErrNoRows {
		global.SendError(w, "user_not_exist", http.StatusNotFound)
		return
	}

	// Format the response
	newPassword := global.NewPassword{global.RandomPasswd(12)}
	json, _ := json.Marshal(newPassword)

	// Send the response back
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}