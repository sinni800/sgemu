package Data

import . "launchpad.net/gobson/bson"

type Player struct {  
	ID     string "_id"
	UserID string
	Name   string

	Faction int32

	Avatar    byte
	Tactics   byte
	Clout     byte
	Education byte
	MechApt   byte

	Map  int16 
	X, Y int16
} 

func (p *Player) Influence() byte {
	return p.Clout/2;
}


func RegisterPlayer(plyaer *Player) {
	e := CPlayers.Insert(plyaer)
	if e != nil {
		panic(e)
	}
}

func GetPlayerByUserID(id string) *Player {
	p := new(Player)
	e := CPlayers.Find(M{"userid": id}).One(p)
	if e != nil {
		panic(e)
	}
	return p
}

func SavePlayer(Player *Player) {
	e := CPlayers.Update(M{"_id": Player.ID}, M{"$set": M{"x": Player.X, "y": Player.Y}} )
	if e != nil {
		panic(e) 
	}
}