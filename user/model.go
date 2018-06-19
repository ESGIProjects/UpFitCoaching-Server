package user

import (
	"database/sql"
	"errors"
)

type Info struct {
	Id			int64	`json:"id"`
	Type		*int	`json:"type"`
	Mail		string	`json:"mail"`
	FirstName	string	`json:"firstName"`
	LastName	string	`json:"lastName"`
	Sex			*int	`json:"sex"`
	City		string	`json:"city"`
	PhoneNumber	string	`json:"phoneNumber"`
	Address		string	`json:"address,omitempty"`
	BirthDate	string	`json:"birthDate,omitempty"`
	Coach		*Info	`json:"coach,omitempty"`
}

func get(db *sql.DB, query *sql.Row, userInfo *Info) (string, error) {
	var password string
	var birthDate, address sql.NullString
	var coachId sql.NullInt64

	// Make the query
	row := query.Scan(&userInfo.Id, &userInfo.Type, &userInfo.Mail, &password, &userInfo.FirstName, &userInfo.LastName, &userInfo.Sex, &userInfo.City, &userInfo.PhoneNumber, &address, &birthDate, &coachId)

	// Does the user exist?
	if row == sql.ErrNoRows {
		return "", errors.New("no_user")
	}

	// Updates model with nullable fields
	userInfo.Address = address.String
	userInfo.BirthDate = birthDate.String

	if coachId.Valid {
		// Create a new UserInfo struct with the coach data
		coach := Info{}
		_, err := GetFromId(db, &coach, coachId.Int64)

		if err == nil {
			userInfo.Coach = &coach
		}
	}

	// Returns password for login operations, and no error
	return password, nil
}

func GetFromId(db *sql.DB, userInfo *Info, id int64) (string, error) {
	query := `SELECT * FROM users
	NATURAL LEFT JOIN coaches
	NATURAL LEFT JOIN clients
	WHERE id = ?
	`

	row := db.QueryRow(query, id)

	return get(db, row, userInfo)
}

func GetFromMail(db *sql.DB, userInfo *Info, mail string) (string, error) {
	query := `SELECT * FROM users
	NATURAL LEFT JOIN coaches
	NATURAL LEFT JOIN clients
	WHERE mail = ?
	`

	row := db.QueryRow(query, mail)

	return get(db, row, userInfo)
}

func GetListFromQuery(db *sql.DB, rows *sql.Rows) (map[int64]Info, error) {
	users := make(map[int64]Info)

	for rows.Next() {
		var userId int64
		rows.Scan(&userId)

		userInfo := Info{}
		_, err := GetFromId(db, &userInfo, userId)
		if err != nil {
			print(err.Error())
			continue
		}

		users[userId] = userInfo
	}

	return users, nil
}