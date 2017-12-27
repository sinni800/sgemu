package Data

import (
	C "github.com/hjf288/sgemu/Core"
	"encoding/xml"
	//"encoding/xml"
	"fmt"
	"log"
	"os"
)

type Group byte

const (
	Weapons   = Group(4)
	Engines   = Group(1)
	Misc      = Group(7)
	Armors    = Group(5)
	Bonus     = Group(6)
	Specials  = Group(8)
	Storage   = Group(2)
	Computers = Group(3)
)

func (g Group) String() string {
	return GroupNames[g]
}

var (
	ranksPath = "./../../bin/sg_ranks.xml"
	unitsPath = "./../../bin/sg_units.xml"
	itemPath  = "./../../bin/sg_items.xml"
	shopPath  = "./../../bin/sg_shop.xml"
	bindsPath = "./../../bin/sg_binds.xml"
	Shopdata  = new(ShopData)

	Units     = make(map[string]*UnitData)
	Divisions = map[string]DType{"Infantry": Infantry,
		"Mobile":   Mobile,
		"Aviation": Aviation,
		"Organic":  Organic,
		"Other":    Other,
		"":         Other}

	Ranks        = make(map[byte]*RankData)
	Items        = make(map[uint16]*ItemData)
	ItemsByGroup = make(map[uint16][]*ItemData)
	Binds        = make(map[string]*BindingGroup)

	GroupNames = map[Group]string{Engines: "Engines", Weapons: "Weapons", Misc: "Misc", Armors: "Armors", Bonus: "Bonus", Specials: "Specials", Storage: "Storage", Computers: "Computers"}
)

type RankData struct {
	Level    byte   `xml:",attr"`
	Infantry string `xml:",attr"`
	Mobile   string `xml:",attr"`
	Aviation string `xml:",attr"`
	Organic  string `xml:",attr"`
	Unk      byte   `xml:",attr"`
}

type BindingFile struct {
	XMLName xml.Name        `xml:"Binds"`
	Groups  []*BindingGroup `xml:"BindGroup"`
}

type BindingGroup struct {
	XMLName xml.Name       `xml:"BindGroup"`
	UID     string         `xml:",attr"`
	Binds   []*BindingData `xml:"Bind"`
}

type BindingData struct {
	UID       string `xml:",attr"`
	ID        uint16 `xml:",attr"`
	GroupType Group  `xml:",attr"`
	Unk       int16  `xml:",attr"`
	Unk2      int16  `xml:",attr"`
}

type ShopData struct {
	XMLName   xml.Name    `xml:"Shop"`
	ShopUnits []*ShopUnit `xml:"Units>Unit"`
}

type ShopUnit struct {
	XMLName xml.Name `xml:"Unit"`
	Name    string
	Money   int32
	Ore     int32
	Silicon int32
	Uranium int32
	Sulfur  byte
}

type ItemDataGroup struct {
	GID      uint16      `xml:",attr"`
	ItemData []*ItemData `xml:"Item"`
}

type ItemData struct {
	Name string `xml:",attr"`
	//Description   string  //not needed
	Group         string  `xml:",attr"`
	ID            uint16  `xml:",attr"`
	GID           uint16  `xml:",attr"`
	TL            uint16  `xml:",attr"`
	Weight        uint16  `xml:",attr"`
	Space         uint16  `xml:",attr"`
	Complexity    byte    `xml:",attr"`
	EnergyUse     byte    `xml:",attr"` //also Energy-regen
	EnergyMax     uint16  `xml:",attr"`
	Damage        byte    `xml:",attr"`
	Range         float32 `xml:",attr"`
	CD            float32 `xml:",attr"`
	Effectiveness uint16  `xml:",attr"`
	Health        uint16  `xml:",attr"`
	Power         uint16  `xml:",attr"`
	Armor         uint16  `xml:",attr"`
	ItemType      byte    `xml:",attr"`
	ItemSubType   byte    `xml:",attr"`
	EnergyDrain   int8    `xml:",attr"`
	Unk1          int8    `xml:",attr"`
	EnergyType    int8    `xml:",attr"`
	Unk2          int8    `xml:",attr"`
	Unk3          int8    `xml:",attr"`
	WeaponType    byte    `xml:",attr"`
	ViewRange     float32 `xml:",attr"`
	GroupType     Group   `xml:",attr"`
	ComplexityMax uint16  `xml:",attr"`
	XpBonus       uint16  `xml:",attr"`
}

func (item *ItemData) String() string {
	return fmt.Sprintf(
		"Name:%s\tGroup:%s\tGroupType:%d\tID:%d\tGID:%d\tTL:%d\tWeight:%d\tSpace:%d\tComplexity:%d\tEnergyUse:%d\tEnergyMax:%d\tDamage:%d\tRange:%f\tCD:%f\tEffectiveness:%d\tHealth:%d\tPower:%d\tArmor:%d\tItemType:%d\tItemSubType:%d\tUnk1:%d\tUnk2:%d\tUnk3:%d\tEnergyType:%d\tEnergyDrain:%d\tWeaponType:%d\tViewRange:%f\tComplexityMax:%d\tXpBonus:%d\t",
		item.Name, item.Group, item.GroupType, item.ID, item.GID, item.TL, item.Weight, item.Space, item.Complexity, item.EnergyUse, item.EnergyMax, item.Damage, item.Range, item.CD, item.Effectiveness, item.Health, item.Power, item.Armor, item.ItemType, item.ItemSubType, item.Unk1, item.Unk2, item.Unk3, item.EnergyType, item.EnergyDrain, item.WeaponType, item.ViewRange, item.ComplexityMax, item.XpBonus)
}

func LoadData() {
	log.Println("Loading data...")

	units := make(chan bool)
	shop := make(chan bool)
	items := make(chan bool)
	binds := make(chan bool)
	ranks := make(chan bool)

	go LoadItems(items)
	go LoadUnits(units)
	go LoadRanks(ranks)
	go LoadShop(shop)
	go LoadBinds(binds)

	<-items
	log.Println("Loaded", len(Items), "Items!")
	<-units
	log.Println("Loaded", len(Units), "Units!")
	<-ranks
	log.Println("Loaded", len(Ranks), "Ranks!")
	<-shop
	log.Println("Loaded", len(Shopdata.ShopUnits), "Shop units!")
	<-binds
	log.Println("Loaded", len(Binds), "Bind groups!")
}

func LoadBinds(Done chan bool) {
	defer func() {
		if x := recover(); x != nil {
			log.Printf("%v\n%s", x, C.PanicPath())
			Done <- false
		} else {
			Done <- true
		}
	}()
	f, e := os.Open(bindsPath)
	if e != nil {
		log.Panicln(e)
	}

	defer f.Close()

	bf := &BindingFile{}

	e = xml.NewDecoder(f).Decode(bf)
	if e != nil {
		log.Panicln(e)
	}

	for _, group := range bf.Groups {
		Binds[group.UID] = group
	}
}

func LoadUnits(Done chan bool) {
	defer func() {
		if x := recover(); x != nil {
			log.Printf("%v\n%s", x, C.PanicPath())
			Done <- false
		} else {
			Done <- true
		}
	}()

	type dummyXML struct {
		XMLName       xml.Name         `xml:"Units"`
		UnitGroupData []*UnitGroupData `xml:"UnitGroup"`
	}

	f, e := os.Open(unitsPath)
	if e != nil {
		log.Panicln(e)
	}

	defer f.Close()

	dum := &dummyXML{}

	e = xml.NewDecoder(f).Decode(dum)
	if e != nil {
		log.Panicln(e)
	}

	for _, group := range dum.UnitGroupData {
		d, e := Divisions[group.Division]
		if !e {
			d = Other
		}
		for _, unit := range group.Units {
			unit.DType = d
			Units[unit.Name] = unit
		}
	}
}

func LoadRanks(Done chan bool) {
	defer func() {
		if x := recover(); x != nil {
			log.Printf("%v\n%s", x, C.PanicPath())
			Done <- false
		} else {
			Done <- true
		}
	}()

	f, e := os.Open(ranksPath)
	if e != nil {
		log.Panicln(e)
	}

	defer f.Close()

	type dummyXML struct {
		XMLName xml.Name    `xml:"Ranks"`
		Ranks   []*RankData `xml:"Rank"`
	}

	l := &dummyXML{}

	e = xml.NewDecoder(f).Decode(l)
	if e != nil {
		log.Panicln(e)
	}

	for _, rank := range l.Ranks {
		Ranks[rank.Level] = rank
	}
}

func LoadShop(Done chan bool) {
	defer func() {
		if x := recover(); x != nil {
			log.Printf("%v\n%s", x, C.PanicPath())
			Done <- false
		} else {
			Done <- true
		}
	}()
	f, e := os.Open(shopPath)
	if e != nil {
		log.Panicln(e)
	}

	defer f.Close()

	e = xml.NewDecoder(f).Decode(Shopdata)
	if e != nil {
		log.Panicln(e)
	}
}

func LoadItems(Done chan bool) {
	defer func() {
		if x := recover(); x != nil {
			log.Printf("%v\n%s", x, C.PanicPath())
			Done <- false
		} else {
			Done <- true
		}
	}()

	f, e := os.Open(itemPath)
	if e != nil {
		log.Panicln(e)
	}

	defer f.Close()

	type xmlitems struct {
		XMLName       xml.Name         `xml:"Items"`
		ItemDataGroup []*ItemDataGroup `xml:"ItemGroup"`
	}

	items := new(xmlitems)
	e = xml.NewDecoder(f).Decode(items)
	if e != nil {
		log.Panicln(e)
	}

	for _, ig := range items.ItemDataGroup {
		ItemsByGroup[ig.GID] = ig.ItemData
		for _, it := range ig.ItemData {
			Items[it.ID] = it
		}
	}
}
