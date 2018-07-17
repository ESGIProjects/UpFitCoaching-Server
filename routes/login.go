// Author: Kévin Le
// Version 1.0

package routes

import (
	"net/http"
	"server/user"
	"encoding/json"
	"server/global"
	"strconv"
	"database/sql"
	"server/message"
	"time"
	"server/auth"
)

const uniqueCoachId = 15

func ExistingMail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	db := global.OpenDB()

	// Get field from request
	mail := r.PostFormValue("mail")

	// Get user from DB
	userInfo := user.Info{}
	_, err := user.GetFromMail(db, &userInfo, mail)

	// If user exists
	if err != nil {
		print(err.Error())
		w.WriteHeader(http.StatusOK)
	} else {
		global.SendError(w, "user_already_exists", http.StatusFound)
	}
}

func SignUp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	db := global.OpenDB()

	// Get fields from request
	mail := r.PostFormValue("mail")
	userType, _ := strconv.Atoi(r.PostFormValue("type"))
	password := r.PostFormValue("password")
	firstName := r.PostFormValue("firstName")
	lastName := r.PostFormValue("lastName")
	sex, _ := strconv.Atoi(r.PostFormValue("sex"))
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

		global.SendError(w, "user_already_exists", http.StatusFound)
		return
	}

	// Start DB transaction
	tx, err := db.Begin()
	if err != nil {
		db.Close()

		println(err.Error())
		global.SendError(w, "user_insert_failed", http.StatusNotModified)
		return
	}

	// Insert the user
	res, err := tx.Exec("INSERT INTO users (type, mail, password, firstName, lastName, sex, city, phoneNumber) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", userType, mail, password, firstName, lastName, sex, city, phoneNumber)
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
		global.SendError(w, "user_insert_failed", http.StatusNotModified)
		return
	}

	// Insert information depending on user type
	if userType == 2 {
		_, err = tx.Exec("INSERT INTO coaches (id, address) VALUES(?, ?)", id, address)
	} else {
		_, err = tx.Exec("INSERT INTO clients (id, birthDate, coachId) VALUES(?, ?, ?)", id, birthDate, uniqueCoachId)
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

	// Generate auth token
	token, err := auth.CreateToken(id)
	if err != nil {
		db.Close()

		println(err.Error())
		global.SendError(w, "user_insert_failed", http.StatusNotModified)
		return
	}

	// Format the response
	json := make(map[string]interface{})
	json["id"] = id
	json["token"] = token
	userInfo.Id = id

	// If client, retrieve coach data and send first message
	if userType == 0 {
		coach := user.Info{}
		user.GetFromId(db, &coach, uniqueCoachId)
		userInfo.Coach = &coach
		json["coach"] = coach

		currentTime := time.Now().Local()

		messageInfo := message.Info{}
		messageInfo.Sender = userInfo
		messageInfo.Receiver = coach
		messageInfo.Date = currentTime.Format("2006-01-02 15:04:05")
		messageInfo.Content = "Nouveau client"

		message.Save(db, messageInfo)
	}

	db.Close()
	global.SendJSON(w, json, http.StatusCreated)
}

func SignIn(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
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
		global.SendError(w, "user_wrong_password", http.StatusBadRequest)
		return
	}

	// Generate auth token
	token, err := auth.CreateToken(userInfo.Id)
	if err != nil {
		db.Close()

		println(err.Error())
		global.SendError(w, "token_creation_error", http.StatusInternalServerError)
		return
	}

	json := make(map[string]interface{})
	json["token"] = token
	json["user"] = userInfo

	global.SendJSON(w, json, http.StatusOK)
}

func UpdateProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	db := global.OpenDB()
	defer db.Close()

	// Get user ID from request
	userId, _ := strconv.Atoi(r.PostFormValue("userId"))

	// Retrieve user from DB
	userInfo := user.Info{}
	_, err := user.GetFromId(db, &userInfo, int64(userId))
	if err != nil {
		db.Close()

		println(err.Error())
		global.SendError(w, "parameter_error", http.StatusBadRequest)
		return
	}

	// Update fields from request
	userInfo.Mail = r.PostFormValue("mail")
	userInfo.FirstName = r.PostFormValue("firstName")
	userInfo.LastName = r.PostFormValue("lastName")
	userInfo.City = r.PostFormValue("city")
	userInfo.PhoneNumber = r.PostFormValue("phoneNumber")

	// Potentially nil fields
	password := r.PostFormValue("password")
	address := r.PostFormValue("address")

	// Start DB transaction
	tx, err := db.Begin()
	if err != nil {
		db.Close()

		println(err.Error())
		global.SendError(w, "user_insert_failed", http.StatusNotModified)
		return
	}

	// Update the user
	_, err = tx.Exec("UPDATE users SET mail = ?, firstName = ?, lastName = ?, city = ?, phoneNumber = ? WHERE id = ?", userInfo.Mail, userInfo.FirstName, userInfo.LastName, userInfo.City, userInfo.PhoneNumber, userInfo.Id)
	if err != nil {
		tx.Rollback()
		db.Close()

		println(err.Error())
		global.SendError(w, "user_insert_failed", http.StatusNotModified)
		return
	}

	// Update the password if needed
	if password != "" {
		_, err = tx.Exec("UPDATE users SET password = ? WHERE id = ?", password, userInfo.Id)
		if err != nil {
			tx.Rollback()
			db.Close()

			println(err.Error())
			global.SendError(w, "user_insert_failed", http.StatusNotModified)
			return
		}
	}

	// Update the coach's adress if needed
	if address != "" {
		_, err = tx.Exec("UPDATE coaches SET address = ? WHERE id = ?", address, userInfo.Id)
		if err != nil {
			tx.Rollback()
			db.Close()

			println(err.Error())
			global.SendError(w, "user_insert_failed", http.StatusNotModified)
			return
		}

		userInfo.Address = address
	}

	// Commit work to database
	if tx.Commit() != nil {
		db.Close()

		println(err.Error())
		global.SendError(w, "user_insert_failed", http.StatusNotModified)
		return
	}

	// Send the updated user
	global.SendJSON(w, userInfo, http.StatusOK)
}

func Forgot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
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