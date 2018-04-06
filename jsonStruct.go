package main

// Misc

type ErrorMessage struct {
	message	string
}

// User management

type UserInfo struct {
	id			int64
	userType	int
	mail		string
	firstName	string
	lastName	string
	birthDate	string
	city		string
	phoneNumber	string
}

type CoachInfo struct {
	id			int64
	mail		string
	firstName	string
	lastName	string
	address		string
	city		string
	phoneNumber	string
}

type NewPassword struct {
	password	string
}

// Messages

type Message struct {
	Id			string
	IdSender	string
	IdReceiver	string
	Body		string
	Timestamp	string
}