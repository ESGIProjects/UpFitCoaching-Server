package routes

import (
	"net/http"
	"server/global"
	"strconv"
	"server/followUp"
)

func GetLastAppraisal(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	db := global.OpenDB()
	defer db.Close()

	// Get userId from request
	userId, err := strconv.Atoi(r.URL.Query().Get("userId"))
	if err != nil {

	}

	// Retrieve last appraisal from user
	appraisal := followUp.GetFromUserId(db, int64(userId))

	if appraisal == nil {
		global.SendError(w, "appraisal_not_found", http.StatusNoContent)
	} else {
		global.SendJSON(w, appraisal, http.StatusOK)
	}
}

func GetMeasurements(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	db := global.OpenDB()
	defer db.Close()

	// Get userId from request
	userId, err := strconv.Atoi(r.URL.Query().Get("userId"))
	if err != nil {
		print(err.Error())
		db.Close()
		global.SendError(w, "parameter_error", http.StatusBadRequest)

		return
	}

	// Retrieve measurements from user
	measurements, err := followUp.GetAllMeasurements(db, int64(userId))
	if err != nil {
		print(err.Error())
		db.Close()
		global.SendError(w, "internal_error", http.StatusInternalServerError)

		return
	}

	global.SendJSON(w, measurements, http.StatusOK)
}