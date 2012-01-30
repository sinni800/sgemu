package LoginServer

import (
	"Core"
	. "SG"
	"net"
	"strconv"
)

type LoginPacketFunc func(c *LClient, p *SGPacket)

var (
	Handler map[int]LoginPacketFunc
	Server  *LServer
)

type LServer struct {
	Core.CoreServer
	WANAddr *net.TCPAddr
	WANIP   string
}

func (serv *LServer) Start(name, ip string, port int, wanip string) (err error) {
	err = Core.Start(serv,name,ip,port)
	if err != nil { 
		return err
	}
	
	serv.WANIP = wanip
	serv.WANAddr, err = net.ResolveTCPAddr("tcp", serv.WANIP+":"+strconv.Itoa(port))
	if err != nil {
		serv.Log.Printf("Server start failed %s", err.Error())
		return err
	}
	
	go startRPC()
	
	return nil
}
  
func (serv *LServer) OnSetup() {
	serv.CoreServer.OnSetup()
}

func init() {
	Handler = make(map[int]LoginPacketFunc)
	Handler[CSM_WELCOME] = OnWelcome
	Handler[CM_WELCOME2] = OnWelcome2
	Handler[CSM_REGISTER] = OnRegister
	Handler[CM_LOGIN] = OnLogin
	Handler[CM_LWELCOME] = OnLoginWelcome
	Handler[CSM_FACTION_DATA] = OnFactionDataRequest
	Handler[CM_PLANET_DATA] = OnPlanetDataRequest
	Handler[CM_REGISTER_DONE] = OnRegisterDone
	Handler[CS_FRIEND_SELECT] = OnFriendSelect
}

func (serv *LServer) OnConnect(socket *net.TCPConn) {
	serv.Log.Printf("Client connected to LoginServer!")
	client := new(LClient)
	client.Server = serv
	Core.SetupClient(client, socket, serv)
}
