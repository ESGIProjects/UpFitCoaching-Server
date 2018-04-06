package main

// Misc

type ErrorMessage struct {
	Message	string
}

// User management

type UserInfo struct {
	Id			int64
	UserType	int
	Mail		string
	FirstName	string
	LastName	string
	BirthDate	string
	City		string
	PhoneNumber	string
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
	Id			string
	IdSender	string
	IdReceiver	string
	Body		string
	Timestamp	string
}