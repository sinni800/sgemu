package Data

import (
	"launchpad.net/gobson/bson"
	"launchpad.net/mgo"
	"log"
	"encoding/hex"
) 
  
var (
	Session  *mgo.Session
	CUsers   *mgo.Collection
	CPlayers *mgo.Collection
)

func InitializeDatabase() bool {
	log.Printf("Connecting to MongoDB...\n")
	session, err := mgo.Mongo("localhost")
	if err != nil {
		log.Println(err)
		return false
	}
	session.SetSyncTimeout(30 * 1000000000)
	err = session.Ping()
	if err != nil {
		log.Println(err)
		return false
	}
	log.Println("Connected!")
	Session = session

	c := session.DB("SGEmu").C("Users")
	p := session.DB("SGEmu").C("Players")
	CPlayers = &p
	CUsers = &c
	index := mgo.Index{
		Key:        []string{"id", "user", "email"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err = CUsers.EnsureIndex(index)
	if err != nil {
		log.Println(err)
		return false
	}

	n, _ := CUsers.Find(nil).Count()
	log.Printf("%d Users found!\n", n)

	index = mgo.Index{
		Key:        []string{"id", "userid"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err = CPlayers.EnsureIndex(index)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}


func ClearDatabase() {
	if Session != nil {
		CUsers.RemoveAll(bson.M{})
		CPlayers.RemoveAll(bson.M{})
	}
}

func NewID() string {
	return hex.EncodeToString([]byte(string(bson.NewObjectId())))
}
