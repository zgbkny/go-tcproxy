package utils
import (
	"bytes"
	"encoding/binary"
)

const MAX_ID uint64 = 4294967295
const MIN_ID uint64 = 0

func Uint32ToBytes(ui32 uint32) []byte {
	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.BigEndian, ui32)
	return buf.Bytes()
}

func BytesToUint32(data []byte) uint32 {
	buf := bytes.NewBuffer(data)
	var ui32 uint32
	binary.Read(buf, binary.BigEndian, &ui32)
	return ui32
}

func Uint64ToBytes(ui64 uint64) []byte {
    buf := bytes.NewBuffer([]byte{})
    binary.Write(buf, binary.BigEndian, ui64)
    return buf.Bytes()
}

func BytesToUint64(data []byte) uint64 {
    buf := bytes.NewBuffer(data)
    var ui64 uint64
    binary.Read(buf, binary.BigEndian, &ui64)
    return ui64
}

func Int16ToBytes(i int) []byte {
	buf := make([]byte, 2)
	buf[1] = byte(i & 0xff)
	buf[0] = byte((i >> 8) & 0xff)
	return buf
}

func BytesToInt16(data []byte) int {

	var i int
	i = 0
	i = int(byte(i) | data[0])
	i = i << 8
	i = int(byte(i) | data[1]) + i
	return i
}
