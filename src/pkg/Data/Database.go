package Data

import (
	"encoding/hex"
	"launchpad.net/gobson/bson"
	"launchpad.net/mgo"
	"log"
)

var (
	Session  *mgo.Session
	CUsers   *mgo.Collection
	CPlayers *mgo.Collection
)

func InitializeDatabase() {
	log.Printf("Connecting to MongoDB...\n")
	session, err := mgo.Mongo("localhost")
	if err != nil {
		log.Panicf("Connecting to MongoDB has been failed! err:%v\n", err)
	}
	session.SetSyncTimeout(30 * 1000000000)
	err = session.Ping()
	if err != nil {
		log.Panicf("Connecting to MongoDB has been failed! err:%v\n", err)
	}
	log.Println("Connected!")
	Session = session
}

func CreateDatabase() {
	c := Session.DB("SGEmu").C("Users")
	p := Session.DB("SGEmu").C("Players")
	CPlayers = &p
	CUsers = &c
	index := mgo.Index{
		Key:        []string{"_id", "user", "email"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err := CUsers.EnsureIndex(index)
	if err != nil {
		panic(err)
	}

	n, _ := CUsers.Find(nil).Count()
	log.Printf("%d Users found!\n", n)

	index = mgo.Index{
		Key:        []string{"_id", "userid"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err = CPlayers.EnsureIndex(index)
	if err != nil {
		panic(err)
	}
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

func NewIID(c *mgo.Collection) uint32 {
	type dummyID struct {
		Seq uint32
	}

	d := dummyID{0}

	change := mgo.Change{Update: bson.M{"$inc": bson.M{"seq": 1}}, New: true}
	e := c.Find(bson.M{"_id": "users"}).Modify(change, &d)
	if e != nil {
		log.Panicf("Could not generate NewIID! err:%v\n", e)
	}
	return d.Seq
}

func AddAutoIncrementingField(c *mgo.Collection) {
	i, e := c.Find(bson.M{"_id": "users"}).Count()
	if e != nil {
		log.Panicf("Could not Add Auto Incrementing Field! err:%v\n", e)
	}
	if i > 0 {
		return
	}
	c.Insert(bson.M{"_id": "users", "seq": uint32(0)})
}
