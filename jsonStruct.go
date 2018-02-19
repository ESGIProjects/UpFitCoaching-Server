package main

type Id struct {
	Id	string
}

type Connection struct {
	Id			string
	Firstname	string
	Lastname	string
	Birthdate	string
	City		string
	Mail 		string
	Tel 		string
}

type NewPasswd struct {
	Passwd 	string
}

type Message struct {
	Id			string
	IdSender	string
	IdReceiver	string
	Body		string
	Timestamp	string
}