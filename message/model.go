// Author: Jason Pierna
// Version 1.0

package message

import (
	"server/user"
	"database/sql"
)

type Info struct {
	Id				int64		`json:"id,omitempty"`
	Sender			user.Info	`json:"sender"`
	Receiver		user.Info	`json:"receiver"`
	Date			string		`json:"date"`
	Content			string		`json:"content"`
}

func get(db *sql.DB, query *sql.Row, messageInfo *Info) (error) {
	var sender, receiver int64

	row := query.Scan(&messageInfo.Id, &sender, &receiver, &messageInfo.Date, &messageInfo.Content)
	if row == sql.ErrNoRows {
		return row
	}

	// Sender
	senderInfo := user.Info{}
	_, err := user.GetFromId(db, &senderInfo, sender)
	if err != nil {
		return err
	}
	messageInfo.Sender = senderInfo

	// Receiver
	receiverInfo := user.Info{}
	_, err = user.GetFromId(db, &receiverInfo, sender)
	if err != nil {
		return err
	}
	messageInfo.Receiver = receiverInfo

	return nil
}

func GetFromUserId(db *sql.DB, id int) (*sql.Rows, error) {
	query := `SELECT * FROM messages
	WHERE sender = ? OR receiver = ?
	ORDER BY date DESC`

	return db.Query(query, id, id)
}

func GetUsersList(db *sql.DB, id int) (map[int64]user.Info, error) {
	query := `SELECT sender AS id FROM messages WHERE receiver = ?
	UNION SELECT receiver AS id FROM messages WHERE sender = ?
	UNION SELECT ? AS id;
	`

	rows, err := db.Query(query, id, id, id)
	if err != nil {
		return nil, err
	}

	return user.GetListFromQuery(db, rows)
}

func Save(db *sql.DB, messageInfo Info) (sql.Result, error) {
	query := `INSERT INTO messages
	(sender, receiver, date, content)
	VALUES (?, ?, ?, ?);
	`

	return db.Exec(query, messageInfo.Sender.Id, messageInfo.Receiver.Id, messageInfo.Date, messageInfo.Content)
}