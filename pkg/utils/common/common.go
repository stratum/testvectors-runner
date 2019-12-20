package common

import "encoding/binary"

func GetUint32(i uint32) []byte {
	byteSlice := make([]byte, 4)
	binary.BigEndian.PutUint32(byteSlice, i)
	return byteSlice
}
