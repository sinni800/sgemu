package GameServer

import (
	C "Core"
	D "Data"
	. "SG"
	//R "reflect"
)

type GClient struct {
	C.CoreClient
	Key    byte
	packet *SGPacket

	Disconnecting bool
	ID            uint32
	Player        *D.Player
	Units         map[uint32]*D.Unit
	Server        *GServer
	Map           *Map
}

func (client *GClient) StartRecive() {
	defer client.OnDisconnect()
	callback := func(p *SGPacket) { client.ParsePacket(p) }

	for {
		err := client.packet.ReadPacketFromStream(client.Socket, callback)
		if err != 0{
			return
		}
	}
}

func (client *GClient) OnConnect() {
	defer client.OnDisconnect()
	
	client.Disconnecting = false
	userID, q := D.LoginQueue.Check(client.IP)
	if !q {
		return
	}

	id, r := client.Server.IDG.Next()

	if !r {
		return
	}

	client.Log().Println_Debug("ID " + userID)
	client.Player = D.GetPlayerByUserID(userID)

	client.Units = make(map[uint32]*D.Unit)

	for _, unitdb := range client.Player.UnitsData {
		id, r := client.Server.IDG.Next()
		if !r {
			client.Log().Println_Warning("No more ids left - server is full!")
			return
		}
		name, e := D.Units[unitdb.Name]
		if !e {
			client.Log().Println_Warning("Unit name does not exists")
			continue
		}
		client.Units[id] = &D.Unit{unitdb, id, client.Player, name}
	}

	client.Log().Println("name " + client.Player.Name)

	client.packet = NewPacket()
	client.packet.Index = 0
	client.ID = id

	
	Server.Run.Funcs <- func() {
		client.Server.Maps[0].Run.Add(func() {
			client.Server.Maps[0].OnPlayerJoin(client)
			client.SendWelcome()
		})
	}
	client.StartRecive()
}

func (client *GClient) OnDisconnect() {	
	if x := recover(); x != nil {
		client.Log().Printf("panic : %v \n %s", x, C.PanicPath())
	}
	
	if client.Disconnecting {
		return
	}
	client.Disconnecting = true
	
	client.Socket.Close()
	
	if client.Map != nil {
		client.Map.Run.Add(
		func() {
			client.Map.OnLeave(client)
		})
	}
	if client.Units != nil {
		for id, _ := range client.Units {
			client.Server.IDG.Return(id)
		}
	}
	if client.Player != nil {
		client.Server.IDG.Return(client.ID)
		client.Server.DBRun.Funcs <- func() { D.SavePlayer(client.Player) }
	}
	client.MainServer.Server().Log.Println("Client Disconnected! %s", client.Socket.RemoteAddr())
}

func (client *GClient) Send(p *SGPacket) {
	if !p.Encrypted {
		op := p.Buffer[3]
		if op > 13 || (op > 0 && op < 3) || (op > 3 && op < 11) {
			p.WSkip(2)
			EncryptPacket(p.Buffer[:p.Index], client.Key)
			p.Encrypted = true
			client.Key++
		}
		p.WriteLen()
	}
	client.Socket.Write(p.Buffer[:p.Index])
}

func (client *GClient) SendRaw(p *SGPacket) {
	p.WriteLen()
	client.Socket.Write(p.Buffer[:p.Index])
}

func (client *GClient) SendWelcome() {

	packet := NewPacket2(13)
	packet.WriteHeader(0x02)
	packet.Write([]byte{0x00, 0x00})
	//client.Send(packet)
	
 
	packet = NewPacket2(14)
	packet.WriteHeader(0x5A)
	packet.Write([]byte{0x1E, 0x06, 0x00})
	//client.Send(packet)
	
	packet = NewPacket2(16)
	packet.WriteHeader(0x66)
	packet.Write([]byte{0x02, 0x01, 0x02, 0x00, 0x02})
	//client.Send(packet)
 
	packet = NewPacket2(20)
	packet.WriteHeader(0x7D)
	packet.Write([]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
	//client.Send(packet)
	

	SendPlayerStats(client)

	packet = NewPacket2(10 + len(client.Units)*100)
	packet.WriteHeader(SM_UNIT_STAT)
	packet.WriteByte(1)
	packet.WriteByte(byte(len(client.Units)))	
	for _, unit := range client.Units {
		unit.WriteToPacket(packet)
	}
	client.Send(packet)

	for _, unit := range client.Units {
		SendUnitInventory(client, unit)
	}

	SendPlayerInventory(client)

	//send map info
	packet = NewPacket2(198)
	packet.WriteHeader(0x17)
	packet.WriteBytes([]byte{0x00, 0x01, 0x87, 0x0A, 0x01, 0x00, 0x00, 0x00, 0x0C, 0x00, 0x00, 0x7B, 0x48, 0x98, 0xE4, 0x7B, 0x49, 0xF8, 0x74, 0x00, 0x0D, 0x00}) //, 0x00, 0x00, 0x0D, 0x42, 0x01, 0xC0, 0x0C, 0x40, 0x00, 0x00, 0x06, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0D, 0x41, 0x0C, 0x00, 0x0B, 0x40, 0x00, 0x00, 0x06, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0D, 0x40, 0x04, 0x00, 0x08, 0x80, 0x00, 0x00, 0x06, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0D, 0x3F, 0x0A, 0x60, 0x07, 0xA0, 0x00, 0x00, 0x07, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0D, 0x3E, 0x03, 0xA0, 0x02, 0x80, 0x00, 0x00, 0x06, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0D, 0x3D, 0x0B, 0xE0, 0x02, 0x40, 0x00, 0x00, 0x06, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0D, 0x3C, 0x07, 0x80, 0x01, 0x60, 0x1D, 0x4C, 0x00, 0x00, 0x00, 0x00, 0x0D, 0x00, 0x00, 0x0D, 0x3B, 0x06, 0x20, 0x04, 0x40, 0x1D, 0x4C, 0x00, 0x00, 0x00, 0x00, 0x0D, 0x00, 0x00, 0x0D, 0x3A, 0x07, 0x40, 0x08, 0x80, 0x1D, 0x4C, 0x00, 0x00, 0x00, 0x00, 0x0D, 0x00, 0x00, 0x0D, 0x39, 0x09, 0x40, 0x0C, 0x00, 0x1D, 0x4C, 0x00, 0x00, 0x00, 0x00, 0x0D, 0x00, 0x00, 0x0D, 0x38, 0x07, 0xC0, 0x0E, 0xA0, 0x1D, 0x4C, 0x00, 0x00, 0x00, 0x00, 0x0D})
	client.Send(packet)

	//send player name 
	packet = NewPacket2(28 + len(client.Map.Players)*13)
	packet.WriteHeader(SM_PLAYER_NAME)
	packet.WriteInt16(int16(len(client.Map.Players)))
	for _, s := range client.Map.Players {
		packet.WriteString(s.Player.Name)
		packet.WSkip(2)
	}
	client.Send(packet)

	//send spawn palyer
	client.Map.OnPlayerAppear(client)

	packet = NewPacket2(18)
	packet.WriteHeader(0x0E)
	packet.WriteBytes([]byte{0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
	client.Send(packet)

	packet = NewPacket2(13)
	packet.WriteHeader(0x0E)
	packet.WriteBytes([]byte{0x05, 0x00})
	client.Send(packet)

	packet = NewPacket2(13)
	packet.WriteHeader(0x3E)
	packet.WriteBytes([]byte{0x00, 0x00})
	//client.Map.Send(packet)

	//SendCustomChatPacket(client, "***Merry Christmas***!", Red)
	//SendCustomChatPacket(client, "***Merry Christmas***!", Green) 
}

func (client *GClient) Log() *C.Logger {
	return Server.Log
}

func (client *GClient) ParsePacket(p *SGPacket) {
	header := p.ReadByte()
	
	fnc, exist := Handler[int(header)]
	
	if !exist {
		client.Log().Printf_Warning("isnt registred : %s", p)
		return
	}
	//client.Log().Printf("Header(%d) len(%d) : % #X\n %s", header, len(p.Buffer), p.Buffer, p.Buffer)
	//client.Log().Printf("Handle %s\n", R.TypeOf(fnc))
	fnc(client, p)	
}
