package main
import (
	"log"
    "session"
    "packet"
    "tcptunnel"
    "sync"
    "net"
    "os"
)

var LOG *log.Logger
var sessionCount uint32				// 产生sessionId
var tunnelCount uint32              // 产生tunnelId
var idSessionMap map[uint32]*session.Session
var idTunnelMap map[uint32]*tcptunnel.TcpTunnel
var TunnelIdAllocLock *sync.Mutex

const REUSE_RATE = 5

func releaseSession(id uint32, flag bool) {
    
}

func allocTunnelId() uint32 {
    LOG.Println("client allocTunnelId")
    TunnelIdAllocLock.Lock()
    defer TunnelIdAllocLock.Unlock()
    var id uint32
    for k, v := range idTunnelMap {
        id = k
        if v.ReuseRate < REUSE_RATE {
            v.ReuseRate++
            return k
        }
    }
    id++
    tt := tcptunnel.CreateNewClientTunnel(id, onData, LOG)
    if tt != nil {
        idTunnelMap[id] = tt
    } else {
        LOG.Println("tunnel connect error")
    } 
    return id
}

func onData(p *packet.Packet) int {
    LOG.Println("client onData")
    s, ok := idSessionMap[p.SessionId]
	if !ok {
		return -1
	}
	if (p.Length == 0) {
		releaseSession(s.GetId(), false)
        return 0
	}
    
    s.Slock.Lock()
    processWrite(s, p.GetPacket())
    s.Slock.Unlock()
    return 0
}

func processWrite(s *session.Session, data []byte) {
    LOG.Println("client processWrite")
    conn := *s.C
	id := s.GetId()
	index := 0

	for {
		length, err := conn.Write(data[index:])
		if err != nil {
			releaseSession(id, true)
			//return -1
		}
		if length != len(data) {
			index += length
		} else {
			break
		}
	}
	//return 0
}

func processRead(s *session.Session) {
    LOG.Println("client processRead")
    conn := *s.C
	id := s.GetId()
    tunnelId := s.GetTunnelId()
    tt := idTunnelMap[tunnelId]
	for {
        
		buf := make([]byte, 4096)
		length, err := conn.Read(buf[96:])
		if err != nil {
			LOG.Println("client read error", err)
			releaseSession(id, true)
			break
		}
		
        p := packet.ConstructPacket(buf[:length + 96], id, LOG)  
        tt.SendPacket(p)
	}
	conn.Close()
}


func processNewAcceptedConn(conn net.Conn) *session.Session {
    LOG.Println("client processNewAcceptedConn")
    s := session.CreateNewSession(sessionCount, LOG)
    s.C = &conn;
    // 分配tcptunnel
    s.SetTunnelId(allocTunnelId())
    
    idSessionMap[sessionCount] = s
    sessionCount++
    return s
}

func initListen() {
    LOG.Println("initListen")
    // create listener
	listener, err := net.Listen("tcp", "0.0.0.0:9000")
	if err != nil {
		return
	}

	// listen and accept connections from clients
	for {
		conn, err := listener.Accept()
		if err != nil {
			return
		}
        LOG.Println("Accept")
		s := processNewAcceptedConn(conn)
		// load balance, then process conn
		go processRead(s)
	}
}

func Run() {
    initListen()
}

func main() {
	//tcpServer.TcpServer()
    fileName := "client_debug.log"
    logFile,err  := os.Create(fileName)
    defer logFile.Close()
    if err != nil {
        LOG.Fatalln("open file error !")
    }
    LOG = log.New(logFile,"[Debug]",log.Llongfile)
    
    LOG.Println("client")
    
    sessionCount = 0
    
    idSessionMap = map[uint32]*session.Session{}
    idTunnelMap = map[uint32]*tcptunnel.TcpTunnel{}
    
    TunnelIdAllocLock = new(sync.Mutex)
    Run()
}
