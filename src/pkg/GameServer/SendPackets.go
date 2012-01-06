package GameServer

import (
	. "SG"
	. "Data"
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
	packet.WriteHeader(0x20)
	packet.WriteByte(0)
	packet.WriteByte(c.Player.Avatar)
	packet.WriteString(c.Player.Name)
	packet.WriteString("") // guild
	packet.WriteUInt32(c.ID)
	packet.WriteUInt32(13)

	//
	//c.Log().Printf("%d %d", c.Player.X, c.Player.Y)
	//

	packet.WriteInt16(c.Player.X)
	packet.WriteInt16(c.Player.Y)
	packet.WriteInt16(28)
	packet.WriteByte(0)
	packet.WriteByte(0)
	packet.WriteInt16(0)
	return packet
}

func SendPlayerLeave(c *GClient) {
	packet := NewPacket2(44)
	packet.WriteHeader(0x19)
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

func ProfileInfo(c *GClient, p *Player) *SGPacket {
	c.Log().Println("ProfileInfo")

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
