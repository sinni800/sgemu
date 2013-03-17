package Core

import (
	"net"
)

type Client interface {
	OnConnect()
	Client() *CoreClient
} 

type CoreClient struct {
	Socket     *net.TCPConn
	MainServer Server
	IP         string
}

func SetupClient(iClnt Client, socket *net.TCPConn, iServ Server) {

	client := iClnt.Client()
	if client == nil {
		iServ.Server().Log.Printf("The client struct is nil")
		return
	}

	client.Socket = socket
	client.MainServer = iServ
	ip, _, _ := net.SplitHostPort(socket.RemoteAddr().String())
	client.IP = ip
	go iClnt.OnConnect()
}

func (client *CoreClient) StartRecive() {
	Buffer := make([]byte, 1024)
	
	for {
		l, err := client.Socket.Read(Buffer)
		if err != nil {
			client.OnDisconnect()
			return
		}
		client.MainServer.Server().Log.Printf("Packet len %d", l)
	}
}

func (client *CoreClient) Client() *CoreClient {
	return client
}

func (client *CoreClient) OnConnect() {
	go client.StartRecive()
}

func (client *CoreClient) OnDisconnect() {
	client.Socket.Close()
	client.MainServer.Server().Log.Println("Client Disconnected!")
}
