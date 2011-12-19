package Extractor

import (
	"os"
	. "encoding/binary"
	"strings"
	"errors"
)

type DatFile struct {
	Files map[string]uint32
	File *os.File
}

func NewDatFile(path string) (dat *DatFile,err error){
	dat = &DatFile{}
	
	file, e := os.Open(path)
	if e != nil {
		return nil,e
	}
	
	size := uint32(0)
	e = Read(file, LittleEndian, &size)
	if e != nil {
		return nil,err
	}
	
	dat.Files = make(map[string]uint32, size)
	dat.File = file
	 
	nameb := make([]byte, 13)
	 
	for i:=uint32(0);i<size;i++ {
		offset := uint32(0)
		e = Read(file, LittleEndian, &offset)
		if e != nil {
			return nil,err
		}
		
		e = Read(file, LittleEndian, &nameb)
		if e != nil {
			return nil,err
		}
		
		j := 0
		for ;j<len(nameb);j++ {
			if nameb[j] == 0x00 {
				break
			}
		}
		
		name := string(nameb[:j])
		
		dat.Files[strings.ToLower(name)] = offset
	}
	 
	return dat,nil
}

func (dat *DatFile) Close() {
	dat.File.Close()
}

func (dat *DatFile) SeekToFile(file string) (coffset int64, e error) {
	file = strings.ToLower(file)
	offset, exists := dat.Files[file] 
	if exists {
		coffset, e = dat.File.Seek(int64(offset), 0)
		return coffset, e
	} 
	return -1,errors.New("cant find specified file")
}

func (dat *DatFile) FileSize(file string) (int64,error) {
	file = strings.ToLower(file)
	offset, exists := dat.Files[file] 
	if exists {
		found := uint32(0)
		for t,off := range dat.Files {
			if t != file {
				if off >= offset {
					if found == 0 || off < found {
						found = off
					}
				}
			}
		}
		if found == 0 {
			target,e := dat.File.Seek(0, 2)
			if e != nil {
				return -1,e
			}
			
			return target-int64(offset),nil
		}
		
		return int64(found)-int64(offset),nil
	} 
	return -1,errors.New("cant find specified file")
}