package tcptunnel

import (
    "log"
    "packet"
    "sync"
    "net"
    "utils"
)

var tunnelIdCount uint32
var tunnelIdCountLock *sync.Mutex

type TcpTunnel struct {
    LOG                 *log.Logger
    
    send                chan []byte
    
    id                  uint32
    
    lock                *sync.Mutex
    
    OnDataF             func(*packet.Packet) int
    
    ReuseRate           int
    
    C                   *net.Conn
}

func Init() {
    tunnelIdCountLock = new(sync.Mutex)
    tunnelIdCount = 0
}

func CreateNewClientTunnel(OnDataF func(*packet.Packet) int, LOG *log.Logger) *TcpTunnel {
    //LOG.Println("CreateNewClientTunnel")

    tt := new(TcpTunnel)
    
    tunnelIdCountLock.Lock()
    tt.id = tunnelIdCount
    tunnelIdCount++
    tunnelIdCountLock.Unlock()
    
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

func CreateNewServerTunnel(OnDataF func(*packet.Packet) int, conn *net.Conn, LOG *log.Logger) *TcpTunnel {
    tt := new(TcpTunnel)
    
    tunnelIdCountLock.Lock()
    tt.id = tunnelIdCount
    tunnelIdCount++
    tunnelIdCountLock.Unlock()
    
    tt.send = make(chan []byte)
    tt.lock = new(sync.Mutex)
    tt.OnDataF = OnDataF
    tt.LOG = LOG
    tt.C = conn
    go tt.recvData()
    go tt.sendData()
    return tt
}

func (tt *TcpTunnel)GetId() uint32 {
    return tt.id
}

func (tt *TcpTunnel)connectToRemote() bool {
    //tt.LOG.Println("tcptunnel connectToRemote")
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
		data := make([]byte, 6)
        headerDataIndex := 0
        for {
            n, err := (*tt.C).Read(data[headerDataIndex:])
            if err != nil {
                return
            } 
            if n < 6 - headerDataIndex {
                headerDataIndex += n
            } else {
                break
            }
        }
		
        
        //tt.LOG.Println("tcptunnel recvData header len:", n)
        sessionId := utils.BytesToUint32(data[:4])
        len := utils.BytesToInt16(data[4:6])   
        p := packet.CreateNewPacket(tt.LOG)
        p.SessionId = sessionId  
        p.Length = len
		if len != 0 {
            realData := make([]byte, len)
            realDataIndex := 0
            for {
                n, err := (*tt.C).Read(realData[realDataIndex:])
                if err != nil {
                    tt.LOG.Println("tcptunnel recvData error:", err)
                    return
                }
                if n < len - realDataIndex {
                    realDataIndex += n;
                } else {
                    break
                }
            } 
            p.RawData = realData
        }
		go tt.OnDataF(p)
	}
}

func (tt *TcpTunnel)sendData() {
    for {
		//tt.LOG.Println("sendData")
		data, ok := <-tt.send
        //tt.LOG.Println("sendData len:", len(data))
		if !ok {
            tt.LOG.Println("tcptunnel sendData error:", ok)
			break
		}
        
        index := 0
        for {
            length, err := (*tt.C).Write(data[index:])
            if err != nil {
                tt.LOG.Println("", err, length)
            }
            
            if length + index < len(data) {
                index += length
            } else {
                break
            }
        }
	}
}



