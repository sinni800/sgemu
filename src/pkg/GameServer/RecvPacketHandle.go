package GameServer

import (
	C "Core"
)  


func OnWelcome(c *GClient, p *C.Packet) {
	c.Log().Println("OnWelcome Packet")
}

func OnChat(c *GClient, p *C.Packet) {
	c.Log().Println("OnChat Packet")

	p.ReadByte() //type?
	text := p.ReadString()

	//Help Text
	//packet := C.NewPacket2(30 + len(text))
	//packet.WriteHeader(CSM_CHAT)
	//packet.WriteByte(0x15)
	//packet.WriteString(fmt.Sprintf("[%s] %s", c.Player.Name, text))
	//packet.Write([]byte{0x46, 0xFA, 0xC8}) //Color

	SendNormalChat(c, text)
}

func OnPing(c *GClient, p *C.Packet) {
	packet := C.NewPacket2(20)
	packet.WriteHeader(SM_PONG)
	packet.WriteInt16(p.ReadInt16())
	packet.WriteInt16(p.ReadInt16())

	c.Send(packet)
}

func OnMove(c *GClient, p *C.Packet) {

	p.RSkip(6)
	tp := p.ReadByte() //type

	if tp == 0x16 {
		c.Log().Println("OnMove Packet")
		p.RSkip(4) //id
		c.Player.X = p.ReadInt16()
		c.Player.Y = p.ReadInt16()

		packet := C.NewPacket2(50)
		packet.WriteHeader(CSM_MOVE)
		packet.WriteInt16(0)

		packet.WriteInt16(1752)
		packet.WriteInt16(c.Player.X + c.Player.Y)

		packet.WriteInt16(1)
		packet.WriteInt16(0x0c)
		packet.WriteByte(0x1a)

		c.Log().Printf("%x %x", c.Player.X, c.Player.Y)

		packet.WriteUInt32(c.ID)
		packet.WriteInt16(c.Player.X)
		packet.WriteInt16(c.Player.Y)

		packet.WriteInt16(0)
		packet.WriteByte(0)

		Server.Run.Funcs <- func() { c.Map.Send(packet) }
	}  
} 

func OnShopRequest(c *GClient, p *C.Packet) {
	c.Log().Println("OnShopRequest Packet")
	
	p.ReadInt32()
	action := p.ReadByte()
	
	switch (action) {
		case 1:
			SendShopInformation(c);
		case 2:
			p.ReadByte() // shop unit id
			p.ReadString() // unit name
			p.ReadByte() // unkown
			c.Log().Println("Buying units is not supported yet!");
		default:
			c.Log().Println("Unkown shop action");
	}
	c.Log().Println(p)
}


func OnProfileRequest(c *GClient, p *C.Packet) {
	p.ReadByte()
	p.ReadInt32()
}