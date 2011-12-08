package Core

import (
	"errors"
	"log"
	"net"
	"os"
	"strconv"
)

type IServer interface {
	OnConnect(socket *net.TCPConn)
	OnSetup()
	GetServer() *Server
	GetIServer() IServer
}

type Server struct {
	Socket   *net.TCPListener
	Addr     *net.TCPAddr
	IP       string
	Port     int
	Log      *log.Logger
	Status   int
	OnClient chan *net.TCPConn
	Name     string
	IServ    IServer
} 

func (serv *Server) OnConnect(socket *net.TCPConn) {
	serv.Log.Printf("Client connected!")
	client := new(Client)
	SetupClient(client, socket, serv)
}

func (serv *Server) OnSetup() {
	go serv.handleClients()
	serv.Log.Printf("Server:[%s] has started on %s:%d", serv.Name, serv.IP, serv.Port)
}

func (serv *Server) GetServer() *Server {
	return serv
}

func (serv *Server) GetIServer() IServer {
	return serv.IServ
}

func Start(iServ IServer, name, ip string, port int) (err error) {

	serv := iServ.GetServer()
	if serv == nil {
		serv.Log.Printf("The Server struct is nil")
		return errors.New("The Server struct is nil")
	}
 
	serv.Log = log.New(os.Stderr, "["+name+"]", log.Ltime|log.Lshortfile)
	serv.IServ = iServ
	serv.IP = ip
	serv.Status = -1
	serv.Port = port
	serv.Name = name
	serv.Addr, err = net.ResolveTCPAddr("tcp", serv.IP+":"+strconv.Itoa(serv.Port))
	if err != nil {
		serv.Log.Printf("Server start failed %s", err.Error())
		return err
	}

	serv.Socket, err = net.ListenTCP("tcp", serv.Addr)

	if err != nil {
		serv.Log.Printf("Server start failed %s", err.Error())
		return err
	}

	serv.Status = 0
	iServ.OnSetup()

	return nil
}

func (serv *Server) handleClients() {

	if serv == nil {
		serv.Log.Printf("The Server struct is nil")
		return
	}

	if serv.Status == -1 {
		serv.Log.Printf("Server is down!")
		return
	}

	iServ := serv.GetIServer()

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

func (serv *Server) acceptClients(in chan *net.TCPConn) {
	for {
		c, err := serv.Socket.AcceptTCP()
		if err != nil {
			serv.Log.Printf("Server accept failed %s", err.Error())
		}
		in <- c
	}
}
