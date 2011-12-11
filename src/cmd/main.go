package main

import (
	"fmt"
	D "Data"
	C "Core"
	LS "LoginServer"
	GS "GameServer"
	"log" 
)          
      
func main() {
	defer OnClose()
	
	log.SetFlags(log.Ltime | log.Lshortfile)
	log.SetPrefix("[Log]")
 
	D.InitializeDatabase()
	D.CreateDatabase()
	
	D.LoadData()  
   
	LS.Server = new(LS.LServer)
	GS.Server = new(GS.GServer)
	C.Start(LS.Server, "LoginServer", "127.0.0.1", 3000)
	C.Start(GS.Server, "GameServer", "127.0.0.1", 13010)    
	 
	CMD()
} 

func OnClose() {
	if x := recover(); x != nil {
		log.Println(x)
	} 
	cmd := ""
	fmt.Println("Press enter to quit...")
	fmt.Scanln(&cmd)
}
 
func CMD() {
	for {
		cmd := ""
		fmt.Scanln(&cmd)
		switch cmd {
		case "exit":
			return
		case "cleardb":
			D.ClearDatabase()
			log.Println("Database has been cleared!")
		}
	}
}
