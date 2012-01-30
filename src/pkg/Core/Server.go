package Core

import (
	"errors"
	"log"
	"net"
	"os"
	"strconv"
)

type Server interface {
	OnConnect(socket *net.TCPConn)
	OnSetup()
	Server() *CoreServer
}

type CoreServer struct {
	Socket   *net.TCPListener
	Addr     *net.TCPAddr
	Log      *Logger
	Name     string
	server	 Server
} 

func (serv *CoreServer) OnConnect(socket *net.TCPConn) {
	serv.Log.Printf("Client connected!")
	client := new(CoreClient)
	SetupClient(client, socket, serv)
}

func (serv *CoreServer) OnSetup() {
	serv.Log.Printf("Server:[%s] has started on %s", serv.Name, serv.Addr)
}

func (serv *CoreServer) Server() *CoreServer {
	return serv
}

func Start(iServ Server, name, ip string, port int) (err error) {

	serv := iServ.Server()
	if serv == nil {
		serv.Log.Printf("The Server struct is nil")
		return errors.New("The Server struct is nil")
	}
 	 
	serv.Log = NewLogger(os.Stderr, "["+name+"]", log.Ltime|log.Lshortfile)
	serv.Name = name
	serv.server = iServ
	serv.Addr, err = net.ResolveTCPAddr("tcp", ip+":"+strconv.Itoa(port))
	if err != nil {
		serv.Log.Printf("Server start failed %s", err.Error())
		return err
	}

	serv.Socket, err = net.ListenTCP("tcp", serv.Addr)

	if err != nil {
		serv.Log.Printf("Server start failed %s", err.Error())
		return err
	}

	iServ.OnSetup()

	return nil
}

func (serv *CoreServer) AcceptClients() {

	if serv == nil {
		serv.Log.Printf("The Server struct is nil")
		return
	}

	iServ := serv.server

	conn_in := make(chan *net.TCPConn, 20)
	go serv.acceptClients(conn_in)
	for {
		select {
		case c := <-conn_in:
			c.SetNoDelay(true)
			iServ.OnConnect(c)
		}
	}
}

func (serv *CoreServer) acceptClients(in chan *net.TCPConn) {
	for {
		c, err := serv.Socket.AcceptTCP()
		if err != nil {
			serv.Log.Printf("Server accept failed %s", err.Error())
		}
		in <- c
	}
}
