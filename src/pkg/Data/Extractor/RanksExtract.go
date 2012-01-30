package Extractor

import (
	"Data/xml"
	//"encoding/xml"
	//. "SG"
	. "encoding/binary"
	"os"
	//"fmt" 
	. "Data"
	"log"
)

func ExtractRanks(path string, outpath string, RanksExtractDone chan bool) {
	defer Panic()
	defer func() {
		RanksExtractDone <- true
	}()

	outRanks, e := os.Create(outpath + RanksOut)
	if e != nil {
		log.Panicln(e)
	}

	defer outRanks.Close()

	f, e := os.Open(RanksPath)
	if e != nil {
		log.Panicln(e)
	}

	ReadRanks(f)

	type dummyXML struct {
		XMLName       xml.Name 	   `xml:"Ranks"`
		Ranks 	      []*RankData  `xml:"Rank"`
	}
	
	l := &dummyXML{}
	l.Ranks = RanksData
	e = xml.NewEncoder(outRanks).Encode(l)
	if e != nil {
		log.Panicln(e)
	} 
}

func ReadRanks(file *os.File) {

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
	ranks := uint16(0)

	e = Read(file, BigEndian, &u)
	checkError(e, "header")
	e = Read(file, BigEndian, &u2)
	checkError(e, "header2")
	e = Read(file, BigEndian, &ranks)
	checkError(e, "ranks size")

	RanksData = make([]*RankData, ranks)

	m := make(map[byte]*RankData)

	for i := uint16(0); i < ranks; i++ {

		dtype := DType(0)
		e = Read(file, BigEndian, &dtype)
		checkError(e, "dtype")

		level := byte(0)

		e = Read(file, BigEndian, &level)
		checkError(e, "dtype")

		r, exist := m[level]
		if !exist {
			r = &RankData{Level: level}
			m[level] = r
		}
		
		e = Read(file, BigEndian, &r.Unk)
		checkError(e, "dtype")

		l := byte(0)

		e = Read(file, BigEndian, &l)
		checkError(e, "dtype")

		nameb := make([]byte, l)

		e = Read(file, BigEndian, &nameb)
		checkError(e, "rank name")

		switch dtype {
		case Infantry:
			r.Infantry = string(nameb)
		case Mobile:
			r.Mobile = string(nameb)
		case Aviation:
			r.Aviation = string(nameb)
		case Organic:
			r.Organic = string(nameb)
		}
	} 
	
	i := 0
	for _,rank := range m {
		RanksData[i] = rank
		i++
	}
}
