package network

import (
	"unsafe"

	"github.com/StoneTrench/go-mc-clone/game/events"
)

/*

All numbers are big-endian.


*/

type EventBus = events.EventBus

const MAX_SAFE_BUFFER_SIZE = 1 << 9
const NETWORK_PROTOCOL = "udp"

// PROTOCOL_VERSION is incremented on breaking changes.
const PROTOCOL_VERSION = 0

const (
	// This peer has no peer id yet.
	PeerId_ConnectionRequest PeerId = 0
	// Peer id of the server.
	PeerId_Server PeerId = 1
	// Largest peer id which may not be used.
	PeerId_MaxStatic PeerId = PeerId_Server
	// Next largest peer id which may be used after the ones which cannot be used.
	PeerId_MaxStaticNext PeerId = PeerId_MaxStatic + 1
	// Largest peer id which may be used.
	PeerId_Max PeerId = (1 << (unsafe.Sizeof(PeerId(0)) * 8)) - 1
)

const MAX_CLIENTS_ON_SERVER = int(PeerId_Max - PeerId_MaxStatic)

const PROTOCOL_ID uint32 = 0x836414d1
const PACKET_HEADER_SIZE = int(unsafe.Sizeof(PeerId(0)) + unsafe.Sizeof(PROTOCOL_ID))

type PeerId uint32

// ParsePacket returns ok as true if it parsed the packet,
// or false if it didn't parse any packet.
func ParsePacket(buf []byte) (peer_id PeerId, payload []byte, ok bool) {
	if len(buf) < PACKET_HEADER_SIZE || GetUint[uint32](buf, 0) != PROTOCOL_ID {
		return 0, nil, false
	}
	peer_id = GetUint[PeerId](buf, 4)
	return peer_id, buf[PACKET_HEADER_SIZE:], true
}

func CreatePacket(peer_id PeerId, payload []byte) []byte {
	packet := make([]byte, PACKET_HEADER_SIZE+len(payload))
	PutUint(packet, 0, PROTOCOL_ID)
	PutUint(packet, 4, peer_id)
	copy(packet[PACKET_HEADER_SIZE:], payload)
	return packet
}

func PutUint[T ~uint8 | ~uint16 | ~uint32 | ~uint64](b []byte, offset int, value T) {
	switch unsafe.Sizeof(value) {
	case 1:
		b[offset] = byte(value)
	case 2:
		_ = b[offset+1]
		e := uint16(value)
		b[offset] = byte(e >> 8)
		b[offset+1] = byte(e)
	case 4:
		_ = b[offset+3]
		e := uint32(value)
		b[offset] = byte(e >> 24)
		b[offset+1] = byte(e >> 16)
		b[offset+2] = byte(e >> 8)
		b[offset+3] = byte(e)
	case 8:
		_ = b[offset+7]
		e := uint64(value)
		b[offset] = byte(e >> 56)
		b[offset+1] = byte(e >> 48)
		b[offset+2] = byte(e >> 40)
		b[offset+3] = byte(e >> 32)
		b[offset+4] = byte(e >> 24)
		b[offset+5] = byte(e >> 16)
		b[offset+6] = byte(e >> 8)
		b[offset+7] = byte(e)
	}
}

func GetUint[T ~uint8 | ~uint16 | ~uint32 | ~uint64](b []byte, offset int) T {
	var value T
	switch unsafe.Sizeof(value) {
	case 1:
		return T(b[offset])
	case 2:
		_ = b[offset+1]
		return T(
			uint16(b[offset])<<8 |
				uint16(b[offset+1]),
		)
	case 4:
		_ = b[offset+3]
		return T(
			uint32(b[offset])<<24 |
				uint32(b[offset+1])<<16 |
				uint32(b[offset+2])<<8 |
				uint32(b[offset+3]),
		)
	case 8:
		_ = b[offset+7]
		return T(
			uint64(b[offset])<<56 |
				uint64(b[offset+1])<<48 |
				uint64(b[offset+2])<<40 |
				uint64(b[offset+3])<<32 |
				uint64(b[offset+4])<<24 |
				uint64(b[offset+5])<<16 |
				uint64(b[offset+6])<<8 |
				uint64(b[offset+7]),
		)
	}
	panic("invalid type in GetUint")
}
