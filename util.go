package main

import (
	"database/sql"
	"time"
	"math/rand"
)

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := "root"
	dbName := "upfit"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

func randomPasswd(strlen int) string {
	var random = rand.New(rand.NewSource(time.Now().UnixNano()))
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ@&-_0123456789"
	passwd := make([]byte, strlen)
	for i := range passwd {
		passwd[i] = chars[random.Intn(len(chars))]
	}
	return string(passwd)
}
