package LoginServer

import (
	C "Core"
	D "Data"
	"log"
)
 
type LClient struct {
	C.Client
	Key      byte
	packet   *C.Packet
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

		for client.packet.Index > 2 {
			p := client.packet
			size := p.Index
			p.Index = 0
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

			if size >= l+3 {
				temp := client.packet.Buffer[:l+3]
				op := client.packet.Buffer[3]

				if op > 13 || (op > 1 && op < 5) || (op > 6 && op < 13) {
					var sumCheck bool
					temp, sumCheck = C.DecryptPacket(temp)
					if !sumCheck {
						client.Log().Println("Packet sum check failed!")
						return
					}
				} else {
					temp = temp[3:]
				}
				client.ParsePacket(C.NewPacketRef(temp))
				client.packet.Index = 0
				if size > l+3 {
					client.packet.Index = size - (l + 3)
					copy(client.packet.Buffer, client.packet.Buffer[l+3:size])
				} else {
					//keeping the user under 4048k use to save memory
					if cap(client.packet.Buffer) > 4048 {
						client.Buffer = make([]byte, 1024)
						client.packet = C.NewPacketRef(client.Buffer) 
					}
				}
			} else {
				break
			}
		}
	}
}


func (client *LClient) OnConnect() {
	client.packet = C.NewPacketRef(client.Buffer)
	client.packet.Index = 0
	client.SendWelcome()
	client.StartRecive()
}

func (client *LClient) OnDisconnect() {
	if x := recover(); x != nil {
			client.Log().Printf("panic : %v",x)
	}
	client.Socket.Close()
	client.MainServer.GetServer().Log.Println("Client Disconnected!")
}

func (client *LClient) Send(p *C.Packet) {
	if !p.Encrypted {
		op := p.Buffer[3]
		if op > 13 || (op > 0 && op < 3) || (op > 3 && op < 11) {
			p.WSkip(2)
			C.EncryptPacket(p.Buffer[:p.Index], client.Key)
			p.Encrypted = true
			client.Key++
		}
		p.WriteLen()
	}
	client.Socket.Write(p.Buffer[:p.Index])
}

func (client *LClient) SendRaw(p *C.Packet) {
	p.WriteLen()
	client.Socket.Write(p.Buffer[:p.Index])
}


func (client *LClient) SendWelcome() {
	SendWelcome(client)
}

func (client *LClient) Log() *log.Logger {
	return Server.Log
}

func (client *LClient) ParsePacket(p *C.Packet) {
	header := p.ReadByte()

	fnc, exist := Handler[int(header)]
	if !exist {
		client.Log().Printf("Header(%d) len(%d) isnt registred : % #X\n %s", header, len(p.Buffer), p.Buffer, p.Buffer)
		return
	}
	//client.MainServer.GetServer().Log.Printf("Header(%d) len(%d) : % #X\n %s", header, len(p.Buffer), p.Buffer, p.Buffer)
	fnc(client, p)
}
