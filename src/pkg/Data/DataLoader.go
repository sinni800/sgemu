package Data

import (
	"Data/xml"
	"bufio"
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

var (
	dataPath = "../sg_data.xml"
	itemPath = "../sg_items.xml"
	shopPath = "../sg_shop.xml"
	bindsPath = "../sg_binds.xml"
	Gamedata = new(Data)
	Shopdata = new(ShopData)

	Units     = make(map[string]*UnitData)
	Divisions = map[string]DType{"Infantry": Infantry,
		"Mobile":   Mobile,
		"Aviation": Aviation,
		"Organic":  Organic,
		"":         Other}

	Ranks = make(map[byte]*RankData)
	Items = make(map[uint16]*ItemData)
	Binds = make(map[string]*BindingGroup)

	GroupNames = map[Group]string{Engines: "Engines", Weapons: "Weapons", Misc: "Misc", Armors: "Armors", Bonus: "Bonus", Specials: "Specials", Storage: "Storage", Computers: "Computers"}
)

type Data struct {
	XMLName xml.Name     `xml:"data"`
	Groups  []*GroupData `xml:"unitslist>group"`
	Ranks   []*RankData  `xml:"rankslist>rank"`
}

type RankData struct {
	Level    byte   `xml:"attr"`
	Infantry string `xml:"attr"`
	Mobile   string `xml:"attr"`
	Aviation string `xml:"attr"`
	Organic  string `xml:"attr"`
}

type BindingFile struct {
	XMLName xml.Name        `xml:"BindingFile"`
	Groups  []*BindingGroup `xml:">BindingGroup"`
}
 
type BindingGroup struct {
	XMLName xml.Name       `xml:"BindingGroup"`
	UID     string         `xml:"attr"`
	Binds   []*BindingData
}

type BindingData struct {
	XMLName   xml.Name `xml:"Bind"`
	UID       string   `xml:"attr"`
	ID        uint16   `xml:"attr"`
	GroupType Group    `xml:"attr"`
	Unk       int16    `xml:"attr"`
	Unk2      int16    `xml:"attr"`
}

type GroupData struct {
	XMLName  xml.Name `xml:"group"`
	ID       string   `xml:"attr"`
	Division string   `xml:"attr"`
	Name     string
	Units    []*UnitData `xml:"unitlist>unit"`
}

type UnitData struct {
	UID         string `xml:"attr"`
	Influence   byte   `xml:"attr"`
	Space       string `xml:"attr"`
	Health      string `xml:"attr"`
	Armor       string `xml:"attr"`
	ViewRange   string `xml:"attr"`
	Speed       string `xml:"attr"`
	UnitType    string `xml:"attr"`
	Slots       string `xml:"attr"`
	Max_Weight  uint16 `xml:"attr"`
	ViewType    string `xml:"attr"`
	U1          string `xml:"attr"`
	U2          string `xml:"attr"`
	Name        string
	Description string
	DType       DType
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

type ItemDataList struct {
	XMLName xml.Name    `xml:"ItemList"`
	Items   []*ItemData `xml:"Items>Item"`
}

type ItemData struct {
	XMLName       xml.Name `xml:"Item"`
	Name          string
	Description   string
	Group         string `xml:"attr"`
	ID            uint16 `xml:"attr"`
	GID           uint16 `xml:"attr"`
	TL            uint16
	Weight        uint16
	Space         uint16
	Complexity    byte
	EnergyUse     byte //also Energy-regen
	EnergyMax     uint16
	Damage        byte
	Range         float32
	CD            float32
	Effectiveness uint16
	Health        uint16
	Power         uint16
	Armor         uint16
	ItemType      byte
	ItemSubType   byte
	EnergyDrain   int8
	Unk1          int8
	EnergyType    int8
	Unk2          int8
	Unk3          int8
	WeaponType    byte
	ViewRange     float32
	GroupType     Group `xml:"attr"`
	ComplexityMax uint16
	XpBonus       uint16
}

func (item *ItemData) String() string {
	return fmt.Sprintf(
		"Name:%s\tDescription:%s\tGroup:%s\tGroupType:%d\tID:%d\tGID:%d\tTL:%d\tWeight:%d\tSpace:%d\tComplexity:%d\tEnergyUse:%d\tEnergyMax:%d\tDamage:%d\tRange:%f\tCD:%f\tEffectiveness:%d\tHealth:%d\tPower:%d\tArmor:%d\tItemType:%d\tItemSubType:%d\tUnk1:%d\tUnk2:%d\tUnk3:%d\tEnergyType:%d\tEnergyDrain:%d\tWeaponType:%d\tViewRange:%f\tComplexityMax:%d\tXpBonus:%d\t",
		item.Name, item.Description, item.Group, item.GroupType, item.ID, item.GID, item.TL, item.Weight, item.Space, item.Complexity, item.EnergyUse, item.EnergyMax, item.Damage, item.Range, item.CD, item.Effectiveness, item.Health, item.Power, item.Armor, item.ItemType, item.ItemSubType, item.Unk1, item.Unk2, item.Unk3, item.EnergyType, item.EnergyDrain, item.WeaponType, item.ViewRange, item.ComplexityMax, item.XpBonus)

}

func LoadData() {
	log.Println("Loading data...")

	units := make(chan bool)
	shop := make(chan bool)
	items := make(chan bool)
	binds := make(chan bool)

	go LoadItems(items)
	go LoadUnitsAndRanks(units)
	go LoadShop(shop)
	go LoadBinds(binds)

	<-items 
	log.Println("Loaded", len(Items), "Items!")
	<-units
	log.Println("Loaded", len(Units), "Units!")
	log.Println("Loaded", len(Ranks), "Ranks!")
	<-shop
	log.Println("Loaded", len(Shopdata.ShopUnits), "Shop units!")
	<-binds
	log.Println("Loaded", len(Binds), "Bind groups!")

}

func LoadBinds(Done chan bool) {
	defer func() {
		if x := recover(); x != nil {
			log.Printf("%v\n", x)
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
 
	bf := BindingFile{}

	e = xml.Unmarshal(f, &bf)
	if e != nil {
		log.Panicln(e)
	}
	
	for _,group := range bf.Groups {
		Binds[group.UID] = group
	}
}

func LoadUnitsAndRanks(Done chan bool) {
	defer func() {
		if x := recover(); x != nil {
			log.Printf("%v\n", x)
			Done <- false
		} else {
			Done <- true
		}
	}()
	f, e := os.Open(dataPath)
	if e != nil {
		log.Panicln(e)
	}

	defer f.Close()

	e = xml.Unmarshal(f, Gamedata)
	if e != nil {
		log.Panicln(e)
	}

	for _, group := range Gamedata.Groups {
		d, e := Divisions[group.Division]
		if !e {
			d = Other
		}
		for _, unit := range group.Units {
			unit.DType = d
			Units[unit.Name] = unit
		}
	}

	for _, rank := range Gamedata.Ranks {
		Ranks[rank.Level] = rank
	}
}

func LoadShop(Done chan bool) {
	defer func() {
		if x := recover(); x != nil {
			log.Printf("%v\n", x)
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

	e = xml.Unmarshal(f, Shopdata)
	if e != nil {
		log.Panicln(e)
	}
}

func LoadItems(Done chan bool) {
	defer func() {
		if x := recover(); x != nil {
			log.Printf("%v\n", x)
			Done <- false
		} else {
			Done <- true
		}
	}()
	ItemsDataList := new(ItemDataList)
	f, e := os.Open(itemPath)
	if e != nil {
		log.Panicln(e)
	}

	defer f.Close()

	e = xml.Unmarshal(f, ItemsDataList)
	if e != nil {
		log.Panicln(e)
	}

	for _, item := range ItemsDataList.Items {
		Items[item.ID] = item
	}
}

func OutputShopBinary() {
	defer func() {
		if x := recover(); x != nil {
			log.Printf("%v\n", x)
		}
	}()
	f, e := os.Open("../shop.bin")

	if e != nil {
		log.Panicln(e)
	}

	defer f.Close()

	format := `<Unit>
	<Name>%s</Name>
	<Money>1</Money>
	<Ore>0</Ore>
	<Silicon>0</Silicon>
	<Uranium>0</Uranium>
	<Sulfur>0</Sulfur>
</Unit> 
`

	r := bufio.NewReader(f)
	bytes := [20]byte{}
	buff := bytes[:]
	for i := 0; i < 51; i++ {
		r.Read(buff[:7])
		s, e2 := r.ReadString(0)
		if e2 != nil {
			log.Panicln(e)
		}
		fmt.Printf(format, s[:len(s)-1])
		r.Read(buff[:16])
	}
}
