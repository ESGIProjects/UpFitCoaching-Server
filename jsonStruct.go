package main

// Misc

type ErrorMessage struct {
	Message	string	`json:"message"`
}

// User management

type UserInfo struct {
	Id			int64	`json:"id"`
	Type		int		`json:"type"`
	Mail		string	`json:"mail,omitempty"`
	FirstName	string	`json:"firstName,omitempty"`
	LastName	string	`json:"lastName,omitempty"`
	City		string	`json:"city,omitempty"`
	PhoneNumber	string	`json:"phoneNumber,omitempty"`
	Address		string	`json:"address,omitempty"`
	BirthDate	string	`json:"birthDate,omitempty"`
}

type NewPassword struct {
	Password	string
}

// Messages

type Message struct {
	Id				int64	`json:"id,omitempty"`
	Sender			int64	`json:"sender"`
	Receiver		int64	`json:"receiver"`
	Date			string	`json:"date"`
	Content			string	`json:"content"`
}