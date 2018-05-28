package routes

import (
	"net/http"
	"server/global"
	"server/forum"
	"strconv"
)

func GetForums(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
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