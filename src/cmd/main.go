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
	log.SetFlags(log.Ltime | log.Lshortfile)
	log.SetPrefix("[Log]")

	if !D.InitializeDatabase() {
		log.Println("Shutting down..")
		return
	}
	
	D.LoadData()  
   
	LS.Server = new(LS.LServer)
	GS.Server = new(GS.GServer)
	C.Start(LS.Server, "LoginServer", "127.0.0.1", 3000)
	C.Start(GS.Server, "GameServer", "127.0.0.1", 13010)    
	 
	CMD()
} 
 
func CMD() {
	for {
		cmd := ""
		fmt.Scanln(&cmd)
		switch cmd {
		case "exit":
			OnClose()
			return
		case "cleardb":
			D.ClearDatabase()
			log.Println("Database has been cleared!")
		}
	}
}

func OnClose() {

}
