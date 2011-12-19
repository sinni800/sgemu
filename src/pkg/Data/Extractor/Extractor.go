package Extractor

import (
	"Data/xml"
	. "encoding/binary"
	"os"
	//"fmt" 
	. "Core"
	. "Data"
	"bufio"
	"log"
	"strconv"
	"strings"
)

var (
	ItemsPath     = "./IINF.udf"
	ItemsDescPath = "./hlp.dat"
	NTTPath       = "./ntt.dat"

	ItemsOut = "sg_items.xml"
	BindsOut = "sg_binds.xml"

	ItemsData     []*ItemData
	BindingGroups []*BindingGroup

	ItemExtractDone = make(chan bool)
	NttExtractDone  = make(chan bool)
)

//Path: Game folder
//outpath: xmls output path 
func ReadFiles(path string, outpath string) {
	defer Panic()

	//Float16Bits
	//Float16FromBits
	/* 
		log.Printf("%x\n", Float16Bits2(3.1))
		log.Printf("%x\n", Float16Bits2(1.6))
		log.Printf("%x\n", Float16Bits(45))
		log.Printf("%x\n", Float16Bits(46))
		log.Printf("%f\n", Float16FromBits2(Float16Bits2(3.1)))
		log.Printf("%f\n", Float16FromBits2(Float16Bits2(1.6)))
		log.Printf("%f\n", Float16FromBits(Float16Bits(45)))
		log.Printf("%f\n", Float16FromBits(Float16Bits(46)))
	*/

	p, e := os.Getwd()
	if e != nil {
		log.Panicln(e)
	}

	outpath = p + outpath

	e = os.Chdir(path)
	if e != nil {
		log.Panicln(e)
	}

	go ExtractItems(path, outpath)
	go ExtractNtt(path, outpath)

	<-ItemExtractDone
	<-NttExtractDone

}

func ExtractNtt(path string, outpath string) {
	defer Panic()
	defer func() {
		NttExtractDone <- true
	}()

	nttFile, e := NewDatFile(NTTPath)
	if e != nil {
		log.Panicln(e)
	}

	defer nttFile.Close()

	offset, e := nttFile.SeekToFile("equip.txt")
	if e != nil {
		log.Panicln(e)
	}

	_ = offset

	reader := bufio.NewReader(nttFile.File)

	fullsize, e := nttFile.FileSize("equip.txt")
	if e != nil {
		log.Panicln(e)
	}
	bytes := make([]byte, fullsize)
	reader.Read(bytes)

	fulltext := string(bytes)
	groups := strings.Split(fulltext, "-1\r\n")

	BindingGroups = make([]*BindingGroup, len(groups))

	for j := 0; j < len(BindingGroups); j++ {
		split1 := strings.Split(groups[j], "\r\n")
		groupName := split1[0]

		bg := BindingGroup{}
		bg.Binds = make([]*BindingData, len(split1)-1)

		for i := 1; i < len(split1)-1; i++ {

			bind := BindingData{}
			bind.UID = groupName
			bg.UID = groupName

			split2 := strings.Split(split1[i], "\t")

			num, e := strconv.ParseUint(split2[0], 10, 32)
			if e != nil {
				log.Panicln(e)
			}
			bind.ID = uint16(num)

			found := false

			gr := strings.ToLower(split2[1])
			if gr == "engine" {
				gr = "engines"
			} else if gr == "computer" {
				gr = "computers"
			} else if gr == "weapon" {
				gr = "weapons"
			} else if gr == "ammo" {
				gr = "bonus"
			} else if gr == "special" {
				gr = "specials"
			} else if gr == "armor" {
				gr = "armors"
			}

			for gid, name := range GroupNames {
				if strings.ToLower(name) == gr {
					bind.GroupType = gid
					found = true
					break
				}
			}

			if !found {
				log.Println("Group", split2[1], "doesn't exists!")
			}

			num, e = strconv.ParseUint(split2[2], 10, 32)
			if e != nil {
				log.Panicln(e)
			}
			bind.Unk = int16(num)

			num, e = strconv.ParseUint(split2[3], 10, 32)
			if e != nil {
				log.Panicln(e)
			}
			bind.Unk2 = int16(num)

			bg.Binds[i-1] = &bind
		}

		BindingGroups[j] = &bg
	}

	outBinds, e := os.Create(outpath + BindsOut)
	if e != nil {
		log.Panicln(e)
	}

	defer outBinds.Close()

	l := BindingFile{}
	l.Groups = BindingGroups

	e = xml.Marshal(outBinds, l)
	if e != nil {
		log.Panicln(e)
	}
}

func ExtractItems(path string, outpath string) {
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
	l := ItemDataList{}
	l.Items = ItemsData

	e = xml.Marshal(outItems, l)
	if e != nil {
		log.Panicln(e)
	}
}

func Panic() {
	if x := recover(); x != nil {
		log.Printf("Panic extractor %v\n", x)
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

		e = Read(file, LittleEndian, &item.ID)
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
