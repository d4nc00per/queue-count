package main

import (
	"time"

	"gopkg.in/mgo.v2"
)

// GetDbClient will start a new session and return a client
func GetDbClient() *DbClient {
	sess, err := mgo.Dial("localhost:27017")

	if err != nil {
		panic(err)
	}

	db := sess.DB("test")

	return &DbClient{
		session: sess,
		Db:      db,
	}
}

// DbClient holds the database session
type DbClient struct {
	session *mgo.Session
	Db      *mgo.Database
}

//Close will close the connection
func (that *DbClient) Close() {
	that.session.Close()
}

// LogMsg holds a log message
type LogMsg struct {
	Occured time.Time
	Message string
}

//Log will log a message
func (that *DbClient) Log(msg string) {
	that.Db.C("log").Insert(&LogMsg{Occured: time.Now(), Message: msg})
}

//Queues will return the queues collection
func (that *DbClient) Queues() *mgo.Collection {
	return that.Db.C("queues")
}
