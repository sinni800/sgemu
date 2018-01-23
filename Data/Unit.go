package Data

import "strings"
import "code.google.com/p/sgemu/SG"

type UnitGroupData struct {
	ID       uint16      `xml:",attr"`
	Division string      `xml:",attr"`
	Name     string      `xml:",attr"`
	Units    []*UnitData `xml:"Unit"`
}

type UnitData struct {
	Name       string  `xml:",attr"`
	UID        string  `xml:",attr"`
	IID        uint16  `xml:",attr"`
	GID        uint16  `xml:",attr"`
	Influence  byte    `xml:",attr"`
	Space      uint16  `xml:",attr"`
	Health     uint16  `xml:",attr"`
	Armor      uint16  `xml:",attr"`
	ViewRange  float32 `xml:",attr"`
	Speed      float32 `xml:",attr"`
	UnitType   string  `xml:",attr"`
	Slots      uint16  `xml:",attr"`
	UnitWeight uint16  `xml:",attr"`
	Max_Weight uint16  `xml:",attr"`
	ViewType   string  `xml:",attr"`
	U1         uint16  `xml:",attr"`
	U2         uint8   `xml:",attr"`
	U3         uint32  `xml:",attr"`
	U4         uint8   `xml:",attr"`
	U5         int16   `xml:",attr"`
	U6         uint8   `xml:",attr"`
	U7         uint32  `xml:",attr"`
	U8         uint16  `xml:",attr"`
	DType      DType   `xml:",attr"`
}

type Unit struct {
	*UnitDB
	ID    uint32
	Owner *Player
	Data  *UnitData
	X, Y  int16
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

	var items [9]*Item
	bitems, exits := Binds[unit.UID]
	if exits {
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

func (unit *Unit) WriteToPacket(packet *SG.SGPacket) {
	packet.WriteUInt32(unit.ID)
	packet.WriteUInt16(unit.Data.IID) //unit id
	packet.WriteUInt16(5)
	packet.WriteUInt32(unit.XP)        //xp
	packet.WriteInt32(0)               //xp modifier
	packet.WriteUInt32(unit.TotalXP()) //xp total
	packet.WriteByte(unit.Level)
	packet.WriteByte(0)                      //hot key related
	packet.WriteByte(0)                      //hot key related
	packet.WriteUInt16(unit.HP)              //hp
	packet.WriteUInt16(unit.Data.Health)     //max hp
	packet.WriteUInt16(unit.MaxWeight())     //max-weight?
	packet.WriteUInt16(8)                    //space?
	packet.WriteUInt16(0x48)                 //weight?
	packet.WriteUInt16(8)                    //space?
	packet.WriteUInt16(unit.Data.UnitWeight) //unit-weight? unit.Data.UnitWeight
	packet.WriteUInt16(0x30)                 //speed *10
	packet.WriteUInt16(0x12c)
	packet.WriteByte(1)
	packet.WriteUInt16(unit.Data.Armor) //armor?
	packet.WriteUInt16(0)
	packet.WriteUInt16(100)
	packet.WriteUInt16(0x62)  //fire power?
	packet.WriteUInt16(0x168) //range * 2 / 10
	packet.WriteUInt16(0xc8)  //cooldown * 100
	packet.WriteUInt16(0x62)  //fire power?
	packet.WriteUInt16(0x168) //range * 2 / 10
	packet.WriteUInt16(0xc8)  //cooldown * 100
	packet.WriteUInt64(0x9000006)
	packet.WriteUInt16(unit.Kills) //kills
	packet.WriteString(unit.CustomName)
	packet.WriteString(unit.Name)
}
