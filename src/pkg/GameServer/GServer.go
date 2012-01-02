package GameServer

import (
	"Core"
	. "SG"
	"Data"
	"net"
)

type GamePacketFunc func(c *GClient, p *SGPacket)

var (
	Handler map[int]GamePacketFunc
	Server  *GServer
)

type GServer struct {
	Core.Server
	Maps  map[int]*Map
	IDG   *Core.IDGen
	Run   *Core.Runner
	DBRun *Core.Runner
	Sdr   *Core.Scheduler
}

func (serv *GServer) OnSetup() {
	serv.Server.OnSetup()
	serv.Maps = make(map[int]*Map)
	serv.IDG = Core.NewIDG()
	serv.Maps[0] = NewMap()
	serv.Run = Core.NewRunner()
	serv.DBRun = Core.NewRunner()
	serv.Sdr = Core.NewScheduler()
	serv.Run.Start()
	serv.DBRun.Start()
	serv.Sdr.Start()

	serv.Sdr.AddMin(func() { serv.SavePlayers() }, 1)
}

func init() {
	Handler = make(map[int]GamePacketFunc)
	Handler[CSM_CHAT] = OnChat
	Handler[CM_PING] = OnPing
	Handler[CSM_MOVE] = OnMove
	Handler[CM_PROFILE] = OnProfileRequest
	Handler[CM_LEAVE_PROFILE] = OnProfileLeave
	Handler[CM_SHOP_REQUEST] = OnShopRequest
	Handler[CSM_GAME_ENTER] = OnGameEnter
}

func (serv *GServer) SavePlayers() {
	serv.Log.Printf("Saving Players...")
	for _, m := range serv.Maps {
		m.Run.Add(func() {
			for _, c := range m.Players {
				serv.DBRun.Add(func() { Data.SavePlayer(c.Player) })
			}
		})
	}
	serv.Sdr.AddMin(func() { serv.SavePlayers() }, 1)
}

func (serv *GServer) OnShutdown() {
	serv.Run.StopAndWait()
	serv.Log.Printf("GServer runner stopped!")
	serv.DBRun.StopAndWait()
	serv.Log.Printf("GServer DB runner stopped!")
	for id, m := range serv.Maps {
		m.Run.StopAndWait()
		serv.Log.Printf("Mapid %d runner stopped!", id)
		for _, c := range m.Players {
			Data.SavePlayer(c.Player)
		}
	}
	serv.Server.Socket.Close()
	serv.Log.Printf("GServer socket closed!")
}

func (serv *GServer) OnConnect(socket *net.TCPConn) {
	serv.Log.Printf("Client connected to GServer!")
	client := new(GClient)
	client.Server = serv
	Core.SetupClient(client, socket, serv)
}
