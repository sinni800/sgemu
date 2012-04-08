//Package Extractor provides function to extract game data to xmls.
package Extractor

import (
	. "Data"
	"bufio"
	"log"
	"os"
)

var (
	ItemsPath     = "./IINF.udf"
	UnitsPath     = "./UNF.bkm"
	ItemsDescPath = "./hlp.dat"
	NTTPath       = "./ntt.dat"
	RanksPath     = "./RNF.udf"

	HelperPath = "../addon.xml"

	ItemsOut = "sg_items.xml"
	BindsOut = "sg_binds.xml"
	UnitsOut = "sg_units.xml"
	RanksOut = "sg_ranks.xml"

	ItemsData     []*ItemData
	BindingGroups []*BindingGroup
	UnitGroups    []*UnitGroupData
	RanksData     []*RankData
	fileHelper 	  *os.File
)

//Path: Game folder.
//
//outpath: xmls output path.
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
	
	fileHelper, e = os.Open(HelperPath)
	if e != nil {
		log.Panicln(e)
	}

	outpath = p + outpath

	e = os.Chdir(path)
	if e != nil {
		log.Panicln(e)
	}

	ItemExtractDone := make(chan bool)
	NttExtractDone := make(chan bool)
	UnitsExtractDone := make(chan bool)
	RanksExtractDone := make(chan bool)
	
	go ExtractItems(path, outpath, ItemExtractDone)
	go ExtractNtt(path, outpath, NttExtractDone)
	go ExtractUnits(path, outpath, UnitsExtractDone)
	go ExtractRanks(path, outpath, RanksExtractDone)
	
	<-ItemExtractDone
	<-NttExtractDone
	<-UnitsExtractDone
	<-RanksExtractDone


	fileHelper.Close();
}

func Panic() {
	if x := recover(); x != nil {
		log.Printf("Panic extractor %v\n", x)
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
		panic(e)
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
		log.Printf(format, s[:len(s)-1])
		r.Read(buff[:16])
	}
}
