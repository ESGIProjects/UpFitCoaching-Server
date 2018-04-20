package main

// Misc

type ErrorMessage struct {
	Message	string	`json:"message"`
}

// User management

type UserInfo struct {
	Id			int64	`json:"id"`
	UserType	int		`json:"type,omitempty"`
	Mail		string	`json:"mail,omitempty"`
	FirstName	string	`json:"firstName,omitempty"`
	LastName	string	`json:"lastName,omitempty"`
	BirthDate	string	`json:"birthDate,omitempty"`
	City		string	`json:"city,omitempty"`
	PhoneNumber	string	`json:"phoneNumber,omitempty"`
}

type CoachInfo struct {
	Id			int64
	Mail		string
	FirstName	string
	LastName	string
	Address		string
	City		string
	PhoneNumber	string
}

type NewPassword struct {
	Password	string
}

// Messages

type Message struct {
	FromUserId		int		`json:"fromId"`
	FromUserType	int		`json:"fromType"`
	ToUserId		int		`json:"toId"`
	ToUserType		int		`json:"toType"`
	Date			string	`json:"date"`
	Content			string	`json:"content"`
}