// Author: Kévin Le
// Version 1.0

package prescription

import (
	"database/sql"
	"server/user"
	"encoding/json"
)

type Info struct {
	Id			int64			`json:"id"`
	User		user.Info		`json:"user"`
	Date		string			`json:"date"`
	Exercises	[]interface{}	`json:"exercises"`
}

func GetPrescriptions(db *sql.DB, userId int64) ([]Info, error) {
	query := `SELECT * FROM prescriptions WHERE userId = ?`

	rows, err := db.Query(query, userId)
	if err != nil {
		return nil, err
	}

	var userInfo user.Info
	_, err = user.GetFromId(db, &userInfo, userId)
	if err != nil {
		return nil, err
	}

	prescriptions := make([]Info, 0)

	for rows.Next() {
		prescription := Info{}
		var exercises []byte

		rows.Scan(&prescription.Id, &userId, &prescription.Date, &exercises)
		prescription.User = userInfo

		err = json.Unmarshal(exercises, &prescription.Exercises)
		if err != nil {
			continue
		}

		prescriptions = append(prescriptions, prescription)
	}

	return prescriptions, nil
}