package main

import "server/user"

// Messages

type Message struct {
	Id				int64			`json:"id,omitempty"`
	Sender			user.Info	`json:"sender"`
	Receiver		user.Info	`json:"receiver"`
	Date			string			`json:"date"`
	Content			string			`json:"content"`
}