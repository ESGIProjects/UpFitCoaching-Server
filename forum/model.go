package forum

import (
	"server/user"
	"database/sql"
	"errors"
)

type Info struct {
	Id		int64		`json:"id"`
	Name	string		`json:"name"`
}

type Thread struct {
	Id			int64		`json:"id"`
	Title		string		`json:"title"`
	Forum		Info		`json:"forum"`
	LastUpdated	string		`json:"lastUpdated,omitempty"`
	LastUser	*user.Info	`json:"lastUser,omitempty"`
}

type Post struct {
	Id			int64			`json:"id"`
	Thread		Thread			`json:"thread"`
	User		user.Info		`json:"user"`
	Date		string			`json:"date"`
	Content		string			`json:"content"`
}

func GetForumFromId(db *sql.DB, forumInfo *Info, id int64) (error) {
	query := "SELECT * FROM forums WHERE id = ?"

	row := db.QueryRow(query, id).Scan(&forumInfo.Id, &forumInfo.Name)
	if row == sql.ErrNoRows {
		return errors.New("no_forum")
	}

	return nil
}

func GetForums(db *sql.DB) ([]Info, error) {
	query := `SELECT * FROM forums`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	forums := make([]Info, 0)

	for rows.Next() {
		forumInfo := Info{}
		rows.Scan(&forumInfo.Id, &forumInfo.Name)

		forums = append(forums, forumInfo)
	}

	return forums, nil
}

func GetThreadFromId(db *sql.DB, thread *Thread, id int64) (error) {
	query := "SELECT * FROM threads WHERE id = ?"
	var forumId int64

	row := db.QueryRow(query, id).Scan(&thread.Id, &thread.Title, &forumId)
	if row == sql.ErrNoRows {
		return errors.New("no_thread")
	}

	forumInfo := Info{}
	err := GetForumFromId(db, &forumInfo, forumId)
	if err != nil {
		return err
	}

	thread.Forum = forumInfo

	// Get last updated
	var lastUpdated string
	var lastUserId int64
	row = db.QueryRow("SELECT userId, date FROM posts WHERE threadId = ? ORDER BY date DESC LIMIT 0,1", thread.Id).Scan(&lastUserId, &lastUpdated)
	if row == sql.ErrNoRows {
		return errors.New("no_last_updated")
	}

	thread.LastUpdated = lastUpdated

	// Get last user
	lastUser := user.Info{}
	_, err = user.GetFromId(db, &lastUser, lastUserId)
	if err != nil {
		return errors.New("no_last_user")
	}

	thread.LastUser = &lastUser

	return nil
}


func GetThreadsFromForum(db *sql.DB, forumInfo Info) ([]Thread, error) {
	query := `SELECT id, title FROM threads WHERE forumId = ?`

	rows, err := db.Query(query, forumInfo.Id)
	if err != nil {
		return nil, err
	}

	threads := make([]Thread, 0)

	for rows.Next() {
		thread := Thread{}
		rows.Scan(&thread.Id, &thread.Title)
		thread.Forum = forumInfo

		// Get last updated
		var lastUpdated string
		var lastUserId int64
		row := db.QueryRow("SELECT userId, date FROM posts WHERE threadId = ? ORDER BY date DESC LIMIT 0,1", thread.Id).Scan(&lastUserId, &lastUpdated)
		if row != sql.ErrNoRows {
			thread.LastUpdated = lastUpdated
		}

		// Get last user
		lastUser := user.Info{}
		_, err = user.GetFromId(db, &lastUser, lastUserId)
		if err == nil {
			thread.LastUser = &lastUser
		}

		threads = append(threads, thread)
	}

	return threads, nil
}

func GetPostsFromThread(db *sql.DB, thread Thread) ([]Post, error) {
	query := "SELECT id, userId, date, content FROM posts WHERE threadId = ?"

	rows, err := db.Query(query, thread.Id)
	if err != nil {
		return nil, err
	}

	posts := make([]Post, 0)

	for rows.Next() {
		post := Post{}
		var userId int64

		rows.Scan(&post.Id, &userId, &post.Date, &post.Content)

		userInfo := user.Info{}
		_, err = user.GetFromId(db, &userInfo, userId)
		if err != nil {
			continue
		}

		post.User = userInfo
		post.Thread = thread

		posts = append(posts, post)
	}

	return posts, nil
}