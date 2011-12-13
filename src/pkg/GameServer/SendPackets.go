package GameServer

import (
	C "Core"
	. "Data"
)  
 
func SendNormalChat(c *GClient, text string) {

	packet := C.NewPacket2(30 + len(text))
	packet.WriteHeader(CSM_CHAT)
	packet.WriteByte(0)
	packet.WriteUInt32(c.ID)
	packet.WriteString(text)
	packet.WriteColor(C.Red) 
	packet.WriteByte(0)

	Server.Run.Funcs <- func() { c.Map.Send(packet) }
}  

func SendHelpChat(c *GClient, text string) {	
	c.Send(HelpChatPacket(text))
}

func HelpChatPacket(text string) (*C.Packet){
	packet := C.NewPacket2(30 + len(text))
	packet.WriteHeader(CSM_CHAT)
	packet.WriteByte(0x15)
	packet.WriteString(text)
	packet.WriteColor(C.HelpColor) 
	return packet
}

func PlayerAppear(c *GClient) *C.Packet {
	packet := C.NewPacket2(44)
	packet.WriteHeader(0x20)
	packet.WriteByte(0)
	packet.WriteByte(c.Player.Avatar)
	packet.WriteString(c.Player.Name)
	packet.WriteString("") // guild
	packet.WriteUInt32(c.ID)
	packet.WriteUInt32(11)
	 
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
	packet := C.NewPacket2(44)
	packet.WriteHeader(0x19)
	packet.WriteUInt32(c.ID)
	packet.WriteUInt32(0)
	Server.Run.Funcs <- func() { c.Map.Send(packet) }
}

func SendShopInformation(c *GClient) {
	packet := C.NewPacket2(1024) //change this sheet
	packet.WriteHeader(SM_SHOP_RESPONSE)
	
	packet.WriteByte(0x40) 
	 
	 
	units := Shopdata.ShopUnits
	  
	for i:=0;i<len(units);i++ {
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
	
	packet.WriteInt32(0) //Num of units you own
	
	/*
	packet.WriteInt32(0x29C63F) //unit id - need to send player unit list first
	packet.WriteInt32(0)
	packet.WriteInt32(0)
	packet.WriteInt32(0)
	packet.WriteInt32(0)
	packet.WriteInt32(0)
	packet.WriteInt32(0)
	packet.WriteInt32(0)
	packet.WriteInt32(0)
	packet.WriteInt32(0)
	packet.WriteInt32(0)
	*/
	
	c.Send(packet)

}

func ProfileInfo(c *GClient, p *Player) *C.Packet { 
	c.Log().Println("ProfileInfo")
	
	packet := C.NewPacket2(200)
	 
	packet.WriteHeader(SM_PROFILE)
	if (p != c.Player) {
		packet.WriteByte(1)
		packet.WriteString(p.Name)
		packet.WriteByte(p.Avatar) 
		packet.WriteInt32(0)
		packet.WriteInt16(-1)
	} else {
		packet.WriteByte(0)
	}
	
	
	packet.WriteInt16(0)
	 
	for i := 0;i<4;i++ {	
		packet.WriteByte(p.Divisions[i].Level)
		packet.WriteString(p.Divisions[i].Rank)
		packet.WriteByte(0x0c)
		packet.WriteByte(1)
		packet.WriteUInt32(p.Divisions[i].XP)
		packet.WriteUInt32(p.Divisions[i].TotalXP()) 
	}
	   
	packet.WriteByte(p.Reincatnation) //Reincatnation
	
	packet.WriteUInt32(p.Prestige) //prestige
	
	packet.WriteUInt32(p.Honor) //honor
	packet.WriteUInt32(p.TotalHonor()) //honor total
	  
	packet.WriteByte(0)  //Medal
	
	packet.WriteUInt16(p.RecordWon) //record
	packet.WriteUInt16(p.RecordLost) //record total
	 
	if (p == c.Player) {
		packet.WriteByte(p.Tactics)
		packet.WriteByte(p.Clout)
		packet.WriteByte(p.Education)
		packet.WriteByte(p.MechApt)
		packet.WriteUInt16(p.Points)   
	} else {
		packet.WriteByte(0)
		packet.WriteInt16(0)
	}
	
	packet.WriteString("Troololo")
	packet.WriteByte(0) 
	
	return packet
}