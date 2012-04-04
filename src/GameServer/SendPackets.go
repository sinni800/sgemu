package GameServer

import (
	. "../Data"
	. "../SG"
)

func SendNormalChat(c *GClient, text string) {

	packet := NewPacket2(30 + len(text))
	packet.WriteHeader(CSM_CHAT)
	packet.WriteByte(0)
	packet.WriteUInt32(c.ID)
	packet.WriteString(text)
	packet.WriteColor(White)
	packet.WriteByte(0)

	Server.Run.Funcs <- func() { c.Map.Send(packet) }
}

func SendHelpChat(c *GClient, text string) {
	c.Send(HelpChatPacket(text))
}

func HelpChatPacket(text string) *SGPacket {
	packet := NewPacket2(30 + len(text))
	packet.WriteHeader(CSM_CHAT)
	packet.WriteByte(0x15)
	packet.WriteString(text)
	packet.WriteColor(HelpColor)
	return packet
}

func SendCustomChatPacket(c *GClient, text string, color Color) {
	c.Send(CustomChatPacket(text, color))
}

func CustomChatPacket(text string, color Color) *SGPacket {
	packet := NewPacket2(30 + len(text))
	packet.WriteHeader(CSM_CHAT)
	packet.WriteByte(3)
	packet.WriteString(text)
	packet.WriteColor(color)
	return packet
}

func PlayerAppear(c *GClient) *SGPacket {
	packet := NewPacket2(44)
	packet.WriteHeader(SM_PLAYER_APPEAR)
	packet.WriteByte(0)
	packet.WriteByte(c.Player.Avatar)
	packet.WriteString(c.Player.Name)
	packet.WriteString("") // guild
	packet.WriteUInt32(c.ID)
	packet.WriteUInt32(13)

	c.Log().Printf_Debug("Player[%s] Appeared at (%d,%d)", c.Player.Name, c.Player.X, c.Player.Y)

	packet.WriteInt16(c.Player.X)
	packet.WriteInt16(c.Player.Y)
	packet.WriteInt16(28)
	packet.WriteByte(0)
	packet.WriteByte(0)
	packet.WriteInt16(0)
	return packet
}

func SendMapData(client *GClient) {
	//send map info
	packet := NewPacket2(198)
	packet.WriteHeader(SM_MAP_LOAD)
	packet.WriteUInt32(client.Map.MapID)
	if client.Map.Type == BaseZone {
		packet.WriteByte(1) //force base
	} else {
		packet.WriteByte(0)
	}
	packet.WriteInt32(12)
	packet.WriteByte(0)
	packet.WriteByte(0)
	packet.WriteInt32(2068355300)
	packet.WriteInt32(2068445300)
	
	packet.WriteByte(0)
	//	 00 - non battle
	//	 01 - start/end/banned mode
	//	 02 - same as 00?
	//	 03 - spectator
	//	 04 - battle
	//	 05 - same as 03?
	//	 06 - same as 00?
	//	 07 - same as 04?
	//	 FF - same as 00?
	packet.WriteByte(0x0D)
	packet.WriteByte(0x00) //number of things on map 
	client.Send(packet)

	//17 D9 0A A0 00 00 00 00 00 00 01 FD 73 AC 33 FD 75 32 D3 04 09 00 - alien cave
	//17 D9 0A A0 00 00 00 00 00 02 01 02 6E D2 2F 02 70 58 CF 04 09 00 - alien cave with units
	//00 01 89 62 01 00 00 00 0D 00 00 FD 73 D3 12 FD 75 32 A2 00 0D 00 - main base
	//00 01 89 62 01 00 00 00 0D 00 00 FD 73 D3 05 FD 75 32 95 00 0D 00 - main base
	//
	//                                                            0D - base
	//														      09 - battle
	//														      0c - empty field
	//														   00 - non battle
	//														   01 - start/end/banned mode
	//														   02 - same as 00?
	//														   03 - spectator
	//														   04 - battle
	//														   05 - same as 03?
	//														   06 - same as 00?
	//														   07 - same as 04?
	//														   FF - same as 00?

	//Zoo Alien Cave Packet
	//packet = NewPacket2(183)
	//packet.WriteHeader(0x17)																												   //change 0x00 to 0x0a
	//packet.Write([]byte{0x17, 0xD9, 0x0A, 0xA0, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0xFD, 0x73, 0xAC, 0x33, 0xFD, 0x75, 0x32, 0xD3, 0x04, 0x09, 0x00, 0x00, 0x00, 0x15, 0xCC, 0x02, 0x00, 0x01, 0x40, 0x00, 0x34, 0x05, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x15, 0xCB, 0x07, 0x20, 0x05, 0x00, 0x00, 0x11, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x15, 0xCA, 0x0B, 0x80, 0x0C, 0x20, 0x00, 0x0A, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x15, 0xC9, 0x0B, 0x80, 0x02, 0x60, 0x00, 0x0A, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x15, 0xC8, 0x03, 0x00, 0x0C, 0x20, 0x00, 0x20, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x15, 0xC7, 0x04, 0xA0, 0x06, 0x60, 0x00, 0x00, 0x06, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x15, 0xC6, 0x0B, 0xC0, 0x09, 0x80, 0x00, 0x00, 0x06, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x15, 0xC5, 0x0D, 0x80, 0x02, 0x80, 0x00, 0x00, 0x07, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x15, 0xC4, 0x0C, 0x40, 0x01, 0xA0, 0x00, 0x00, 0x07, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x15, 0xC3, 0x05, 0x40, 0x0A, 0xC0, 0x0B, 0xB8, 0x00, 0x00, 0x00, 0x00, 0x00})
	//client.Send(packet)
	
	packet = NewPacket2(150)
	packet.WriteHeader(0x18)
	packet.WriteInt16(0x0113) //unit id
	packet.WriteByte(0x01) //unit space index
	packet.WriteByte(0)
	packet.WriteString("Shade") //unit custom name
	packet.WriteString(client.Player.Name) //player name 
	packet.Write([]byte{0x06, 0x47, 0x19, 0xF0, 0x00, 0x00, 0x30, 0x01, 0x2C, 0x01, 0x00, 0x09, 0x00, 0x00, 0x00, 0x00, 0x00, 0x58, 0x01, 0x68, 0x00, 0xC8, 0x00, 0x58, 0x01, 0x68, 0x00, 0xC8, 0x01, 0x06, 0x00, 0x00, 0x00, 0x00, 0x09, 0x00, 0x00, 0x06, 0x02, 0x01, 0x00, 0x2B, 0x00, 0x00, 0x02, 0x01, 0x4A, 0x00, 0x04, 0x00, 0x18, 0xEC, 0xBE, 0x00, 0x00, 0x00, 0x0D, 0x00, 0x00, 0x00, 0x0D, 0x01, 0x97, 0x0C, 0x13, 0x02, 0x53, 0x00, 0x00, 0x00, 0x00})
	client.Send(packet)
	//18 01 13 06 00 05 53 68 61 64 65 0A 65 76 69 6C 6D 6F 72 74 61 6C 06 47 19 F0 00 00 30 01 2C 01 00 09 00 00 00 00 00 58 01 68 00 C8 00 58 01 68 00 C8 01 06 00 00 00 00 09 00 00 06 02 01 00 2B 00 00 02 01 4A 00 04 00 18 EC BE 00 00 00 0D 00 00 00 0D 01 97 0C 13 02 53 00 00 00 00

}

func SendPlayerLeave(c *GClient) {
	packet := NewPacket2(44)
	packet.WriteHeader(SM_PLAYER_LEAVE)
	packet.WriteUInt32(c.ID)
	packet.WriteUInt32(0)
	Server.Run.Funcs <- func() { c.Map.Send(packet) }
}

func SendShopInformation(c *GClient) {
	packet := NewPacket2(1024) //change this sheet
	packet.WriteHeader(SM_SHOP_RESPONSE)

	packet.WriteByte(0x40)

	units := Shopdata.ShopUnits

	for i := 0; i < len(units); i++ {
		packet.WriteInt32(int32(i))
		if i == 0 {
			packet.WriteByte(byte(len(units)))
			packet.WriteByte(1)
		}

		//this map should be removed from here
		u, exist := Units[units[i].Name]

		if exist {
			packet.WriteByte(c.Player.Divisions[u.DType].Influence(c.Player))
			packet.WriteByte(u.Influence)
		} else {
			packet.WriteByte(0)
			packet.WriteByte(0)
		}
		packet.WriteString(units[i].Name)
		packet.WriteInt32(units[i].Money)

		packet.WriteInt32(units[i].Ore)
		packet.WriteInt32(units[i].Silicon)
		packet.WriteInt32(units[i].Uranium)
		packet.WriteByte(units[i].Sulfur)
	}

	packet.WriteUInt32(uint32(len(c.Units))) //Num of units you own
	for id, _ := range c.Units {
		packet.WriteUInt32(id)
		packet.WriteUInt64(0)
		packet.WriteUInt64(0)
		packet.WriteUInt32(0)
		packet.WriteUInt32(2)
		packet.WriteUInt64(0)
		packet.WriteUInt64(0)
	}

	packet.WriteByte(0)

	c.Send(packet)

}

func SendPlayerStats(client *GClient) {
	packet := NewPacket2(77)
	packet.WriteHeader(SM_PLAYER_STATS)
	packet.WriteUInt32(client.ID)
	packet.WriteUInt32(12)
	packet.WriteUInt32(12)
	packet.WriteByte(9)
	packet.WriteUInt32(0)

	packet.WriteInt32(client.Player.Money)
	packet.WriteInt32(client.Player.Ore)
	packet.WriteInt32(client.Player.Silicon)
	packet.WriteInt32(client.Player.Uranium)
	packet.WriteByte(client.Player.Sulfur)
	packet.WriteInt32(6)
	packet.WriteByte(client.Player.Tactics)
	packet.WriteByte(client.Player.Clout)
	packet.WriteByte(client.Player.Education)
	packet.WriteByte(client.Player.MechApt)
	packet.WriteBytes([]byte{
		0x30, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xFF, 0xFF, 0x00, 0x01, 0x00, 0x00, 0x01, 0x19, 0x00})
	//  0x00, 0x00, 0x00, 0x0C, 0x00, 0x00, 0x00, 0x0C, 0x07, 0x00, 0x00, 0x00, 0x01, 0x00, 0x04, 0x95, 0xD4, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x06, 0x0F, 0x0A, 0x0A, 0x05, 
	//  0x30, 0x00, 0x00, 0x00, 0x02, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x64, 0x00, 0x01, 0x00, 0x00, 0x01, 0x19, 0x00})
	client.Send(packet)
}

func SendUnitStats(c *GClient, unit *Unit) {
	packet := NewPacket2(100)
	packet.WriteHeader(SM_UNIT_STAT)
	packet.WriteByte(4)
	unit.WriteToPacket(packet)
	c.Send(packet)
}

func SendNewUnit(c *GClient, unit *Unit) {
	packet := NewPacket2(110)
	packet.WriteHeader(SM_UNIT_STAT)
	packet.WriteByte(2)
	unit.WriteToPacket(packet)
	c.Send(packet)
}

func SendRemoveUnit(c *GClient, unitID uint32) {
	packet := NewPacket2(20)
	packet.WriteHeader(SM_UNIT_STAT)
	packet.WriteByte(3)
	packet.WriteUInt32(unitID)
	packet.WriteByte(0)
	c.Send(packet)
}

func SendNewUnits(c *GClient, units []*Unit) {
	packet := NewPacket2(10 + len(units)*100)
	packet.WriteHeader(SM_UNIT_STAT)
	packet.WriteByte(1)
	packet.WriteByte(byte(len(units)))
	for _, unit := range units {
		unit.WriteToPacket(packet)
	}
	c.Send(packet)
}

func SendNewUnitsMap(c *GClient, units map[uint32]*Unit) {
	packet := NewPacket2(10 + len(units)*100)
	packet.WriteHeader(SM_UNIT_STAT)
	packet.WriteByte(1)
	packet.WriteByte(byte(len(units)))
	for _, unit := range units {
		unit.WriteToPacket(packet)
	}
	c.Send(packet)
}

func SendUnitInventory(c *GClient, unit *Unit) {
	packet := NewPacket2(20 + len(unit.Items)*2)
	packet.WriteHeader(SM_INVENTORY_UPDATE)
	packet.WriteUInt32(unit.ID)

	i := packet.Index
	packet.WriteByte(0)

	mi := byte(0)
	for _, item := range unit.Items {
		if item != nil {
			packet.WriteUInt16(item.ID)
			mi++
		}
	}

	packet.Buffer[i] = mi

	packet.WriteByte(0)
	c.Send(packet)
}

func SendPlayerInventory(c *GClient) {
	packet := NewPacket2(20 + len(c.Player.Items)*2)
	packet.WriteHeader(SM_INVENTORY_UPDATE)
	packet.WriteUInt32(c.ID)
	packet.WriteByte(byte(len(c.Player.Items)))

	for _, item := range c.Player.Items {
		if item != nil {
			packet.WriteUInt16(item.ID)
		}
	}
	packet.WriteByte(0)
	c.Send(packet)
}

func ProfileInfo(c *GClient, p *Player) *SGPacket {
	c.Log().Println_Debug("ProfileInfo packet")

	packet := NewPacket2(200)

	packet.WriteHeader(SM_PROFILE)
	if p != c.Player {
		packet.WriteByte(1)
		packet.WriteString(p.Name)
		packet.WriteByte(p.Avatar)
		packet.WriteInt32(0)
		packet.WriteInt16(-1)
	} else {
		packet.WriteByte(0)
	}

	packet.WriteInt16(0)

	for i := 0; i < 4; i++ {
		packet.WriteByte(p.Divisions[i].Level)
		if p.Divisions[i].Rank != "" {
			packet.WriteString(p.Divisions[i].Rank) //Custom rank
		} else {
			switch i {
			case 0:
				packet.WriteString(Ranks[p.Divisions[i].Level].Infantry)
				break
			case 1:
				packet.WriteString(Ranks[p.Divisions[i].Level].Mobile)
				break
			case 2:
				packet.WriteString(Ranks[p.Divisions[i].Level].Aviation)
				break
			case 3:
				packet.WriteString(Ranks[p.Divisions[i].Level].Organic)
				break
			}
		}
		packet.WriteByte(0x0c)
		packet.WriteByte(1)
		packet.WriteUInt32(p.Divisions[i].XP)
		packet.WriteUInt32(p.Divisions[i].TotalXP())
	}

	packet.WriteByte(p.Reincatnation)

	packet.WriteUInt32(p.Prestige)

	packet.WriteUInt32(p.Honor)
	packet.WriteUInt32(p.TotalHonor())

	packet.WriteByte(0) //Medal

	packet.WriteUInt16(p.RecordWon)
	packet.WriteUInt16(p.RecordLost)

	if p == c.Player {
		packet.WriteByte(p.Tactics)
		packet.WriteByte(p.Clout)
		packet.WriteByte(p.Education)
		packet.WriteByte(p.MechApt)
		packet.WriteUInt16(p.Points)
	} else {
		packet.WriteByte(0)
		packet.WriteInt16(0)
	}

	packet.WriteString("") //superior
	packet.WriteByte(0)

	return packet
}
