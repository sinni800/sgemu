package LoginServer

import (
	. "SG"
	GS "GameServer"
)

func SendMessage(c *LClient, errcode int, msg string) {
	packet := NewPacket2(len(msg) + 20)
	packet.WriteHeader(CSM_REGISTER)
	switch errcode {
	case 0:
		break
	case 1:
		packet.WriteByte(0)
		break
	default:
		packet.WriteByte(4)
		packet.WriteString(msg)
		break
	}
	packet.WSkip(2)
	c.Send(packet)
}

func SendToGameServer(c *LClient, username string) {
	packet := NewPacket2(20)
	packet.WriteHeader(SM_SENDIP)
	packet.Index--
	ip := []byte(GS.Server.Addr.IP.To4())
	packet.Write([]byte{ip[3], ip[2], ip[1], ip[0]})
	packet.WriteUInt16(uint16(GS.Server.Port))
	packet.WriteByte(0x0c)
	packet.WriteByte(1)
	packet.WriteString(username)

	c.Send(packet)
}

func SendWelcome(c *LClient) {
	packet := NewPacket2(40)
	packet.WriteHeader(0x7E)
	packet.Buffer[4] = 0x1B
	packet.WriteRawString("SERVER CONNECTED\n")
	c.SendRaw(packet)

	c.Log().Println("Welcome packet sent!")
}
