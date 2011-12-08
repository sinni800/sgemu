package Data

import (
	"encoding/xml"
	"os"
	"bufio"
	"fmt"
)

var (
	dataPath = "../sg_data.xml"
	shopPath = "../sg_shop.xml"
	Gamedata	= new(Data)
	Shopdata	= new(ShopData)
	Units		= make(map[string]*UnitData)
)
  
type Data struct {
	XMLName xml.Name `xml:"data"`
	Groups []*Group	`xml:"unitslist>group"`
} 
   
type Group struct { 
	XMLName xml.Name `xml:"group"`
	ID string `xml:"attr"`
	Division string `xml:"attr"`
	Name string 
	Units []*UnitData	`xml:"unitlist>unit"`
}   
  
type UnitData struct {
	UID string `xml:"attr"`
	Influence byte `xml:"attr"` 
	Space string `xml:"attr"`
	Health string `xml:"attr"`
	Armor string `xml:"attr"`
	ViewRange string `xml:"attr"`
	Speed string `xml:"attr"`
	UnitType string `xml:"attr"`
	Slots string `xml:"attr"`
	Max_Weight string `xml:"attr"`
	ViewType string `xml:"attr"`
	U1 string `xml:"attr"`
	U2 string `xml:"attr"`
	Name string
	Description string
} 
  
type ShopData struct { 
	XMLName xml.Name `xml:"Shop"` 
	ShopUnits []*ShopUnit `xml:"Units>Unit"`
}                 
  
type ShopUnit struct {
	Name string
	Money int32
	Ore int32
	Silicon int32
	Uranium int32
	Sulfur byte
}
 
func LoadData() {
	f, e := os.Open(dataPath)
	if e != nil {
		panic(e) 
	} 
	e = xml.Unmarshal(f, Gamedata)
	if e != nil {  
		panic(e) 
	} 
	
	for _,group := range Gamedata.Groups {
		for _,unit := range group.Units {
			Units[unit.Name] = unit
		} 
	}
	
	f.Close()
	
	f, e = os.Open(shopPath)
	if e != nil {
		panic(e) 
	} 
	e = xml.Unmarshal(f, Shopdata)
	if e != nil {  
		panic(e) 
	} 
	f.Close()
}

func OutputShopBinary() {
	f,e := os.Open("../shop.bin")
	
	if e != nil {  
		panic(e) 
	} 
	
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
	for i := 0;i<51;i++ {
		r.Read(buff[:7])
		s,e2 := r.ReadString(0)
		if e2 != nil {  
			panic(e) 
		} 
		fmt.Printf(format,s[:len(s)-1])
		r.Read(buff[:16])
	}
}