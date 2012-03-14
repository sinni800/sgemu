package main

import (
	C "Core"
	D "Data"
	LS "LoginServer" 
	"fmt" 
	"log"
	"os"
	"runtime" 
	"SG"
	//"WebAdmin"
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
	
	//Note: do this in the game server?
	D.LoadData() 
	
	SG.ReadConfig()
	
	LS.Server = new(LS.LServer)
	LS.Server.Start("LoginServer", SG.Config.LSConfig.IP, SG.Config.LSConfig.Port, SG.Config.LSConfig.WANIP)
	 
	
	go ListenSignals() 
	
	CMD()   
}

func ListenSignals() {
	/*
	for sig := range signal.Incoming {
		if os.SIGINT == sig || os.SIGTERM == sig || os.SIGKILL == sig || os.SIGQUIT == sig {
			OnClose()
			return
		}
	}
	*/
}  

func OnClose() {
	if Closing {
		return  
	}
	Closing = true

	if x := recover(); x != nil {
		log.Printf("%v %s\n", x, C.PanicPath())
	}

	defer func() {
		if x := recover(); x != nil {
			log.Printf("%v %s\n", x, C.PanicPath())
		}
		cmd := ""
		log.Println("Press enter to quit...")
		fmt.Scanln(&cmd)
		os.Exit(0)
	}() 

	if LS.Server != nil {
		//Do stuff with LS
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
			if Closing {
				return
			}
		}
	}
}
