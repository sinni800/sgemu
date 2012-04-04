package Data

import . "../launchpad.net/mgo/bson"
import "strconv"

//import "container/list"

type DType byte

const (
	Infantry = DType(0)
	Mobile   = DType(1)
	Aviation = DType(2)
	Organic  = DType(3)
	Other    = DType(4)
)

func (d DType) String() string {
	switch d {
		case Infantry:
			return "Infantry"
		case Mobile:
			return "Mobile"
		case Aviation:
			return "Aviation"
		case Organic:
			return "Organic"
	}
	return "Other"
}

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

	Points    uint16
	Divisions []Division

	RecordWon  uint16
	RecordLost uint16
	Prestige   uint32

	Honor uint32 //Level = Honor / 100

	Money   int32
	Ore     int32
	Silicon int32
	Uranium int32
	Sulfur  byte

	Reincatnation byte

	MapID   uint32
	X, Y 	int16

	UnitsData map[string]*UnitDB
	Items     map[string]*Item
}

type Item struct {
	DBID string "_id"
	ID   uint16
	data *ItemData
}

func (item *Item) Data() *ItemData{
	if item.data == nil {
		data,exist := Items[item.ID]
		if exist {
			item.data = data
		} else {
			panic("No such data item id " + strconv.FormatInt(int64(item.ID), 10))
		}
	} 
	return item.data
}

func CreateItem(id uint16) *Item {
	return &Item{NewID(), id, nil}
}

type Division struct {
	Type  DType
	Level byte
	Rank  string
	XP    uint32
}

func (d *Division) Influence(p *Player) byte {
	return (p.Clout / 2) + d.Level
}

func (d *Division) TotalXP() uint32 {
	level3 := uint32(d.Level) * uint32(d.Level) * uint32(d.Level)
	if d.Level >= 49 {
		return uint32((float32(1.5) * float32(level3)) * float32(d.Level-45))
	} else {
		return 6 * level3
	}
	return 0
}

func (p *Player) MaxUnits() byte {
	return 48 + (p.Tactics / 20)
}

func (d *Player) TotalHonor() uint32 {
	return 200
}

func NewPlayer() *Player {
	p := new(Player)

	p.Divisions = make([]Division, 4)
	p.Divisions[Infantry] = Division{Infantry, 1, "", 0}
	p.Divisions[Mobile] = Division{Mobile, 1, "", 0}
	p.Divisions[Aviation] = Division{Aviation, 1, "", 0}
	p.Divisions[Organic] = Division{Organic, 1, "", 0}
	
	p.MapID = 100706;

	p.UnitsData = make(map[string]*UnitDB)
	p.Items = make(map[string]*Item)

	return p
}

func (p *Player) AddUnit(name string) *UnitDB {
	u := CreateUnit(name)
	if u == nil {
		panic("No such unit Shade")
	} else {
		p.UnitsData[u.DBID] = u
	}
	return u
}

func (p *Player) SetDefaultStats() {
	p.Points = 0
	p.Money = 300000

	u := CreateUnit("Shade")
	if u == nil {
		panic("No such unit Shade")
	} else {
		p.UnitsData[u.DBID] = u
	}
	
	u = CreateUnit("Pegasus-Mk9") 
	if u == nil {
		panic("No such unit Pegasus-Mk9")
	} else {
		p.UnitsData[u.DBID] = u
	}
	
	u = CreateUnit("RCShade") 
	if u == nil {
		panic("No such unit RCShade")
	} else {
		p.UnitsData[u.DBID] = u
	}
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
		return nil
	}
	return p
}

func SavePlayer(Player *Player) {
	e := CPlayers.Update(M{"_id": Player.ID}, Player)
	if e != nil {
		panic(e)
	}
}
