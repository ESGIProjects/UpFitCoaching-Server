// Author: Kévin Le
// Version 1.0

package event

import (
	"server/user"
	"database/sql"
)

type Info struct {
	Id			int64		`json:"id"`
	Name		string		`json:"name"`
	Type		int			`json:"type"`
	Status		int			`json:"status"`
	FirstUser	user.Info	`json:"firstUser"`
	SecondUser	user.Info	`json:"secondUser"`
	Start		string		`json:"start"`
	End			string		`json:"end"`
	Created		string		`json:"created"`
	CreatedBy	user.Info	`json:"createdBy"`
	Updated		string		`json:"updated"`
	UpdatedBy	user.Info	`json:"updatedBy"`
}

func get(db *sql.DB, query *sql.Row, eventInfo *Info) (error) {
	var firstUser, secondUser, createdBy, updatedBy int64

	row := query.Scan(&eventInfo.Id, &eventInfo.Name, &eventInfo.Type, &eventInfo.Status, &firstUser, &secondUser, &eventInfo.Start, &eventInfo.End, &eventInfo.Created, &createdBy, &eventInfo.Updated, &updatedBy)
	if row == sql.ErrNoRows {
		return row
	}

	// First user
	firstUserInfo := user.Info{}
	_, err := user.GetFromId(db, &firstUserInfo, firstUser)
	if err != nil {
		return err
	}
	eventInfo.FirstUser = firstUserInfo

	// Second user
	secondUserInfo := user.Info{}
	_, err = user.GetFromId(db, &secondUserInfo, secondUser)
	if err != nil {
		return err
	}
	eventInfo.SecondUser = secondUserInfo

	// Created by
	if createdBy == firstUser {
		eventInfo.CreatedBy = firstUserInfo
	} else {
		eventInfo.CreatedBy = secondUserInfo
	}

	// Updated by
	if updatedBy == firstUser {
		eventInfo.UpdatedBy = firstUserInfo
	} else {
		eventInfo.UpdatedBy = secondUserInfo
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
	WHERE firstUser = ? OR secondUser = ?
	ORDER BY start`

	return db.Query(query, id, id)
}

func GetUsersList(db *sql.DB, id int) (map[int64]user.Info, error) {
	query := `SELECT firstUser AS id FROM events WHERE secondUser = ?
	UNION SELECT secondUser AS id FROM events WHERE firstUser = ?
	UNION SELECT ? AS id;`

	rows, err := db.Query(query, id, id, id)
	if err != nil {
		return nil, err
	}

	return user.GetListFromQuery(db, rows)
}

func Save(db *sql.DB, eventInfo Info) (sql.Result, error) {
	query := `INSERT INTO events
	(name, type, status, firstUser, secondUser, start, end, created, createdBy, updated, updatedBy)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`

	return db.Exec(query, eventInfo.Name, eventInfo.Type, eventInfo.Status, eventInfo.FirstUser.Id, eventInfo.SecondUser.Id, eventInfo.Start, eventInfo.End, eventInfo.Created, eventInfo.CreatedBy.Id, eventInfo.Updated, eventInfo.UpdatedBy.Id)
}

func Update(db *sql.DB, eventInfo Info) (sql.Result, error) {
	query := `UPDATE events
	SET name = ?, type = ?, status = ?, firstUser = ?, secondUser = ?, start = ?, end = ?, created = ?, createdBy = ?, updated = ?, updatedBy = ?
	WHERE id = ?;`

	return db.Exec(query, eventInfo.Name, eventInfo.Type, eventInfo.Status, eventInfo.FirstUser.Id, eventInfo.SecondUser.Id, eventInfo.Start, eventInfo.End, eventInfo.Created, eventInfo.CreatedBy.Id, eventInfo.Updated, eventInfo.UpdatedBy.Id, eventInfo.Id)
}