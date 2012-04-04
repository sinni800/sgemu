package Extractor

import (
	"encoding/xml"
	//"encoding/xml"
	. "../../SG"
	. "encoding/binary"
	"os"
	//"fmt" 
	. "../../Data"
	"log"
)

func ExtractUnits(path string, outpath string, UnitExtractDone chan bool) {
	defer Panic()
	defer func() {
		UnitExtractDone <- true
	}()
	 
	outUnits, e := os.Create(outpath + UnitsOut)
	if e != nil {
		log.Panicln(e)
	}

	defer outUnits.Close()

	f, e := os.Open(UnitsPath)
	if e != nil {
		log.Panicln(e)
	}
	ReadUnits(f)
	

	
	defer outUnits.Close()
	
	ReadUnitsHelper(fileHelper)
	
	for _,group := range UnitGroups {
		switch {
			//case len(group.Units) == 1:
			//	group.Division = "Other"
			case group.ID > 300:
				group.Division = "Organic"
			case group.ID > 200:
				group.Division = "Aviation"
			case group.ID > 100:	
				group.Division = "Mobile"
			case group.ID > 0:		
				group.Division = "Infantry"	
		} 
		for _,unit := range group.Units {
			unit.DType = Divisions[group.Division]
		}
	}
	
	type dummyXML struct {
		XMLName       xml.Name `xml:"Units"`
		UnitGroupData []*UnitGroupData `xml:"UnitGroup"`
	}

	l := &dummyXML{}
	l.UnitGroupData = UnitGroups

	b,e := xml.MarshalIndent(l,"","\t")
	if e != nil {
		log.Panicln(e)
	}
	_, e = outUnits.Write(b)
	if e != nil {
		log.Panicln(e)
	}

}

func ReadUnitsHelper(file *os.File) {
	checkError := func(e error, text string) {
		if e != nil {
			log.Panicf("Read panic %s err:%v ", text, e)
		} 
	}
	_ = checkError
	
	type UnitHelper struct {
		Name        string  `xml:"name,attr"`
		UID         string  `xml:"uid,attr"`
		Influence   byte    `xml:"influence,attr"`
		Slots       uint16  `xml:"slots,attr"`
		UnitWeight  uint16  `xml:"weight,attr"`
	}
	
	type dummyXML struct {
		XMLName       xml.Name `xml:"data"`
		Units 	  []*UnitHelper `xml:"units-list>division>unit"`
	}
	
	l := &dummyXML{}
	e := xml.NewDecoder(file).Decode(l)
	checkError(e, "unmarshal")
	
	for _,unitHelper := range l.Units { 
		for _,group := range UnitGroups {
			for _,unit := range group.Units {
				if unit.Name == unitHelper.Name {
					unit.UID = unitHelper.UID
					unit.Influence = unitHelper.Influence
					unit.Slots = unitHelper.Slots
					unit.UnitWeight = unitHelper.UnitWeight
					goto End
				}
			}
		}
		log.Print("Could not find unit:", unitHelper.Name)
		End:
	}
}

func ReadUnits(file *os.File) {
	
	checkError := func(e error, text string) {
		if e != nil {
			log.Panicf("Read panic %s err:%v ", text, e)
		}
	}
	
	version := uint32(0)
	e := Read(file, LittleEndian, &version)
	checkError(e, "version")
 
	u := uint32(0)
	u2 := uint16(0)
	groups := uint16(0)

	e = Read(file, BigEndian, &u)
	checkError(e, "header")
	e = Read(file, BigEndian, &u2)
	checkError(e, "header2")
	e = Read(file, BigEndian, &groups)
	checkError(e, "units size")

	UnitGroups = make([]*UnitGroupData,  groups)

	for i := uint16(0); i < groups; i++ {
	
		g := &UnitGroupData{}
		g.Units = make([]*UnitData, 0)
		
		e = Read(file, BigEndian, &g.ID)
		checkError(e, "gid")
		 
		l := byte(0)
		e = Read(file, BigEndian, &l)
		checkError(e, "group name length")
		
		nameb := make([]byte, l)

		e = Read(file, BigEndian, &nameb)
		checkError(e, "group name")
		
		g.Name = string(nameb)
		
		UnitGroups[i] = g
	}
	
	units := uint16(0)
	e = Read(file, BigEndian, &units)
	checkError(e, "units size")
	
	for i := uint16(0); i < units; i++ {
		 
		unit := &UnitData{} 
	
		e = Read(file, BigEndian, &unit.IID) 
		checkError(e, "units IID")
		
		e = Read(file, BigEndian, &unit.GID)
		checkError(e, "units GID")
		
		l := byte(0)
		e = Read(file, BigEndian, &l)
		checkError(e, "unit name length")
		
		nameb := make([]byte, l)

		e = Read(file, LittleEndian, &nameb)
		checkError(e, "unit name")
		
		unit.Name = string(nameb)
		
		e = Read(file, BigEndian, &unit.Max_Weight)
		checkError(e, "unit Max_Weight")
		
		e = Read(file, BigEndian, &unit.Space)
		checkError(e, "unit Space")
		
		e = Read(file, BigEndian, &unit.Armor)
		checkError(e, "unit armor")
		
		e = Read(file, BigEndian, &unit.U8)
		checkError(e, "unit u1")
		
		e = Read(file, BigEndian, &unit.U1)
		checkError(e, "unit u1")
		
		e = Read(file, BigEndian, &unit.U2)
		checkError(e, "unit u2")
		
		e = Read(file, BigEndian, &unit.U3)
		checkError(e, "unit u3")
		
		e = Read(file, BigEndian, &unit.U4)
		checkError(e, "unit u4")
		
		e = Read(file, LittleEndian, &unit.U5)
		checkError(e, "unit u5")
		 
		e = Read(file, BigEndian, &unit.Health)
		checkError(e, "unit Health") 
		
		//not sure if needed
		//unit.Health -= 10
		 
		speed := uint16(0)
		e = Read(file, BigEndian, &speed)
		checkError(e, "unit u2")
		
		unit.Speed = Float16FromBits3(speed)
		
		viewrange := uint16(0)
		e = Read(file, BigEndian, &viewrange)
		checkError(e, "unit u2")
		
		unit.ViewRange = Float16FromBits4(viewrange)
		
		e = Read(file, BigEndian, &unit.U6)
		checkError(e, "unit U6")
		
		e = Read(file, BigEndian, &unit.U7)
		checkError(e, "unit U7")
		
		for _,g := range UnitGroups {
			if g.ID == unit.GID {
				g.Units = append(g.Units, unit)
			}
		}
	}
}