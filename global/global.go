package global

import (
	"database/sql"
	"math/rand"
	"time"
	"encoding/json"
	"net/http"
)

type ErrorMessage struct {
	Message	string	`json:"message"`
}

type NewPassword struct {
	Password	string
}

func OpenDB() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := "root"
	dbName := "upfit"

	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}

	db.Exec("SET NAMES utf8mb4")

	return db
}

func RandomPasswd(strlen int) string {
	var random = rand.New(rand.NewSource(time.Now().UnixNano()))
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ@&-_0123456789"
	passwd := make([]byte, strlen)
	for i := range passwd {
		passwd[i] = chars[random.Intn(len(chars))]
	}
	return string(passwd)
}

func SendJSON(w http.ResponseWriter, v interface{}, statusCode int) {
	convertedJSON, _ := json.Marshal(v)
	w.WriteHeader(statusCode)
	w.Write(convertedJSON)
}

func SendError(w http.ResponseWriter, message string, statusCode int) {
	errorMessage := ErrorMessage{message}
	SendJSON(w, errorMessage, statusCode)
}