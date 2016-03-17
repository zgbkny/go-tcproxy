package packet

import (
    "log"
    "utils"
)

/******************************************
 * sessionId + length + data
 ******************************************/

type Packet struct {
    LOG             *log.Logger
    Start           int
    
    RawData			[]byte
    
    SessionId		uint32
    
    Length          int
}

func CreateNewPacket(LOG *log.Logger) *Packet {
    p := new(Packet)
    p.LOG = LOG
    return p
}

func ConstructPacket(rawData []byte, sessionId uint32, LOG *log.Logger) *Packet {
    p := new(Packet)
    p.LOG = LOG
    p.SessionId = sessionId
    p.RawData = rawData
    p.Start = 90;
    
    sessionIdBytes := utils.Uint32ToBytes(sessionId)
    copy(p.RawData[90:], sessionIdBytes)
    
    dataLenBytes := utils.Int16ToBytes(len(p.RawData) - 96)
    copy(p.RawData[94:], dataLenBytes)
    
    return p
}

func (p *Packet)GetPacket() []byte {
	return p.RawData[p.Start:]
}