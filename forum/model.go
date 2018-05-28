package forum

import "server/user"

type Info struct {
	Id		int64		`json:"id"`
	Name	string		`json:"name"`
}

type Thread struct {
	Id		int64		`json:"id,omitempty"`
	Title	string		`json:"title"`
	Forum	*Info		`json:"forum"`
	Count	int64		`json:"count"`
}

type Post struct {
	Id			int64			`json:"id,omitempty"`
	Thread		*Thread			`json:"thread"`
	User		*user.Info		`json:"user"`
	Date		string			`json:"date"`
	Content		string			`json:"content"`
}