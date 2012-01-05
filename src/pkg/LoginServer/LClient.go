package LoginServer

import (
	C "Core"
	. "SG"
	D "Data"
	"log" 
)

type LClient struct {
	C.Client
	Key      byte
	packet   *SGPacket
	TempUser *D.User
	Server   *LServer
}

func (client *LClient) StartRecive() {
	defer client.OnDisconnect()
	for {
		bl, err := client.Socket.Read(client.packet.Buffer[client.packet.Index:])
		if err != nil {
			return
		}
 
		client.packet.Index += bl

		//enter when we recive enough bytes to start reading them
		for client.packet.Index > 2 {
			p := client.packet
			size := p.Index
			p.Index = 0
			
			//Check header byte
			if p.ReadByte() != 0xAA {
				client.Log().Printf("Wrong packet header")
				client.Log().Printf("% #X", p.Buffer[:size])
				return
			}
			l := int(p.ReadUInt16())
			p.Index = size
			
			if len(client.packet.Buffer) < l {
				client.packet.Resize(l)
			}

			//enter when we recived enough packet data
			if size >= l+3 {
				temp := client.packet.Buffer[:l+3]
				op := client.packet.Buffer[3]
				
				//check if packet is encrypted
				if op > 13 || (op > 1 && op < 5) || (op > 6 && op < 13) {
					var sumCheck bool
					temp, sumCheck = DecryptPacket(temp)
					if !sumCheck {
						client.Log().Println("Packet sum check failed!")
						return
					}
				} else {
					temp = temp[3:]
				}
				
				//handle packet
				client.ParsePacket(NewPacketRef(temp))
				client.packet.Index = 0
				
				//enter when we have more than one packet in buffer
				if size > l+3 {
					client.packet.Index = size - (l + 3)
					copy(client.packet.Buffer, client.packet.Buffer[l+3:size])
				} else {
					//enter when we done processing the buffer
					//keeping the user under 4048k use to save memory
					if cap(client.packet.Buffer) > 4048 {
						client.Buffer = make([]byte, 1024)
						client.packet = NewPacketRef(client.Buffer)
					}
				}
			} else {
				//break if we didn't get all the packet bytes
				break
			}
		}
	}
}

func (client *LClient) OnConnect() {
	client.packet = NewPacketRef(client.Buffer)
	client.packet.Index = 0
	client.SendWelcome()
	client.StartRecive()
}

func (client *LClient) OnDisconnect() {
	if x := recover(); x != nil {
		client.Log().Printf("panic : %v \n %s", x, C.PanicPath())
	}
	client.Socket.Close()
	client.MainServer.GetServer().Log.Println("Client Disconnected!")
}

func (client *LClient) Send(p *SGPacket) {
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

func (client *LClient) SendRaw(p *SGPacket) {
	p.WriteLen()
	client.Socket.Write(p.Buffer[:p.Index])
}

func (client *LClient) SendWelcome() {
	SendWelcome(client)
}

func (client *LClient) Log() *log.Logger {
	return Server.Log
}

func (client *LClient) ParsePacket(p *SGPacket) {
	header := p.ReadByte()

	fnc, exist := Handler[int(header)]
	if !exist {
		client.Log().Printf("Header(%d) len(%d) isnt registred : % #X\n %s", header, len(p.Buffer), p.Buffer, p.Buffer)
		return
	}
	//client.MainServer.GetServer().Log.Printf("Header(%d) len(%d) : % #X\n %s", header, len(p.Buffer), p.Buffer, p.Buffer)
	fnc(client, p)
}
