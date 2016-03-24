package packet

import (
    "testing"
)

func createPacket() *Packet {
    return nil
}

func TestGetPacket(t *testing.T) {
    p := createPacket()
    t.Errorf("p data :", p.GetPacket())
}

