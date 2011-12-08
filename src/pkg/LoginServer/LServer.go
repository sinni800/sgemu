package LoginServer

import (
	"Core"
	"net"
)

type LoginPacketFunc func(c *LClient, p *Core.Packet)


var (
	Handler map[int]LoginPacketFunc
	Server  *LServer
)

type LServer struct {
	Core.Server
}

func (serv *LServer) OnSetup() {
	serv.Server.OnSetup()
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
