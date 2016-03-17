package main
import (
	"log"
    "tcptunnel"
    "session"
    "packet"
    "sync"
    "net"
    "os"
)

var LOG *log.Logger
var tunnelCount uint32        // 产生tunnelId
var idSessionMap map[uint32]*session.Session
var idTunnelMap map[uint32]*tcptunnel.TcpTunnel
var TunnelIdAllocLock *sync.Mutex
var lock *sync.Mutex

func getSession(id uint32) *session.Session {
	LOG.Println("server getSession")
	s, ok := idSessionMap[id]
	if !ok {
		lock.Lock()
		defer lock.Unlock()
		s, ok = idSessionMap[id]
		if !ok {
			s = session.CreateNewSession(id, LOG)
			idSessionMap[id] = s
			ok := connectToServer(s)
			if !ok {
				delete (idSessionMap, id)
				s.Destroy(false)
				return nil
			}
			
			go processRead(s)
		}
	}
	return s
}

func releaseSession(id uint32, flag bool) {
    
}

func connectToServer(s *session.Session) bool {
	LOG.Println("server connectToServer")
	conn, err := net.Dial("tcp", "localhost:90")
	if err != nil {
		log.Println("connect to nginx proxy", err)
		return false
	}
	s.C = &conn
    return true
}

func onData(p *packet.Packet) int {
	LOG.Println("server onData")
    s := getSession(p.SessionId)
	if s == nil  {
		return -1
	}
	if (p.Length == 0) {
		releaseSession(s.GetId(), false)
	}
    
    s.Slock.Lock()
    processWrite(s, p.GetPacket())
    s.Slock.Unlock()
    return 0
} 

func processWrite(s *session.Session, data []byte) {
	LOG.Println("server processWrite")
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
}

func processRead(s *session.Session) {
    conn := *s.C
	id := s.GetId()
    tunnelId := s.GetTunnelId()
    tt := idTunnelMap[tunnelId]
	for {
		LOG.Println("server processRead")
		/////////////////////////////////////////////////
		buf := make([]byte, 4096)
		length, err := conn.Read(buf[96:])
		if err != nil {
			LOG.Println("server read error", err)
			releaseSession(id, true)
			break
		}
		/////////////////////////////////////////////////
		
        p := packet.ConstructPacket(buf[:length + 96], id, LOG)  
        tt.SendPacket(p)
	}
	conn.Close()
}

func processNewAcceptedConn(c net.Conn) *tcptunnel.TcpTunnel {
	LOG.Println("processNewAcceptedConn")
    tt := tcptunnel.CreateNewServerTunnel(tunnelCount, onData, &c, LOG)
    idTunnelMap[tunnelCount] = tt
    tunnelCount++
    
    return tt
}

func initListen() {
	LOG.Println("initListen")
    listener, err := net.Listen("tcp", "0.0.0.0:9001")
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
		processNewAcceptedConn(conn)
	}
}

func Run() {
    initListen()
}

func main() {
	//tcpServer.TcpServer()
	
    
    fileName := "server_debug.log"
    logFile,err  := os.Create(fileName)
    defer logFile.Close()
    if err != nil {
        LOG.Fatalln("open file error !")
    }
    LOG = log.New(logFile,"[Debug]",log.Llongfile)
    LOG.Println("server")
    tunnelCount = 0
    
    idSessionMap = map[uint32]*session.Session{}
    idTunnelMap = map[uint32]*tcptunnel.TcpTunnel{}
	lock = new(sync.Mutex)
    
    initListen()
}