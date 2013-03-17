package SG

import (
	"encoding/xml"
	"log"
	"net"
	"os"
	"strconv"
)

var (
	Config  *ConfigData
	GSAddr  *net.TCPAddr
	LSAddr  *net.TCPAddr
	RPCAddr *net.TCPAddr
)

type ConfigData struct {
	XMLName   xml.Name           `xml:"Config"`
	LSConfig  *LoginServerConfig `xml:"LoginServer"`
	GSConfig  *GameServerConfig  `xml:"GameServer"`
	RPCConfig *RPCConfig         `xml:"RPC"`
}

type LoginServerConfig struct {
	IP, WANIP string
	Port      int
}

type GameServerConfig struct {
	IP, WANIP string
	Port      int
}

type RPCConfig struct {
	IP, WANIP string
	Port      int
}

func ReadConfig() {
	f, e := os.Open(configPath)
	if e != nil {
		f, e = CreateConfig()
		if e != nil {
			log.Panicln(e)
		} else {
			Initialize()
			return
		}
	}

	defer f.Close()

	Config = new(ConfigData)
	e = xml.NewDecoder(f).Decode(Config)
	if e != nil {
		log.Panicln(e)
	}

	Initialize()
}

func Initialize() {
	var err error
	LSAddr, err = net.ResolveTCPAddr("tcp", Config.LSConfig.WANIP+":"+strconv.Itoa(Config.LSConfig.Port))
	if err != nil {
		panic(err)
	}

	GSAddr, err = net.ResolveTCPAddr("tcp", Config.LSConfig.WANIP+":"+strconv.Itoa(Config.GSConfig.Port))
	if err != nil {
		panic(err)
	}

	RPCAddr, err = net.ResolveTCPAddr("tcp", Config.RPCConfig.WANIP+":"+strconv.Itoa(Config.RPCConfig.Port))
	if err != nil {
		panic(err)
	}
}

func GetGSIP() {

}

func CreateConfig() (f *os.File, e error) {
	f, e = os.Create(configPath)
	if e != nil {
		return nil, e
	}

	defer f.Close()

	Config = &ConfigData{
		LSConfig: &LoginServerConfig{"127.0.0.1", "127.0.0.1", 3000},
		GSConfig: &GameServerConfig{"127.0.0.1", "127.0.0.1", 13010},
		RPCConfig:  &RPCConfig{"127.0.0.1", "127.0.0.1", 1234},
	}


	b,e := xml.MarshalIndent(Config,"","\t")
	if e != nil {
		log.Panicln(e)
	}
	_, e = f.Write(b)
	if e != nil {
		log.Panicln(e)
	}

	return f, nil
}
