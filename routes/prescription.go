package routes

import (
	"net/http"
	"server/global"
	"strconv"
	"server/prescription"
)

func GetPrescriptions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	db := global.OpenDB()
	defer db.Close()

	// Get userId from request
	userId, err := strconv.Atoi(r.URL.Query().Get("userId"))
	if err != nil {
		println(err.Error())
		db.Close()

		global.SendError(w, "parameter_error", http.StatusBadRequest)
		return
	}

	// Retrieve prescriptions from user
	prescriptions, err := prescription.GetPrescriptions(db, int64(userId))
	if err != nil {
		println(err.Error())
		db.Close()

		global.SendError(w, "internal_error", http.StatusInternalServerError)
		return
	}

	global.SendJSON(w, prescriptions, http.StatusOK)
}

func CreatePrescription(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	db := global.OpenDB()
	defer db.Close()

	// Get fields from request
	userId, _ := strconv.Atoi(r.PostFormValue("userId"))
	date := r.PostFormValue("date")
	exercises := r.PostFormValue("exercises")

	// Inserting prescription into DB
	res, err := db.Exec("INSERT INTO prescriptions (userId, date, exercises) VALUES (?, ?, ?)", userId, date, exercises)
	if err != nil {
		db.Close()

		print(err.Error())
		global.SendError(w, "prescription_insert_failed", http.StatusNotModified)
		return
	}

	// Get the new prescription ID
	prescriptionId, err := res.LastInsertId()
	if err != nil {
		db.Close()

		print(err.Error())
		global.SendError(w, "internal_error", http.StatusInternalServerError)
		return
	}

	json := make(map[string]int64)
	json["id"] = prescriptionId

	global.SendJSON(w, json, http.StatusCreated)
}