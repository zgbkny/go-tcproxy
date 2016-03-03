package tcptunnel

import (
    "log"
    "packet"
    "sync"
    "net"
)

type TcpTunnel struct {
    LOG                 *log.Logger
    
    send                chan []byte
    
    id                  uint32
    
    lock                *sync.Mutex
    
    OnDataF             func(*packet.Packet) int
    
    ReuseRate           int
    
    C                   *net.Conn
}

func CreateNewClientTunnel(id uint32, OnDataF func(*packet.Packet) int, LOG *log.Logger) *TcpTunnel {
    tt := new(TcpTunnel)
    tt.id = id
    tt.send = make(chan []byte)
    tt.lock = new(sync.Mutex)
    tt.OnDataF = OnDataF
    tt.LOG = LOG
    
    if !tt.connectToRemote() {
        return nil
    }
    go tt.recvData() 
    go tt.sendData()
    return tt
}

func CreateNewServerTunnel(id uint32, OnDataF func(*packet.Packet) int, conn *net.Conn, LOG *log.Logger) *TcpTunnel {
    tt := new(TcpTunnel)
    tt.id = id
    tt.send = make(chan []byte)
    tt.lock = new(sync.Mutex)
    tt.OnDataF = OnDataF
    tt.LOG = LOG
    tt.C = conn
    go tt.recvData()
    go tt.sendData()
    return tt
}

func (tt *TcpTunnel)connectToRemote() bool {
	conn, err := net.Dial("tcp", "localhost:9001")
	if err != nil {
		log.Println("connectToRemote", err)
		return false
	}
	tt.C = &conn
	return true
}

func (tt *TcpTunnel)SendPacket(p *packet.Packet) {
    tt.send <- p.GetPacket()
}

func (tt *TcpTunnel)recvData() {
    for {
		//log.Println("udptunnel tunnelReadFromClientProxy")
		data := make([]byte, 4096)
		_, err := (*tt.C).Read(data)
		//log.Println("after read", n)
		if err != nil {
			return
		}
		//log.Println("data len", n)
		//go tt.readPacketFromClientProxy(data[:n])
	}
}

func (tt *TcpTunnel)sendData() {
    for {
		//log.Println("udptunnel tunnelWrite")
		data, ok := <-tt.send
		if !ok {
			break
		}
		//log.Println("connWrite", string(data))
		(*tt.C).Write(data)
	}
}



