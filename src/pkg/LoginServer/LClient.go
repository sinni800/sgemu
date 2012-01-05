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
	
	callback := func(p *SGPacket) { client.ParsePacket(p) }
	
	for {
		client.packet.ReadPacketFromStream(client.Socket,  callback)
	}
}

func (client *LClient) OnConnect() {
	client.packet = NewPacket()
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
