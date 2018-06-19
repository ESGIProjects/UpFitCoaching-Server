package routes

import (
	"net/http"
	"server/global"
	"server/forum"
	"strconv"
)

func GetForums(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	db := global.OpenDB()
	defer db.Close()

	// Retrieve forums
	forums, err := forum.GetForums(db)
	if err != nil {
		db.Close()

		print(err.Error())
		global.SendError(w, "internal_error", http.StatusInternalServerError)
		return
	}

	global.SendJSON(w, forums, http.StatusOK)
}

func GetThreads(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	db := global.OpenDB()
	defer db.Close()

	// Get forumId from request
	forumId, err := strconv.Atoi(r.URL.Query().Get("forumId"))
	if err != nil {
		db.Close()

		print(err.Error())
		global.SendError(w, "internal_error", http.StatusInternalServerError)
		return
	}

	// Get forum from DB
	forumInfo := forum.Info{}
	err = forum.GetForumFromId(db, &forumInfo, int64(forumId))
	if err != nil {
		db.Close()

		print(err.Error())
		global.SendError(w, "internal_error", http.StatusInternalServerError)
		return
	}

	// Get threads from DB
	threads, err := forum.GetThreadsFromForum(db, forumInfo)
	if err != nil {
		db.Close()

		print(err.Error())
		global.SendError(w, "internal_error", http.StatusInternalServerError)
		return
	}

	global.SendJSON(w, threads, http.StatusOK)
}

func GetThread(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	db := global.OpenDB()
	defer db.Close()

	// Get threadId from request
	threadId, err := strconv.Atoi(r.URL.Query().Get("threadId"))
	if err != nil {
		db.Close()

		print(err.Error())
		global.SendError(w, "internal_error", http.StatusInternalServerError)
		return
	}

	// Get thread from DB
	thread := forum.Thread{}
	err = forum.GetThreadFromId(db, &thread, int64(threadId))
	if err != nil {
		db.Close()

		print(err.Error())
		global.SendError(w, "internal_error", http.StatusInternalServerError)
		return
	}

	// Get posts from DB
	posts, err := forum.GetPostsFromThread(db, thread)
	if err != nil {
		db.Close()

		print(err.Error())
		global.SendError(w, "internal_error", http.StatusInternalServerError)
		return
	}

	global.SendJSON(w, posts, http.StatusOK)
}

func CreateThread(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	db := global.OpenDB()
	defer db.Close()

	// Get fields from request
	userId, _ := strconv.Atoi(r.PostFormValue("userId"))
	forumId, _ := strconv.Atoi(r.PostFormValue("forumId"))
	title := r.PostFormValue("title")
	date := r.PostFormValue("date")
	content := r.PostFormValue("content")

	// Inserting thread into DB
	res, err := db.Exec("INSERT INTO threads (title, forumId) VALUES (?, ?)", title, forumId)
	if err != nil {
		db.Close()

		print(err.Error())
		global.SendError(w, "internal_error", http.StatusInternalServerError)
		return
	}

	// Get the new thread ID
	threadId, err := res.LastInsertId()
	if err != nil {
		db.Close()

		print(err.Error())
		global.SendError(w, "internal_error", http.StatusInternalServerError)
		return
	}

	// Inserting post into DB
	res, err = db.Exec("INSERT INTO posts (threadId, userId, date, content) VALUES (?, ?, ?, ?)", threadId, userId, date, content)
	if err != nil {
		db.Close()

		print(err.Error())
		global.SendError(w, "internal_error", http.StatusInternalServerError)
		return
	}

	// Get the new post ID
	postId, err := res.LastInsertId()
	if err != nil {
		db.Close()

		print(err.Error())
		global.SendError(w, "internal_error", http.StatusInternalServerError)
		return
	}

	json := make(map[string]int64)
	json["postId"] = postId
	json["threadId"] = threadId

	global.SendJSON(w, json, http.StatusCreated)
}

func AddPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	db := global.OpenDB()
	defer db.Close()

	// Get fields from request
	userId, _ := strconv.Atoi(r.PostFormValue("userId"))
	threadId, _ := strconv.Atoi(r.PostFormValue("threadId"))
	date := r.PostFormValue("date")
	content := r.PostFormValue("content")

	// Inserting post into DB
	res, err := db.Exec("INSERT INTO posts (threadId, userId, date, content) VALUES (?, ?, ?, ?)", threadId, userId, date, content)
	if err != nil {
		db.Close()

		print(err.Error())
		global.SendError(w, "internal_error", http.StatusInternalServerError)
		return
	}

	// Get the new post ID
	postId, err := res.LastInsertId()
	if err != nil {
		db.Close()

		print(err.Error())
		global.SendError(w, "internal_error", http.StatusInternalServerError)
		return
	}

	json := make(map[string]int64)
	json["postId"] = postId

	global.SendJSON(w, json, http.StatusCreated)
}