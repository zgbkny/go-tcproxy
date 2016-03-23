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
var idSessionMap map[uint32]*session.Session
var idTunnelMap map[uint32]*tcptunnel.TcpTunnel

var IdSessionMapLock *sync.Mutex
var IdTunnelMapLock *sync.Mutex

const REUSE_RATE = 5

/////////////////////////////////////////////////////////////////////

func getSession(id uint32) (*session.Session, bool) {
    IdSessionMapLock.Lock()
    defer IdSessionMapLock.Unlock()
    s, ok := idSessionMap[id]
    
    return s, ok
}

func setSession(id uint32, s *session.Session) {
    IdSessionMapLock.Lock()
    defer IdSessionMapLock.Unlock()
    idSessionMap[id] = s
}

func releaseSession(id uint32, flag bool) {
    IdSessionMapLock.Lock()
    defer IdSessionMapLock.Unlock()
    s, ok := idSessionMap[id]
    if !ok {
        return;
    }
    (*s.C).Close()
    delete(idSessionMap, id)
}

/////////////////////////////////////////////////////////////////////

func getTunnel(id uint32) (*tcptunnel.TcpTunnel, bool) {
    IdTunnelMapLock.Lock()
    defer IdTunnelMapLock.Unlock()
    s, ok := idTunnelMap[id]
    
    return s, ok
}

func setTunnel(id uint32, t *tcptunnel.TcpTunnel) {
    //IdTunnelMapLock.Lock()  
    idTunnelMap[id] = t
    //IdSessionMapLock.Unlock()
}

func releaseTunnel(id uint32, flag bool) {
    IdTunnelMapLock.Lock()
    defer IdTunnelMapLock.Unlock()
    _, ok := idTunnelMap[id]
    if !ok {
        return;
    }
    delete(idTunnelMap, id)
}

/////////////////////////////////////////////////////////////////////

func allocTunnelId() (uint32, bool) {
    LOG.Println("client allocTunnelId")
    IdTunnelMapLock.Lock()
    defer IdTunnelMapLock.Unlock()
    for k, v := range idTunnelMap {
        if v.ReuseRate < REUSE_RATE {
            v.ReuseRate++
            return k, true
        }
    }
    return 0, false
}


func processNewAcceptedConn(conn net.Conn) *session.Session {
    LOG.Println("client processNewAcceptedConn")
    s := session.CreateNewSession(LOG)
    s.C = &conn;

    tunnelId, ok := allocTunnelId()
    if !ok {
 
        tt := tcptunnel.CreateNewClientTunnel(onData, LOG)
        
        if tt != nil {
            setTunnel(tt.GetId(), tt)
            s.SetTunnelId(tt.GetId())
        } else {
            LOG.Println("tunnel connect error")
            return nil
        } 
    } else {
        s.SetTunnelId(tunnelId)
    }
    
    setSession(s.GetId(), s)
    return s
}

//////////////////////////////////////////////////////////////////////

func onData(p *packet.Packet) int {
    LOG.Println("client onData")
    s, ok := getSession(p.SessionId)
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
			return
		}
		if length != len(data) {
			index += length
		} else {
			break
		}
	}
}

func processRead(s *session.Session) {
    LOG.Println("client processRead")
    conn := *s.C
	id := s.GetId()
    tunnelId := s.GetTunnelId()
    tt, _ := getTunnel(tunnelId)
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
    
    emptyBuf := make([]byte, 96)
    emptyP := packet.ConstructPacket(emptyBuf, id, LOG)
    tt.SendPacket(emptyP)
	conn.Close()
}

//////////////////////////////////////////////////////////////////////////

func initListen() {
	listener, err := net.Listen("tcp", "0.0.0.0:9000")
	if err != nil {
		return
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			return
		}
		s := processNewAcceptedConn(conn)
        if s != nil {
            go processRead(s)          
        }
	}
}

func Run() {
    initListen()
}

func main() {
    fileName := "client_debug.log"
    logFile,err  := os.Create(fileName)
    defer logFile.Close()
    if err != nil {
        LOG.Fatalln("open file error !")
    }
    LOG = log.New(logFile,"[Debug]",log.Llongfile)
    
    LOG.Println("client")
    
    
    idSessionMap = map[uint32]*session.Session{}
    idTunnelMap = map[uint32]*tcptunnel.TcpTunnel{}

    IdSessionMapLock = new(sync.Mutex)
    IdTunnelMapLock = new(sync.Mutex)
    
    session.Init()
    tcptunnel.Init()
    
    Run()
}
