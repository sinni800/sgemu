package GameServer

import (
	"Core"
	"net"
)

type GamePacketFunc func(c *GClient, p *Core.Packet)

var (
	Handler map[int]GamePacketFunc
	Server  *GServer
)

type GServer struct {
	Core.Server
	Maps map[int]*Map
	IDG  *Core.IDGen
	Run  *Core.Runner
	DBRun  *Core.Runner
}

func (serv *GServer) OnSetup() {
	serv.Server.OnSetup()
	serv.Maps = make(map[int]*Map)
	serv.IDG = Core.NewIDG()
	serv.Maps[0] = NewMap()
	serv.Run = Core.NewRunner()
	serv.Run.Start()
	serv.DBRun = Core.NewRunner()
	serv.DBRun.Start()
}

func init() {
	Handler = make(map[int]GamePacketFunc)
	Handler[CSM_CHAT] = OnChat
	Handler[CM_PING] = OnPing
	Handler[CSM_MOVE] = OnMove
	Handler[CM_PROFILE] = OnProfileRequest 
	Handler[CM_SHOP_REQUEST] = OnShopRequest
} 

func (serv *GServer) OnConnect(socket *net.TCPConn) {
	serv.Log.Printf("Client connected to GServer!")
	client := new(GClient)
	client.Server = serv
	Core.SetupClient(client, socket, serv)
}
