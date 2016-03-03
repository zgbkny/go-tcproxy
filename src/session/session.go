package session

import (
    "log"
    "sync"
    "net"
)

type Session struct {
    LOG             *log.Logger
	id				uint32
    tunnelId        uint32
    
    Slock			*sync.Mutex
    
    C				*net.Conn
}


func CreateNewSession(sessionId uint32, LOG *log.Logger) *Session {
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