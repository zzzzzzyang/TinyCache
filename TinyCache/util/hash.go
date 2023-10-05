package util

import (
	"encoding/binary"
	"unsafe"
)

func DecodeFixed32(ptr unsafe.Pointer) uint32 {
	buffer := (*[4]byte)(ptr)
	return binary.LittleEndian.Uint32(buffer[:])
}

func Hash(data []byte, seed uint32) uint32 {
	const m = 0xc6a4a793
	const r = 24
	n := len(data)
	limit := n
	h := seed ^ (uint32(n) * m)

	for limit >= 4 {
		w := DecodeFixed32(unsafe.Pointer(&data[0]))
		data = data[4:]
		h += w
		h *= m
		h ^= (h >> 16)
		limit -= 4
	}

	switch limit {
	case 3:
		h += uint32(data[2]) << 16
	case 2:
		h += uint32(data[1]) << 8
	case 1:
		h += uint32(data[0])
		h *= m
		h ^= (h >> r)
	}

	return h
}
