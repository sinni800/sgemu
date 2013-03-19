package GameServer

import (
	. "code.google.com/p/sgemu/Data"
	. "code.google.com/p/sgemu/SG"
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

func OnNameRequest(c *GClient, p *SGPacket) {
	SendPlayerNames(c)
}

func OnMove(c *GClient, p *SGPacket) {
	c.Log().Println_Debug("OnMove Packet")
	p.RSkip(6)
	tp := p.ReadByte() //type

	if tp == 0x16 {
		id := p.ReadUInt32()
		x := p.ReadInt16()
		y := p.ReadInt16()
		if id == c.ID {
			c.Player.X = x
			c.Player.Y = y
			c.Log().Printf_Debug("Player[%s] Moved (%x,%x)", c.Player.Name, x, y)
		} else {
			unit, exists := c.Units[id]
			if !exists {
				panic("moving unkown id")
			}
			unit.X = x
			unit.Y = y
			c.Log().Printf_Debug("Unit[%s] Moved (%x,%x)", unit.CustomName, x, y)
		}

		packet := NewPacket2(50)
		packet.WriteHeader(CSM_MOVE)
		packet.WriteInt16(0)

		packet.WriteUInt32(c.Map.Ticks)
		c.Map.Ticks++

		packet.WriteInt16(1)
		packet.WriteInt16(0x0c)
		packet.WriteByte(0x1a)

		packet.WriteUInt32(id)
		packet.WriteInt16(x)
		packet.WriteInt16(y)

		packet.WriteInt16(0)
		packet.WriteByte(0)

		Server.Run.Funcs <- func() { c.Map.Send(packet) }
	} else {
		c.Log().Printf_Debug("Player[%s] Move packet unkown type %d %s", c.Player.Name, tp, p)
	}
}

func OnUnitEdit(c *GClient, p *SGPacket) {
	id := p.ReadUInt32()
	unit, exist := c.Units[id]
	if exist {
		itemsRemove := p.ReadByte()

		for i := byte(0); i < itemsRemove; i++ {
			id := p.ReadUInt16()
			for i := 0; i < len(unit.Items); i++ {
				item := unit.Items[i]
				if item != nil && item.ID == id {
					unit.Items[i] = nil
					i--
					c.Player.Items[item.DBID] = item
					break
				}
			}
		}

		itemsEquip := p.ReadByte()

		for i := byte(0); i < itemsEquip; i++ {
			id := p.ReadUInt16()
			for dbid, item := range c.Player.Items {
				if item.ID == id {
					t := item.Data().GroupType
					it := unit.Items[t]
					if it == nil {
						unit.Items[t] = item
						delete(c.Player.Items, dbid)
					} else {
						c.Log().Println_Warning("Trying to move item to already equipted slot %s", c.Player.Name)
					}
					break
				}
			}
		}

		name := p.ReadString()
		if len(name) > 0 {
			unit.CustomName = name
		}

		SendPlayerInventory(c)
		SendUnitInventory(c, unit)
		SendUnitStats(c, unit)

		packet := NewPacket2(14)
		packet.WriteHeader(CSM_LAB_ENTER)
		packet.Write([]byte{0x01, 0x00, 0x00})
		c.Send(packet)

	} else {
		//c.Log().Println_Debug("Access to not existed unit %s", c.Player.Name)
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
		uid := p.ReadByte()     // shop unit id
		uname := p.ReadString() // unit name
		p.ReadByte()            // unkown
		//c.Log().Println_Debug("Buying units is not supported yet!")
		OnUnitBuyRequest(c, uid, uname)
	case 3:
		OnUnitSellRequest(c, p.ReadUInt32())
	default:
		c.Log().Println_Debug("Unkown shop action")
	}
	c.Log().Println(p)
}

func OnUnitSellRequest(c *GClient, unitID uint32) {
	if c.RemoveUnit(unitID) {
		//TODO: Fix the sell value
		c.Player.Money += 1
		SendPlayerStats(c)
	}
}

func OnUnitBuyRequest(c *GClient, unitID byte, unitName string) {
	if unitID >= byte(len(Shopdata.ShopUnits)) {
		panic("This unit does not exist")
	}

	u := Shopdata.ShopUnits[unitID]

	unit := c.AddUnit(u.Name, unitName)
	if unit != nil {
		c.Player.Money -= u.Money
		c.Player.Ore -= u.Ore
		c.Player.Silicon -= u.Silicon
		c.Player.Sulfur -= u.Sulfur

		SendPlayerStats(c)
	}
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

func OnMapChangeRequest(c *GClient, p *SGPacket) {
	mapid := p.ReadUInt32()
	x, y := p.ReadInt16(), p.ReadInt16()
	maptype := p.ReadByte()
	p.ReadByte() //dunno

	c.Log().Printf_Debug("Map id:%d [%d,%d] type:%d", mapid, x, y, maptype)
	//0 - base
	//1 - peace
	//2 - battle

	Server.Run.Funcs <- func() {
		m, exists := c.Server.Maps[mapid]

		c.Player.X = x
		c.Player.Y = y

		if !exists {
			m = NewMap(mapid, BattleZone)
			c.Server.Maps[mapid] = m
		}

		c.Map.OnLeave(c)
		c.Map = m
		c.Player.MapID = m.MapID
		c.Log().Printf_Debug("Going to map %d", m.MapID)

		m.Run.Add(func() {

			m.OnPlayerJoin(c)
			if m.Type == BattleZone {
				SendPlayerNamesBattle(c)
			} else {
				SendPlayerNames(c)
			}
			SendPlayerNames(c)
			SendMapData(c)
			if m.Type != BattleZone {
				m.OnPlayerAppear(c)
			}

			packet := NewPacket2(18)
			packet.WriteHeader(0x0E)
			packet.WriteBytes([]byte{0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
			c.Send(packet)

			packet = NewPacket2(13)
			packet.WriteHeader(0x0E)
			packet.WriteBytes([]byte{0x05, 0x00})
			c.Send(packet)
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
