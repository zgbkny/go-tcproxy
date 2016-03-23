package session

import (
    "log"
    "sync"
    "net"
)

var sessionIdCount uint32
var SessionIdCountLock *sync.Mutex

type Session struct {
    LOG             *log.Logger
	id				uint32
    tunnelId        uint32
    
    Slock			*sync.Mutex
    
    C				*net.Conn
}

func Init() {
    SessionIdCountLock = new(sync.Mutex)
    sessionIdCount = 0
}


func CreateNewSession(LOG *log.Logger) *Session {
    s := new(Session)
    s.LOG = LOG
    
    SessionIdCountLock.Lock()
    s.id = sessionIdCount
    sessionIdCount++
    SessionIdCountLock.Unlock()
    
    s.Slock = new(sync.Mutex)
    return s
}


func CreateNewSessionWithId(sessionId uint32, LOG *log.Logger) *Session {
    s := new(Session)
    s.LOG = LOG
    s.id = sessionId
    s.Slock = new(sync.Mutex)
    return s
}

func (s *Session)GetId() uint32 {
    return s.id
}

func (s *Session)GetTunnelId() uint32 {
    return s.tunnelId
}

func (s *Session)SetTunnelId(tunnelId uint32) {
    s.tunnelId = tunnelId
}

func (s *Session)Destroy(flag bool) {
    
}

func ProcessNewDataToServerProxy() {
    
}