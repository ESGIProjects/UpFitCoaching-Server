package event

import (
	"server/user"
	"database/sql"
)

type Info struct {
	Id			int64		`json:"id,omitempty"`
	Name		string		`json:"name,omitempty"`
	Type		*int		`json:"type,omitempty"`
	Client		user.Info	`json:"client,omitempty"`
	Coach		user.Info	`json:"coach,omitempty"`
	Start		string		`json:"start,omitempty"`
	End			string		`json:"end,omitempty"`
	Created		string		`json:"created,omitempty"`
	CreatedBy	user.Info	`json:"createdBy,omitempty"`
	Updated		string		`json:"updated,omitempty"`
	UpdatedBy	user.Info	`json:"updatedBy,omitempty"`
}

func get(db *sql.DB, query *sql.Row, eventInfo *Info) (error) {
	var coach, client, createdBy, updatedBy int64

	row := query.Scan(&eventInfo.Id, &eventInfo.Name, &eventInfo.Type, &client, &coach, &eventInfo.Start, &eventInfo.End, &eventInfo.Created, &createdBy, &eventInfo.Updated, &updatedBy)
	if row == sql.ErrNoRows {
		return row
	}

	// Client
	clientInfo := user.Info{}
	_, err := user.GetFromId(db, &clientInfo, client)
	if err != nil {
		return err
	}
	eventInfo.Client = clientInfo

	// Coach
	coachInfo := user.Info{}
	_, err = user.GetFromId(db, &coachInfo, coach)
	if err != nil {
		return err
	}
	eventInfo.Coach = coachInfo

	// Created by
	if createdBy == client {
		eventInfo.CreatedBy = clientInfo
	} else {
		eventInfo.CreatedBy = coachInfo
	}

	// Updated by
	if updatedBy == client {
		eventInfo.UpdatedBy = clientInfo
	} else {
		eventInfo.UpdatedBy = coachInfo
	}

	return nil
}

func GetFromId(db *sql.DB, eventInfo *Info, id int64) (error) {
	query := `SELECT * FROM events
	WHERE id = ?
	`

	row := db.QueryRow(query, id)

	return get(db, row, eventInfo)
}

func GetFromUserId(db *sql.DB, id int) (*sql.Rows, error) {
	query := `SELECT * FROM events
	WHERE client = ? OR coach = ?
	ORDER BY start`

	return db.Query(query, id, id)
}

func GetConcernedUsers(db *sql.DB, id int) (*sql.Rows, error) {
	query := `SELECT id, type, mail, firstName, lastName, city, phoneNumber, address, birthDate
	FROM users
	NATURAL LEFT JOIN coaches
	NATURAL LEFT JOIN clients
	WHERE id IN
	(SELECT client AS id FROM events WHERE coach = ?
	UNION SELECT coach AS id FROM events WHERE client = ?
	UNION SELECT ? as id);`

	return db.Query(query, id, id, id)
}