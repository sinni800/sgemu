package GameServer

import "net/rpc"
import "net"
import "Data"
import "SG"
import "strconv"
//import "errors"
import "time"

type AcceptorRPC struct {
}

func (ar *AcceptorRPC) Queue(a *Data.InStruct,b *int) error {
	Data.LoginQueue.Add(a.IP, a.ID)
	return nil
}


func Accept(l net.Listener) (error) {
	for {
		c,err := l.Accept() 
		if err != nil {
			return err
		}
		go HandleAuth(c)
	}
	return nil
} 

func HandleAuth(client net.Conn) {
	client.SetReadDeadline(time.Now().Add(time.Second*10))
		
	b := make([]byte, len(SG.RPCKey))
	n, err := client.Read(b)
	if err != nil {
		client.Close()
		return
	}
	
	if n != len(b) {
		client.Write([]byte{0}) 
		client.Close()
	}
	
	for i:=0;i<len(b);i++ {
		if b[i] != SG.RPCKey[i] {
			client.Write([]byte{0}) 
			client.Close()
		}
	}
	
	client.SetReadDeadline(time.Time{}) 
	
	n, err = client.Write([]byte{1}) 
	Server.Log.Println("RPC Client Connected!")
	
	Server.RPCServer.ServeConn(client) 
}

func startRPCServer() {
	defer func() {
		if x := recover(); x != nil {
			Server.Log.Println(x)
		} 
	}()
	
	Server.RPCServer = rpc.NewServer()
	serv := Server.RPCServer
	 
	l, e := net.Listen("tcp", SG.Config.RPCConfig.IP + ":" + strconv.Itoa(SG.Config.RPCConfig.Port))
	if e != nil {
		panic("listen error:" + e.Error())
	}
	 
	acceptor := &AcceptorRPC{}
	serv.Register(acceptor)
	Server.RPCListener = l
	//Server.RPCListener.Accept()
	//serv.ServeCodec()
	go Accept(Server.RPCListener)	
}

