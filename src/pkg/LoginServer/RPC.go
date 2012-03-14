package LoginServer

import "net/rpc"
import "Data"
import "SG"
import "strconv"
import "net"

var RPCClient *rpc.Client

func startRPC() {
	defer func() {
		if x := recover(); x != nil {
			Server.Log.Println(x)
		} 
	}()
 
 
 	var client net.Conn
 
	Server.Log.Println_Info("Trying to connect to the RPC Server...")
	for {
		var err error
		client, err = net.Dial("tcp", SG.Config.RPCConfig.WANIP + ":" + strconv.Itoa(SG.Config.RPCConfig.Port))
		if err == nil { 
			break
		} 
	}
	
	_, e := client.Write(SG.RPCKey)
	if e != nil { 
		panic(e) 
	} 
	var Ok [1]byte 
	_,e = client.Read(Ok[:1])
	if e != nil { 
		panic(e) 
	}
	
	if Ok[0] == 0 {
		Server.Log.Println_Warning("Wrong RPC Key!")
	} else {	
		Server.Log.Println_Info("Conneted to RPC Server!")
	}
	
	RPCClient = rpc.NewClient(client)
	//TODO: handle disconnect
}

func addClient(ip, id string) {
	var b int
	RPCClient.Go("AcceptorRPC.Queue", Data.InStruct{IP:ip, ID:id}, &b, nil)
}