package Extractor

import (
	. "Core/SG"
	"Data/xml"
	. "encoding/binary"
	"os"
	//"fmt" 
	. "Data"
	"log"
)

func ExtractItems(path string, outpath string, ItemExtractDone chan bool) {
	defer Panic()
	defer func() {
		ItemExtractDone <- true
	}()

	outItems, e := os.Create(outpath + ItemsOut)
	if e != nil {
		log.Panicln(e)
	}

	defer outItems.Close()

	f, e := os.Open(ItemsPath)
	if e != nil {
		log.Panicln(e)
	}
	ReadItems(f)

	type dummyXML struct {
		XMLName       xml.Name `xml:"Items"`
		ItemDataGroup []*ItemDataGroup
	}

	gmap := make(map[uint16]*ItemDataGroup)

	for _, item := range ItemsData {
		g, exist := gmap[item.GID]
		if exist {
			g.ItemData = append(g.ItemData, item)
		} else {
			g = &ItemDataGroup{item.GID, make([]*ItemData, 0)}
			gmap[item.GID] = g
			g.ItemData = append(g.ItemData, item)
		}
	}

	l := dummyXML{}
	l.ItemDataGroup = make([]*ItemDataGroup, len(gmap))
	for _, g := range gmap {
		l.ItemDataGroup = append(l.ItemDataGroup, g)
	}

	e = xml.Marshal(outItems, &l)
	if e != nil {
		log.Panicln(e)
	}
}

func ReadItems(file *os.File) {
	version := uint32(0)
	e := Read(file, LittleEndian, &version)
	if e != nil {
		log.Panicln("Read version panic ", e)
	}

	u := uint32(0)
	u2 := uint16(0)
	items := uint16(0)

	e = Read(file, LittleEndian, &u)
	if e != nil {
		log.Panicf("Read unk panic err:%v ", e)
	}
	e = Read(file, LittleEndian, &u2)
	if e != nil {
		log.Panicf("Read unk2 panic err:%v ", e)
	}
	e = Read(file, BigEndian, &items)
	if e != nil {
		log.Panicf("Read size panic err:%v ", e)
	}

	ItemsData = make([]*ItemData, items)

	for i := uint16(0); i < items; i++ {

		item := &ItemData{}
		ItemsData[i] = item
		toPrint := false
		size := byte(0)

		e = Read(file, LittleEndian, &size)
		if e != nil {
			log.Panicf("Read size panic iter:%d err:%v ", i, e)
		}

		nameb := make([]byte, size)

		e = Read(file, LittleEndian, &nameb)
		if e != nil {
			log.Panicf("Read name panic iter:%d err:%v ", i, e)
		}

		item.Name = string(nameb)

		e = Read(file, BigEndian, &item.ID)
		if e != nil {
			log.Panicf("Read name id iter:%d err:%v ", i, e)
		}

		e = Read(file, BigEndian, &item.GID)
		if e != nil {
			log.Panicf("Read name gid iter:%d err:%v ", i, e)
		}

		e = Read(file, LittleEndian, &item.TL)
		if e != nil {
			log.Panicf("Read name tl iter:%d err:%v ", i, e)
		}

		e = Read(file, LittleEndian, &item.Weight)
		if e != nil {
			log.Panicf("Read name weight iter:%d err:%v ", i, e)
		}

		e = Read(file, LittleEndian, &item.Space)
		if e != nil {
			log.Panicf("Read name space iter:%d err:%v ", i, e)
		}

		e = Read(file, LittleEndian, &item.Complexity)
		if e != nil {
			log.Panicf("Read name complexity iter:%d err:%v ", i, e)
		}

		e = Read(file, LittleEndian, &item.Unk3)
		if e != nil {
			log.Panicf("Read name unk1 iter:%d err:%v ", i, e)
		}

		e = Read(file, LittleEndian, &item.EnergyDrain)
		if e != nil {
			log.Panicf("Read name unk1 iter:%d err:%v ", i, e)
		}

		e = Read(file, LittleEndian, &item.Unk1)
		if e != nil {
			log.Panicf("Read name unk1 iter:%d err:%v ", i, e)
		}

		e = Read(file, LittleEndian, &item.EnergyUse)
		if e != nil {
			log.Panicf("Read name energy iter:%d err:%v ", i, e)
		}

		e = Read(file, LittleEndian, &item.GroupType)
		if e != nil {
			log.Panicf("Read name unk2 iter:%d err:%v ", i, e)
		}

		switch item.GroupType {
		case Weapons:
			e = Read(file, LittleEndian, &item.Unk3)
			if e != nil {
				log.Panicf("Read name unk3 iter:%d err:%v ", i, e)
			}

			e = Read(file, LittleEndian, &item.Damage)
			if e != nil {
				log.Panicf("Read name damage iter:%d err:%v ", i, e)
			}

			f16 := uint16(0)

			e = Read(file, BigEndian, &f16)
			if e != nil {
				log.Panicf("Read name rangeu iter:%d err:%v ", i, e)
			}
			item.Range = Float16FromBits(f16)

			e = Read(file, BigEndian, &f16)
			if e != nil {
				log.Panicf("Read name cd iter:%d err:%v ", i, e)
			}
			item.CD = Float16FromBits2(f16)

			e = Read(file, LittleEndian, &item.WeaponType)
			if e != nil {
				log.Panicf("Read name WeaponType iter:%d err:%v ", i, e)
			}

		case Engines: //Good
			e = Read(file, BigEndian, &item.Health)
			if e != nil {
				log.Panicf("Read name Health iter:%d err:%v ", i, e)
			}

			e = Read(file, BigEndian, &item.Power)
			if e != nil {
				log.Panicf("Read name Power iter:%d err:%v ", i, e)
			}

		case Misc: //Good
			e = Read(file, LittleEndian, &item.ItemSubType)
			if e != nil {
				log.Panicf("Read name ItemSubType iter:%d err:%v ", i, e)
			}

			e = Read(file, BigEndian, &item.Effectiveness)
			if e != nil {
				log.Panicf("Read name effect iter:%d err:%v ", i, e)
			}

		case Armors:
			e = Read(file, BigEndian, &item.Health)
			if e != nil {
				log.Panicf("Read name Health iter:%d err:%v ", i, e)
			}

			e = Read(file, BigEndian, &item.Armor)
			if e != nil {
				log.Panicf("Read name Armor iter:%d err:%v ", i, e)
			}

		case Bonus:
			e = Read(file, LittleEndian, &item.ItemType)
			if e != nil {
				log.Panicf("Read name ItemType iter:%d err:%v ", i, e)
			}

			e = Read(file, LittleEndian, &item.ItemSubType)
			if e != nil {
				log.Panicf("Read name ItemSubType iter:%d err:%v ", i, e)
			}

			f16 := uint16(0)
			e = Read(file, BigEndian, &f16)
			if e != nil {
				log.Panicf("Read name ViewRange iter:%d err:%v ", i, e)
			}
			item.ViewRange = Float16FromBits(f16)

		case Specials:
			e = Read(file, BigEndian, &item.ItemType)
			if e != nil {
				log.Panicf("Read name ItemType iter:%d err:%v ", i, e)
			}

			e = Read(file, BigEndian, &item.Effectiveness)
			if e != nil {
				log.Panicf("Read name effect iter:%d err:%v ", i, e)
			}

		case Storage:
			e = Read(file, BigEndian, &item.EnergyMax)
			if e != nil {
				log.Panicf("Read name EnergyMax iter:%d err:%v ", i, e)
			}

			e = Read(file, BigEndian, &item.EnergyType)
			if e != nil {
				log.Panicf("Read name effect iter:%d err:%v ", i, e)
			}

		case Computers:
			e = Read(file, BigEndian, &item.ComplexityMax)
			if e != nil {
				log.Panicf("Read name ComplexityMax iter:%d err:%v ", i, e)
			}

			e = Read(file, BigEndian, &item.XpBonus)
			if e != nil {
				log.Panicf("Read name XpBonus iter:%d err:%v ", i, e)
			}

		default:
			log.Println(item)
			log.Panicf("Unkown type:%d iter:%d err:%v ", item.GroupType, i, e)
		}

		item.Group = GroupNames[item.GroupType]

		toPrint = false //for debugging
		if toPrint {
			log.Println(item)
		}
	}
}
