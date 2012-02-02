package GameServer

import (
	. "SG"
)

func OnWelcome(c *GClient, p *SGPacket) {
	c.Log().Println_Debug("OnWelcome Packet")
}

func OnDisconnectPacket(c *GClient, p *SGPacket) {
	c.Log().Println_Debug("OnDisconnect Packet")
	c.OnDisconnect()
}

func OnChat(c *GClient, p *SGPacket) {
	p.ReadByte() //type?
	text := p.ReadString()

	c.Log().Println_Debug("OnChat Packet ", text)

	SendNormalChat(c, text)
}

func OnPing(c *GClient, p *SGPacket) {
	packet := NewPacket2(20)
	packet.WriteHeader(SM_PONG)
	packet.WriteInt16(p.ReadInt16())
	packet.WriteInt16(p.ReadInt16())

	c.Send(packet)
}

func OnLabraryEnter(c *GClient, p *SGPacket) {
	
	packet := NewPacket2(20)
	packet.WriteHeader(CSM_LAB_ENTER)
	packet.WriteByte(0)
	packet.WriteInt32(p.ReadInt32())
	packet.WriteInt16(0)
	packet.WriteByte(0)

	c.Send(packet)
}

func OnGameEnter(c *GClient, p *SGPacket) {
	typ := p.ReadByte()
	switch typ {
	case 1:
		packet := NewPacket2(21)
		packet.WriteHeader(CSM_GAME_ENTER)
		packet.WriteByte(typ)
		packet.WriteInt32(p.ReadInt32())
		packet.WriteInt32(0)
		packet.WriteByte(0)
		c.Send(packet)
	case 2:
		packet := NewPacket2(17)
		packet.WriteHeader(CSM_GAME_ENTER)
		packet.WriteByte(4)
		packet.WriteByte(1)
		packet.WriteInt32(0)
		c.Send(packet)
	}
}

func OnMove(c *GClient, p *SGPacket) {
	c.Log().Println_Debug("OnMove Packet")
	p.RSkip(6)
	tp := p.ReadByte() //type

	if tp == 0x16 {
		p.RSkip(4) //id
		c.Player.X = p.ReadInt16()
		c.Player.Y = p.ReadInt16()

		packet := NewPacket2(50)
		packet.WriteHeader(CSM_MOVE)
		packet.WriteInt16(0)

		packet.WriteInt16(1752)
		packet.WriteInt16(c.Player.X + c.Player.Y)

		packet.WriteInt16(1)
		packet.WriteInt16(0x0c)
		packet.WriteByte(0x1a)

		c.Log().Printf_Debug("Player[%s] Moved (%x,%x)", c.Player.Name, c.Player.X, c.Player.Y)

		packet.WriteUInt32(c.ID)
		packet.WriteInt16(c.Player.X)
		packet.WriteInt16(c.Player.Y)

		packet.WriteInt16(0)
		packet.WriteByte(0)

		Server.Run.Funcs <- func() { c.Map.Send(packet) }
	}
}

func OnShopRequest(c *GClient, p *SGPacket) {
	c.Log().Println_Debug("OnShopRequest Packet")

	p.ReadInt32()
	action := p.ReadByte()

	switch action {
	case 1:
		SendShopInformation(c)
	case 2:
		p.ReadByte()   // shop unit id
		p.ReadString() // unit name
		p.ReadByte()   // unkown
		c.Log().Println_Debug("Buying units is not supported yet!")
	default:
		c.Log().Println_Debug("Unkown shop action")
	}
	c.Log().Println(p)
}

func OnProfileRequest(c *GClient, p *SGPacket) {
	p.ReadByte()
	id := p.ReadUInt32()
	if id == c.ID {
		c.Send(ProfileInfo(c, c.Player))
	} else {
		c.Map.Run.Add(func() {
			player, exists := c.Map.Players[id]
			if exists {
				c.Send(ProfileInfo(c, player.Player))
			}
		})
	}
}

func OnProfileLeave(c *GClient, p *SGPacket) {
	t := p.ReadByte()
	cl := p.ReadByte()
	e := p.ReadByte()
	m := p.ReadByte()
	totalp := uint16(t + cl + e + m)
	if totalp <= c.Player.Points {
		c.Player.Points -= totalp
		c.Player.Tactics += t
		c.Player.Clout += cl
		c.Player.Education += e
		c.Player.MechApt += m
	}
}
