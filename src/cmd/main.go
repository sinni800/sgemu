package main

import (
	"fmt"
	D "Data"
	C "Core"
	LS "LoginServer"
	GS "GameServer"
	"log" 
	"runtime"
	"os/signal"
	"os"
)       
  
var (
	Closing = false
)      
      
func main() { 
	defer OnClose()
	        
	runtime.GOMAXPROCS(5)
	  
	log.SetFlags(log.Ltime | log.Lshortfile)
	log.SetPrefix("[Log]") 
	 
	D.InitializeDatabase()
	D.CreateDatabase()
	 
	D.LoadData()  
    
	LS.Server = new(LS.LServer)
	GS.Server = new(GS.GServer)
	C.Start(LS.Server, "LoginServer", "127.0.0.1", 3000)
	C.Start(GS.Server, "GameServer", "127.0.0.1", 13010)    
	  
  	go ListenSignals() 
	   
	CMD()  
}     
  
func ListenSignals() {
	for signal := range signal.Incoming {
			_ = signal
			OnClose() 
			return
	}  
} 
   
func OnClose() {
	if Closing { return }; Closing = true
	
	defer func() {
		if x := recover(); x != nil {
			log.Printf("%v\n", x)
		}   
		cmd := ""
		fmt.Println("Press enter to quit...")
		fmt.Scanln(&cmd)
		os.Exit(0)
	}()   
	
	if GS.Server != nil {
		GS.Server.OnShutdown()
	} 
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
		default:
			if Closing { return }
		}
	}
}
