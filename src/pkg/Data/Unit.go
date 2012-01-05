package Data

import "strings"

type UnitGroupData struct {
	ID       uint16   `xml:"attr"`
	Division string   `xml:"attr"`
	Name     string   `xml:"attr"`
	Units    []*UnitData `xml:"unitlist>unit"`
}

type UnitData struct {
	Name        string  `xml:"attr"`
	UID         string  `xml:"attr"`
	IID         uint16  `xml:"attr"`
	GID         uint16  `xml:"attr"`
	Influence   byte    `xml:"attr"`
	Space       uint16  `xml:"attr"`
	Health      uint16  `xml:"attr"`
	Armor       uint16  `xml:"attr"`
	ViewRange   float32 `xml:"attr"`
	Speed       float32 `xml:"attr"`
	UnitType    string  `xml:"attr"`
	Slots       uint16  `xml:"attr"`
	UnitWeight  uint16  `xml:"attr"`
	Max_Weight  uint16  `xml:"attr"`
	ViewType    string  `xml:"attr"`
	U1          uint16  `xml:"attr"`
	U2          uint8   `xml:"attr"`
	U3          uint32  `xml:"attr"`
	U4          uint8   `xml:"attr"`
	U5          int16   `xml:"attr"`
	U6          uint8   `xml:"attr"`
	U7          uint32  `xml:"attr"`
	U8          uint16  `xml:"attr"`
	DType       DType	`xml:"attr"`
}

type Unit struct {
	*UnitDB
	ID    uint32
	Owner *Player
	Data  *UnitData
}

type UnitDB struct {
	DBID       string "_id"
	Level      byte
	HP         uint16
	XP         uint32
	Kills      uint16
	CustomName string
	Name       string
	Items      []*Item
}

func (u *Unit) TotalXP() uint32 {
	n := uint32(u.Level + 1)
	return n * n
}

func CreateUnit(unitName string) *UnitDB {
	unit, e := Units[unitName]
	if !e {
		return nil
	}

	var items [8]*Item
	bitems := Binds[unit.UID]
	for _, bind := range bitems.Binds {
		pitems := ItemsByGroup[bind.ID]
		var cItem *ItemData = nil
		for _, item := range pitems {
			if item.TL == 2 && !strings.Contains(item.Name, "Gold") {
				cItem = item
				items[cItem.GroupType] = CreateItem(cItem.ID)
				break
			}
		}
	}

	return &UnitDB{NewID(), 1, unit.Health, 0, 0, unitName, unitName, items[:]}
}

//Unit Quality 
func (u *Unit) UQ() byte {
	return u.Level + (u.Owner.Clout / 2)
}

//Max Weight
func (u *Unit) MaxWeight() uint16 {
	return uint16(float32(u.Data.Max_Weight) * (1 + (float32(u.Owner.MechApt) / 120)))
}

//Alien Tech Level
func (u *Unit) ATL() byte {
	return u.Level + ((u.Owner.MechApt + u.Owner.Education) / 2)
}

//Tech Level
func (u *Unit) TL() byte {
	return u.Level + (u.Owner.Education / 2)
}
