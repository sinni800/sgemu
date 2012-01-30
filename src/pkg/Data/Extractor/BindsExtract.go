package Extractor

import (
	. "Data"
	"Data/xml"
	//"encoding/xml"
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

func ExtractNtt(path string, outpath string, NttExtractDone chan bool) {
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

	BindingGroups = make([]*BindingGroup, len(groups)-1)

	for j := 0; j < len(BindingGroups); j++ {
		split1 := strings.Split(groups[j], "\r\n")
		groupName := split1[0]

		bg := BindingGroup{}
		bg.Binds = make([]*BindingData, len(split1)-1)

		for i := 1; i < len(bg.Binds); i++ {

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

	e = xml.NewEncoder(outBinds).Encode(l)
	if e != nil {
		log.Panicln(e)
	}
}
