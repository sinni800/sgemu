package Core

import (
	"net"
)

type IClient interface {
	OnConnect()
	OnDisconnect()
	GetClient() *Client
}

type Client struct {
	Socket     *net.TCPConn
	MainServer IServer
	Buffer     []byte
	IP         string
}

func SetupClient(iClnt IClient, socket *net.TCPConn, iServ IServer) {

	client := iClnt.GetClient()
	if client == nil {
		iServ.GetServer().Log.Printf("The client struct is nil")
		return
	}

	client.Socket = socket
	client.MainServer = iServ
	client.Buffer = make([]byte, 1024)
	ip, _, _ := net.SplitHostPort(socket.RemoteAddr().String())
	client.IP = ip
	go iClnt.OnConnect()
}

func (client *Client) StartRecive() {

	l, err := client.Socket.Read(client.Buffer)
	if err != nil {
		client.OnDisconnect()
		return
	}
	client.MainServer.GetServer().Log.Printf("Packet len %d", l)

	go client.StartRecive()
}

func (client *Client) GetClient() *Client {
	return client
}

func (client *Client) OnConnect() {
	go client.StartRecive()
}

func (client *Client) OnDisconnect() {
	client.Socket.Close()
	client.MainServer.GetServer().Log.Println("Client Disconnected!")
}
