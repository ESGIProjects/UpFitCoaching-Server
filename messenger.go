package main

import (
	"net/http"
	"encoding/json"
)

func GetMessages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db := dbConn()
	var id,idSender,idReceiver,body,timestamp string
	test,_ := r.URL.Query()["id"]
	rows,_ := db.Query("SELECT * FROM Message WHERE idSender='" + test[0] + "' OR idReceiver ='" +
		test[0] + "'")
	defer rows.Close()
	Messages := []Message{}
	for rows.Next(){
		Message := Message{}
		rows.Scan(&id,&idSender,&idReceiver,&body,&timestamp)
		Message.Id = id
		Message.IdSender = idSender
		Message.IdReceiver = idReceiver
		Message.Body = body
		Message.Timestamp = timestamp
		Messages = append(Messages, Message)
	}
	json,_ := json.Marshal(Messages)
	w.Write(json)
	w.WriteHeader(http.StatusOK)
}
func SendMessage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db := dbConn()
	_, err := db.Query("INSERT INTO Message (idSender,idReceiver,body,timestamp) VALUES ('" + r.PostFormValue("idSender") +
		"','" +r.PostFormValue("idReceiver") + "','" + r.PostFormValue("body") + "','" +r.PostFormValue("timestamp") + "')")
	if err == nil{
		w.WriteHeader(http.StatusCreated)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}

}
func DeleteMessages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	db := dbConn()
	test,_ := r.URL.Query()["id"]
	_, err := db.Query("DELETE FROM Message WHERE id='" + test[0] + "'")
	if err == nil{
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}
