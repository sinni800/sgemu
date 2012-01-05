package main

import (
	. "Data/Extractor"
	"flag"
	"fmt"
)

//Args: use like extractor.exe "game path"
func main() {
	flag.Parse()
	args := flag.Args()
  
	path := "./"

	if len(args) != 0 {
		path = args[0]
	} else {
		fmt.Println("No game path, Use like extractor.exe \"game path\"")
	} 

	//path = "C:/Games/Shattered Galaxy/" //Default path used for debugging

	ReadFiles(path, "./")
 
	cmd := ""
	fmt.Println("Done! Press enter to quit...")
	fmt.Scanln(&cmd)
} 
