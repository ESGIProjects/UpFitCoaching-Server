package followUp

import (
	"server/user"
	"database/sql"
)

type Appraisal struct {
	Id					int64		`json:"id"`
	User				user.Info	`json:"user"`
	Date				string		`json:"date"`
	Goal				string		`json:"goal"`
	SessionsByWeek		*int		`json:"sessionsByWeek"`
	Contraindication	string		`json:"contraindication"`
	SportAntecedents	string		`json:"sportAntecedents"`
	HelpNeeded			*int		`json:"helpNeeded"`
	HasNutritionist		*int		`json:"hasNutritionist"`
	Comments			string		`json:"comments"`
}

type Measurements struct {
	Id					int64		`json:"id"`
	User				user.Info	`json:"user"`
	Date				string		`json:"date"`
	Weight				*int		`json:"weight"`
	Height				*int		`json:"height"`
	HipCircumference	*int		`json:"hipCircumference"`
	WaistCircumference	*int		`json:"waistCircumference"`
	ThighCircumference	*int		`json:"thighCircumference"`
	ArmCircumference	*int		`json:"armCircumference"`
}

type Test struct {
	Id							int64		`json:"id"`
	User						user.Info	`json:"user"`
	Date						string		`json:"date"`
	WarmUp						*float64	`json:"warmUp"`
	StartSpeed					*float64	`json:"startSpeed"`
	Increase					*float64	`json:"increase"`
	Frequency					*float64	`json:"frequency"`
	KneeFlexibility				*int		`json:"kneeFlexibility"`
	ShinFlexibility				*int		`json:"shinFlexibility"`
	HitFootFlexibility			*int		`json:"hitFootFlexibility"`
	ClosedFistGroundFlexibility	*int		`json:"closedFistGroundFlexibility"`
	HandFlatGroundFlexibility	*int		`json:"handFlatGroundFlexibility"`
}

func GetFromUserId(db *sql.DB, userId int64) (*Appraisal) {
	query := "SELECT * FROM appraisals WHERE userId = ?"

	appraisal := Appraisal{}
	var dbUserId int64

	row := db.QueryRow(query, userId).Scan(&appraisal.Id, &dbUserId, &appraisal.Date, &appraisal.Goal, &appraisal.SessionsByWeek, &appraisal.Contraindication, &appraisal.SportAntecedents, &appraisal.HelpNeeded, &appraisal.HasNutritionist, &appraisal.Comments)
	if row == sql.ErrNoRows {
		return nil
	}

	// Get user
	userInfo := user.Info{}
	_, err := user.GetFromId(db, &userInfo, userId)
	if err == nil {
		appraisal.User = userInfo
	}

	return &appraisal
}

func GetMeasurements(db *sql.DB, userId int64) ([]Measurements, error) {
	query := `SELECT * FROM measurements WHERE userId = ?`

	rows, err := db.Query(query, userId)
	if err != nil {
		return nil, err
	}

	var userInfo user.Info
	_, err = user.GetFromId(db, &userInfo, userId)
	if err != nil {
		return nil, err
	}

	measurements := make([]Measurements, 0)

	for rows.Next() {
		measurement := Measurements{}
		rows.Scan(&measurement.Id, &userId, &measurement.Date, &measurement.Weight, &measurement.Height, &measurement.HipCircumference, &measurement.WaistCircumference, &measurement.ThighCircumference, &measurement.ArmCircumference)
		measurement.User = userInfo

		measurements = append(measurements, measurement)
	}

	return measurements, nil
}

func GetTests(db *sql.DB, userId int64) ([]Test, error) {
	query := `SELECT * FROM tests WHERE userId = ?`

	rows, err := db.Query(query, userId)
	if err != nil {
		return nil, err
	}

	var userInfo user.Info
	_, err = user.GetFromId(db, &userInfo, userId)
	if err != nil {
		return nil, err
	}

	tests := make([]Test, 0)

	for rows.Next() {
		test := Test{}
		rows.Scan(&test.Id, &userId, &test.Date, &test.WarmUp, &test.StartSpeed, &test.Increase, &test.Frequency, &test.KneeFlexibility, &test.ShinFlexibility, &test.HitFootFlexibility, &test.ClosedFistGroundFlexibility, &test.HandFlatGroundFlexibility)
		test.User = userInfo

		tests = append(tests, test)
	}

	return tests, nil
}