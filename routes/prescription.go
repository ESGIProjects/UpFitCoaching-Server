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