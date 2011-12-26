//Package Extractor provides function to extract game data to xmls.
package Extractor

import (
	. "Data"
	"log"
	"os"
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





func Panic() {
	if x := recover(); x != nil {
		log.Printf("Panic extractor %v\n", x)
	}
}


