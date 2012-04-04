package GameServer

import C "../Core"
import . "../SG"

type MapType byte

const (
	BaseZone = MapType(0)
	PeaceZone = MapType(1)
	BattleZone = MapType(2)
)


type Map struct {
	MapID	uint32
	Type	MapType
	Players map[uint32]*GClient
	Run     *C.Runner
}

func NewMap(mapid uint32, typ MapType) *Map {
	m := new(Map)
	m.Players = make(map[uint32]*GClient)
	m.Run = C.NewRunner()
	m.Run.Start()
	m.MapID = mapid
	m.Type = typ
	return m
}

func (m *Map) OnPlayerJoin(c *GClient) {
	c.Map = m
	m.Players[c.ID] = c
}

func (m *Map) OnPlayerAppear(c *GClient) {
	for _, value := range m.Players {
		c.Send(PlayerAppear(value))
	}

	m.SendAllExcept(PlayerAppear(c), c)
}

func (m *Map) OnLeave(c *GClient) {
	SendPlayerLeave(c)
	delete(m.Players, c.ID)
}

func (m *Map) Send(p *SGPacket) {
	m.SendAllExcept(p, nil)
}

func (m *Map) SendAllExcept(p *SGPacket, c *GClient) {
	for _, value := range m.Players {
		if c != value {
			value.Send(p)
		}
	}
}
