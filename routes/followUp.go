// Author: Jason Pierna
// Version 1.0

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
		println(err.Error())
		db.Close()

		global.SendError(w, "parameter_error", http.StatusBadRequest)
		return
	}

	// Retrieve measurements from user
	measurements, err := followUp.GetMeasurements(db, int64(userId))
	if err != nil {
		println(err.Error())
		db.Close()

		global.SendError(w, "internal_error", http.StatusInternalServerError)
		return
	}

	global.SendJSON(w, measurements, http.StatusOK)
}

func GetTests(w http.ResponseWriter, r *http.Request) {
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

	// Retrieve tests for user
	tests, err := followUp.GetTests(db, int64(userId))
	if err != nil {
		println(err.Error())
		db.Close()

		global.SendError(w, "internal_error", http.StatusInternalServerError)
		return
	}

	global.SendJSON(w, tests, http.StatusOK)
}

func CreateAppraisal(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	db := global.OpenDB()
	defer db.Close()

	// Get fields from request
	userId, _ := strconv.Atoi(r.PostFormValue("userId"))
	date := r.PostFormValue("date")
	goal := r.PostFormValue("goal")
	sessionsByWeek, _ := strconv.Atoi(r.PostFormValue("sessionsByWeek"))
	contraindication := r.PostFormValue("contraindication")
	sportAntecedents := r.PostFormValue("sportAntecedents")
	helpNeeded, _ := strconv.ParseBool(r.PostFormValue("helpNeeded"))
	hasNutritionist, _ := strconv.ParseBool(r.PostFormValue("hasNutritionist"))
	comments := r.PostFormValue("comments")

	// Inserting appraisal into DB
	res, err := db.Exec("INSERT INTO appraisals (userId, date, goal, sessionsByWeek, contraindication, sportAntecedents, helpNeeded, hasNutritionist, comments) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)", userId, date, goal, sessionsByWeek, contraindication, sportAntecedents, helpNeeded, hasNutritionist, comments)
	if err != nil {
		db.Close()

		print(err.Error())
		global.SendError(w, "appraisal_insert_failed", http.StatusNotModified)
		return
	}

	// Get the new appraisal ID
	appraisalId, err := res.LastInsertId()
	if err != nil {
		db.Close()

		print(err.Error())
		global.SendError(w, "internal_error", http.StatusInternalServerError)
		return
	}

	json := make(map[string]int64)
	json["id"] = appraisalId

	global.SendJSON(w, json, http.StatusCreated)
}

func CreateMeasurements(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	db := global.OpenDB()
	defer db.Close()

	// Get fields from request
	userId, _ := strconv.Atoi(r.PostFormValue("userId"))
	date := r.PostFormValue("date")
	weight, _ := strconv.ParseFloat(r.PostFormValue("weight"), 64)
	height, _ := strconv.ParseFloat(r.PostFormValue("height"), 64)
	hipCircumference, _ := strconv.ParseFloat(r.PostFormValue("hipCircumference"), 64)
	waistCircumference, _ := strconv.ParseFloat(r.PostFormValue("waistCircumference"), 64)
	thighCircumference, _ := strconv.ParseFloat(r.PostFormValue("thighCircumference"), 64)
	armCircumference, _ := strconv.ParseFloat(r.PostFormValue("armCircumference"), 64)

	// Inserting measurements into DB
	res, err := db.Exec("INSERT INTO measurements (userId, date, weight, height, hipCircumference, waistCircumference, thighCircumference, armCircumference) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", userId, date, weight, height, hipCircumference, waistCircumference, thighCircumference, armCircumference)
	if err != nil {
		db.Close()

		print(err.Error())
		global.SendError(w, "measurements_insert_failed", http.StatusNotModified)
		return
	}

	// Get the new measurements ID
	measurementsId, err := res.LastInsertId()
	if err != nil {
		db.Close()

		print(err.Error())
		global.SendError(w, "internal_error", http.StatusInternalServerError)
		return
	}

	json := make(map[string]int64)
	json["id"] = measurementsId

	global.SendJSON(w, json, http.StatusCreated)
}

func CreateTest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	db := global.OpenDB()
	defer db.Close()

	// Get fields from request
	userId, _ := strconv.Atoi(r.PostFormValue("userId"))
	date := r.PostFormValue("date")

	warmUp, _ := strconv.ParseFloat(r.PostFormValue("warmUp"), 64)
	startSpeed, _ := strconv.ParseFloat(r.PostFormValue("startSpeed"), 64)
	increase, _ := strconv.ParseFloat(r.PostFormValue("increase"), 64)
	frequency, _ := strconv.ParseFloat(r.PostFormValue("frequency"), 64)
	kneeFlexibility, _ := strconv.Atoi(r.PostFormValue("kneeFlexibility"))
	shinFlexibility, _ := strconv.Atoi(r.PostFormValue("shinFlexibility"))
	hitFootFlexibility, _ := strconv.Atoi(r.PostFormValue("hitFootFlexibility"))
	closedFistGroundFlexibility, _ := strconv.Atoi(r.PostFormValue("closedFistGroundFlexibility"))
	handFlatGroundFlexibility, _ := strconv.Atoi(r.PostFormValue("handFlatGroundFlexibility"))

	// Inserting test into DB
	res, err := db.Exec("INSERT INTO tests (userId, date, warmUp, startSpeed, increase, frequency, kneeFlexibility, shinFlexibility, hitFootFlexibility, closedFistGroundFlexibility, handFlatGroundFlexibility) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", userId, date, warmUp, startSpeed, increase, frequency, kneeFlexibility, shinFlexibility, hitFootFlexibility, closedFistGroundFlexibility, handFlatGroundFlexibility)
	if err != nil {
		db.Close()

		print(err.Error())
		global.SendError(w, "test_insert_failed", http.StatusNotModified)
		return
	}

	// Get the new test ID
	testId, err := res.LastInsertId()
	if err != nil {
		db.Close()

		print(err.Error())
		global.SendError(w, "internal_error", http.StatusInternalServerError)
		return
	}

	json := make(map[string]int64)
	json["id"] = testId

	global.SendJSON(w, json, http.StatusCreated)
}
